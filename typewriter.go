// Package typewriter converts typographic ("smart") Unicode characters back to
// their plain ASCII typewriter equivalents, and strips Unicode style variants
// (bold, italic, monospace, etc.) back to plain letters.
//
// Use package-level functions for the common case:
//
//	clean := typewriter.Replace(s)
//	clean := typewriter.ReplaceBytes(b)
//
// Configure with a Config struct:
//
//	r := typewriter.New(typewriter.Config{
//	    Categories: typewriter.Default,
//	    Runs: []typewriter.RunStyle{
//	        {Style: typewriter.Bold, Prefix: "**", Suffix: "**"},
//	    },
//	})
//	clean := r.Replace(s)
package typewriter

import (
	"strings"
	"sync"
	"unicode/utf8"
)

// Category is a bitfield grouping related character conversions.
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

// UnicodeStyle identifies a typographic Unicode style variant.
type UnicodeStyle int

const (
	Bold        UnicodeStyle = iota // 𝗔𝗕𝗖 → ABC
	Italic                          // 𝘈𝘉𝘊 → ABC
	BoldItalic                      // 𝘼𝘽𝘾 → ABC
	Monospace                       // 𝙰𝙱𝙲 → ABC
	Superscript                     // ²⁴  → 24
	Subscript                       // ₂₄  → 24
)

// RunStyle configures how a run of styled Unicode characters is converted.
// The ASCII text of the run is wrapped with Prefix and Suffix.
// Empty Prefix and Suffix strips the run to plain ASCII.
type RunStyle struct {
	Style  UnicodeStyle
	Prefix string
	Suffix string
}

// Config configures a Replacer.
type Config struct {
	Categories Category
	Overrides  map[string]string // from → to; empty to = pass through unchanged
	Runs       []RunStyle
}

// Replacer applies typographic-to-ASCII conversions. Create with New.
type Replacer struct {
	sr     *strings.Replacer
	runs   []RunStyle
	lookup map[rune]styledRune
}

var defaultReplacer = sync.OnceValue(func() *Replacer {
	return New(Config{Categories: Default})
})

// New creates a Replacer from cfg.
func New(cfg Config) *Replacer {
	return &Replacer{
		sr:     buildReplacer(cfg.Categories, cfg.Overrides),
		runs:   cfg.Runs,
		lookup: buildStyleLookup(cfg.Runs),
	}
}

// Replace returns s with all active conversions applied.
func (r *Replacer) Replace(s string) string {
	if len(r.runs) == 0 {
		return r.sr.Replace(s)
	}
	return r.replaceWithRuns(s)
}

// ReplaceBytes returns a copy of b with all active conversions applied.
func (r *Replacer) ReplaceBytes(b []byte) []byte {
	return []byte(r.Replace(string(b)))
}

func (r *Replacer) replaceWithRuns(s string) string {
	var buf strings.Builder
	buf.Grow(len(s))
	segStart := 0

	for i := 0; i < len(s); {
		ru, size := utf8.DecodeRuneInString(s[i:])
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
