// Package typewriter provides a goldmark extension that converts typographic
// ("smart") characters back to their ASCII typewriter equivalents.
//
// It is the complement of goldmark's built-in typographer extension and tools
// like Smarty Pants. Content inside code spans and fenced blocks is left
// untouched by the extension form; use StripBytes to normalise everything
// including code content.
//
// All categories are active by default:
//
//	md := goldmark.New(goldmark.WithExtensions(typewriter.New()))
//
// Opt out of specific categories:
//
//	md := goldmark.New(goldmark.WithExtensions(
//	    typewriter.New(typewriter.WithoutCategory(typewriter.Math)),
//	))
//
// Enable only specific categories:
//
//	md := goldmark.New(goldmark.WithExtensions(
//	    typewriter.New(typewriter.WithCategory(typewriter.Dashes | typewriter.Ellipsis)),
//	))
//
// Override or remove individual conversions:
//
//	md := goldmark.New(goldmark.WithExtensions(
//	    typewriter.New(
//	        typewriter.WithMapping("—", "--"),  // prefer double-dash for em dash
//	        typewriter.WithMapping("×", ""),    // keep × as-is
//	    ),
//	))
package typewriter

import (
	"strings"
	"sync"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
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

// Option configures the extension.
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

// Extension is the goldmark extension.
type Extension struct {
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

// New creates the extension. With no options all Default categories are active.
func New(opts ...Option) *Extension {
	if len(opts) == 0 {
		return &Extension{r: getDefaultReplacer()}
	}
	cfg := &buildConfig{categories: Default}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.categories == Default && len(cfg.overrides) == 0 {
		return &Extension{r: getDefaultReplacer()}
	}
	return &Extension{r: buildReplacer(cfg.categories, cfg.overrides)}
}

// buildReplacer constructs a strings.Replacer from the active categories and
// any per-entry overrides. Overrides take precedence over builtins; an empty
// override value means "exclude this entry". strings.Replacer handles
// longest-match ordering internally — no pre-sorting required.
func buildReplacer(cats Category, overrides map[string]string) *strings.Replacer {
	args := make([]string, 0, len(builtinMappings)*2)
	seen := make(map[string]bool, len(overrides)+len(builtinMappings))

	// Overrides first so they shadow builtins.
	for from, to := range overrides {
		seen[from] = true
		if to != "" {
			args = append(args, from, to)
		}
	}

	// Add builtins for active categories.
	for _, m := range builtinMappings {
		if cats&m.cat == 0 || seen[m.from] {
			continue
		}
		seen[m.from] = true
		args = append(args, m.from, m.to)
	}

	return strings.NewReplacer(args...)
}

// StripBytes applies the extension's conversions directly to raw bytes, without
// any markdown parsing. Use this to normalise a markdown source before a
// subsequent goldmark pass, where the intermediate form must remain valid
// markdown rather than HTML.
func (e *Extension) StripBytes(src []byte) []byte {
	return []byte(e.r.Replace(string(src)))
}

// StripBytes is a package-level convenience that applies Default conversions
// to raw bytes. Equivalent to New().StripBytes(src).
func StripBytes(src []byte) []byte {
	return []byte(getDefaultReplacer().Replace(string(src)))
}

// Extend implements goldmark.Extender.
func (e *Extension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&transformer{r: e.r}, 100),
		),
	)
}
