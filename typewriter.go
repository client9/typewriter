// Package typewriter converts typographic ("smart") Unicode characters back to
// their plain ASCII typewriter equivalents.
//
// It is the complement of goldmark's typographer extension and tools like
// Smarty Pants. Use it as a preprocessor before a smart-typography pass to
// normalise mixed input to a consistent ASCII baseline.
//
// All categories are active by default:
//
//	clean := typewriter.ReplaceBytes(src)
//	clean := typewriter.Replace(s)
//
// Configure with options:
//
//	r := typewriter.New(
//	    typewriter.WithoutCategory(typewriter.Math),
//	    typewriter.WithMapping("—", "--"),
//	)
//	clean := r.Replace(s)
package typewriter

import (
	"strings"
	"sync"
)

// Category is a bitfield that groups related conversions.
type Category uint

const (
	Quotes    Category = 1 << iota // "curly" and «angle» quotes → straight ASCII
	Dashes                         // — en/em dashes → ---/--
	Ellipsis                       // … → ...
	Fractions                      // ½ ¼ ¾ → 1/2 1/4 3/4
	Symbols                        // © ® ™ § ¶ → (c) (r) (tm) S. P.
	Math                           // × ÷ ≠ ≤ ≥ → x / != <= >=
	Ligatures                      // ﬁ ﬂ ﬀ → fi fl ff
	Bullets                        // • · † ‡ → * . * **
	Spaces                         // non-breaking and width-variant Unicode spaces → plain ASCII space

	// Default is all defined categories.
	Default = Quotes | Dashes | Ellipsis | Fractions | Symbols | Math | Ligatures | Bullets | Spaces

	// CategoryAll is an alias for Default. Kept for forward compatibility in
	// case a future opt-in category is added.
	CategoryAll = Default
)

// buildConfig accumulates option values before the replacer is built.
type buildConfig struct {
	categories Category
	overrides  map[string]string // from → to; empty to means "exclude"
}

// Option configures a Replacer.
type Option func(*buildConfig)

// WithCategory sets the active categories to exactly c, replacing the default.
//
//	typewriter.WithCategory(typewriter.Quotes | typewriter.Dashes)
func WithCategory(c Category) Option {
	return func(cfg *buildConfig) { cfg.categories = c }
}

// WithoutCategory removes one or more categories from the active set.
//
//	typewriter.WithoutCategory(typewriter.Math | typewriter.Bullets)
func WithoutCategory(c Category) Option {
	return func(cfg *buildConfig) { cfg.categories &^= c }
}

// WithMapping adds or overrides a single conversion. Set to to an empty
// string to prevent a character from being converted at all.
//
//	typewriter.WithMapping("—", "--")  // prefer -- over ---
//	typewriter.WithMapping("×", "")    // leave × unchanged
func WithMapping(from, to string) Option {
	return func(cfg *buildConfig) {
		if cfg.overrides == nil {
			cfg.overrides = make(map[string]string)
		}
		cfg.overrides[from] = to
	}
}

// Replacer applies typographic-to-ASCII conversions. Create with New.
type Replacer struct {
	r *strings.Replacer
}

var (
	cachedDefault     *strings.Replacer
	cachedDefaultOnce sync.Once
)

func getDefaultReplacer() *strings.Replacer {
	cachedDefaultOnce.Do(func() {
		cachedDefault = buildReplacer(Default, nil)
	})
	return cachedDefault
}

// New creates a Replacer with the given options. With no options all Default
// categories are active.
func New(opts ...Option) *Replacer {
	if len(opts) == 0 {
		return &Replacer{r: getDefaultReplacer()}
	}
	cfg := &buildConfig{categories: Default}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.categories == Default && len(cfg.overrides) == 0 {
		return &Replacer{r: getDefaultReplacer()}
	}
	return &Replacer{r: buildReplacer(cfg.categories, cfg.overrides)}
}

// Replace returns s with all active typographic characters converted to ASCII.
func (r *Replacer) Replace(s string) string {
	return r.r.Replace(s)
}

// ReplaceBytes returns a copy of b with all active typographic characters
// converted to ASCII.
func (r *Replacer) ReplaceBytes(b []byte) []byte {
	return []byte(r.r.Replace(string(b)))
}

// Replace returns s with all Default typographic characters converted to ASCII.
func Replace(s string) string {
	return getDefaultReplacer().Replace(s)
}

// ReplaceBytes returns a copy of b with all Default typographic characters
// converted to ASCII.
func ReplaceBytes(b []byte) []byte {
	return []byte(getDefaultReplacer().Replace(string(b)))
}

// buildReplacer constructs a strings.Replacer from the active categories and
// any per-entry overrides. Overrides take precedence over builtins; an empty
// override value means "exclude this entry". strings.Replacer handles
// longest-match ordering internally.
func buildReplacer(cats Category, overrides map[string]string) *strings.Replacer {
	args := make([]string, 0, len(builtinMappings)*2)
	seen := make(map[string]bool, len(overrides)+len(builtinMappings))

	for from, to := range overrides {
		seen[from] = true
		if to != "" {
			args = append(args, from, to)
		}
	}

	for _, m := range builtinMappings {
		if cats&m.cat == 0 || seen[m.from] {
			continue
		}
		seen[m.from] = true
		args = append(args, m.from, m.to)
	}

	return strings.NewReplacer(args...)
}
