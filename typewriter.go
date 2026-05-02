package typewriter

import (
	"strings"
	"sync"
	"unicode/utf8"
)

// Category is a bitmask that selects groups of character substitutions.
// Combine groups with |; exclude a group from [Default] with &^.
type Category uint

const (
	Quotes    Category = 1 << iota // curly/angle quotes → straight ASCII
	Dashes                         // em/en dashes → ---/--
	Ellipsis                       // … → ...
	Fractions                      // ½ ¼ ¾ → 1/2 1/4 3/4
	Symbols                        // © ® ™ § ¶ → (c) (r) (tm) S. P.
	Math                           // × ÷ ≠ ≤ ≥ → x / != <= >=
	Ligatures                      // ﬁ ﬂ ﬀ → fi fl ff
	Bullets                        // • · † ‡ → * . * **
	Spaces                         // non-breaking and width-variant spaces → plain space

	// Default is all defined categories.
	Default = Quotes | Dashes | Ellipsis | Fractions | Symbols | Math | Ligatures | Bullets | Spaces

	// CategoryAll is an alias for Default, kept for forward compatibility.
	CategoryAll = Default
)

// UnicodeStyle identifies a typographic Unicode style variant used in
// mathematical notation and social-media text.
//
// The zero value is StyleUnknown. Callers must set Style explicitly when
// constructing a [RunStyle]; a zero-value RunStyle{} has no defined style.
type UnicodeStyle int

const (
	StyleUnknown UnicodeStyle = iota // zero value; not a valid style
	Bold                             // sans-serif bold: 𝗔𝗕𝗖 → ABC
	Italic                           // sans-serif italic: 𝘈𝘉𝘊 → ABC
	BoldItalic                       // sans-serif bold-italic: 𝘼𝘽𝘾 → ABC
	Monospace                        // monospace: 𝙰𝙱𝙲 → ABC
	Superscript                      // superscript digits/letters: ²⁴ → 24
	Subscript                        // subscript digits: ₂₄ → 24
)

// RunStyle configures how a contiguous run of styled Unicode characters is
// converted. The recovered ASCII text is wrapped with Prefix and Suffix.
// When both are empty the run is stripped to plain ASCII with no added markup.
//
// Style must be set explicitly; [StyleUnknown] (zero value) matches nothing.
// If the same Style value appears more than once in [Config.Runs], only the
// first occurrence is used; subsequent duplicates are silently ignored.
type RunStyle struct {
	Style  UnicodeStyle // style variant to detect; must not be StyleUnknown
	Prefix string       // prepended to the recovered ASCII text
	Suffix string       // appended to the recovered ASCII text
}

// Config configures a [Replacer]. Pass as a value to [New].
type Config struct {
	// Categories selects which built-in conversion groups are active.
	// Use [Default] to enable all groups.
	Categories Category

	// Overrides adjusts individual mappings before the built-in table is
	// consulted. The key is the Unicode source string; the value is the ASCII
	// target. An empty value suppresses the built-in mapping for that key.
	Overrides map[string]string

	// Runs configures detection of contiguous Unicode-styled character runs
	// (bold, italic, monospace, etc.) and the Prefix/Suffix used to wrap
	// the recovered ASCII text. See [RunStyle].
	Runs []RunStyle
}

// Replacer applies typographic-to-ASCII conversions configured by [New].
// A Replacer is safe for concurrent use by multiple goroutines.
type Replacer struct {
	sr     *strings.Replacer
	runs   []RunStyle
	lookup map[rune]styledRune
}

var defaultReplacer = sync.OnceValue(func() *Replacer {
	return New(Config{Categories: Default})
})

// New returns a Replacer configured by cfg.
func New(cfg Config) *Replacer {
	return &Replacer{
		sr:     buildReplacer(cfg.Categories, cfg.Overrides),
		runs:   cfg.Runs,
		lookup: buildStyleLookup(cfg.Runs),
	}
}

// Replace returns s with all active conversions applied.
// If r is nil, the package-level [Default] replacer is used.
func (r *Replacer) Replace(s string) string {
	if r == nil {
		return defaultReplacer().Replace(s)
	}
	if len(r.runs) == 0 {
		return r.sr.Replace(s)
	}
	return r.replaceWithRuns(s)
}

// ReplaceBytes returns a copy of b with all active conversions applied.
// If r is nil, the package-level [Default] replacer is used.
func (r *Replacer) ReplaceBytes(b []byte) []byte {
	if r == nil {
		return defaultReplacer().ReplaceBytes(b)
	}
	return []byte(r.Replace(string(b)))
}

func (r *Replacer) replaceWithRuns(s string) string {
	var buf strings.Builder
	buf.Grow(len(s))
	segStart := 0

	for i := 0; i < len(s); {
		ru, size := utf8.DecodeRuneInString(s[i:])
		// utf8.RuneError with size==1 means invalid UTF-8; it won't be in
		// lookup so it falls through to the strings.Replacer segment below,
		// which handles arbitrary byte sequences safely.
		sr, ok := r.lookup[ru]
		if !ok {
			i += size
			continue
		}
		// Flush the non-run segment through the char replacer.
		if i > segStart {
			buf.WriteString(r.sr.Replace(s[segStart:i]))
		}
		rs := r.findRunStyle(sr.style)
		buf.WriteString(rs.Prefix)
		buf.WriteRune(sr.ascii)
		i += size
		// Consume the rest of the run.
		for i < len(s) {
			ru2, size2 := utf8.DecodeRuneInString(s[i:])
			sr2, ok2 := r.lookup[ru2]
			if !ok2 || sr2.style != sr.style {
				break
			}
			buf.WriteRune(sr2.ascii)
			i += size2
		}
		buf.WriteString(rs.Suffix)
		segStart = i
	}
	if segStart < len(s) {
		buf.WriteString(r.sr.Replace(s[segStart:]))
	}
	return buf.String()
}

func (r *Replacer) findRunStyle(style UnicodeStyle) RunStyle {
	for _, rs := range r.runs {
		if rs.Style == style {
			return rs
		}
	}
	return RunStyle{}
}

// Replace returns s with all Default conversions applied.
func Replace(s string) string { return defaultReplacer().Replace(s) }

// ReplaceBytes returns a copy of b with all Default conversions applied.
func ReplaceBytes(b []byte) []byte { return defaultReplacer().ReplaceBytes(b) }

func buildReplacer(cats Category, overrides map[string]string) *strings.Replacer {
	args := make([]string, 0, len(builtinMappings)*2)
	blocked := make(map[string]bool, len(overrides))

	for from, to := range overrides {
		blocked[from] = true
		if to != "" {
			args = append(args, from, to)
		}
	}
	for _, m := range builtinMappings {
		if cats&m.cat == 0 || blocked[m.from] {
			continue
		}
		args = append(args, m.from, m.to)
	}
	return strings.NewReplacer(args...)
}
