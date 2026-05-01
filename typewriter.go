// Package typewriter provides a goldmark extension that converts typographic
// ("smart") characters back to their ASCII typewriter equivalents.
//
// It is the complement of goldmark's built-in typographer extension and tools
// like Smarty Pants. Content inside code spans and fenced blocks is left
// untouched.
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
// Opt in to space normalisation (off by default):
//
//	md := goldmark.New(goldmark.WithExtensions(
//	    typewriter.New(typewriter.WithSpaces()),
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
	"sort"
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
	Spaces                         // non-breaking/em/en/thin spaces → space (opt-in)

	// Default is all categories except Spaces, which can change document
	// semantics and is therefore opt-in via WithSpaces.
	Default = Quotes | Dashes | Ellipsis | Fractions | Symbols | Math | Ligatures | Bullets
)

// buildConfig accumulates option values before pairs are compiled.
type buildConfig struct {
	categories Category
	overrides  map[string]string // from → to; empty to means "exclude"
}

// Option configures the extension.
type Option func(*buildConfig)

// WithoutCategory disables one or more default categories.
//
//	typewriter.WithoutCategory(typewriter.Math | typewriter.Bullets)
func WithoutCategory(c Category) Option {
	return func(cfg *buildConfig) { cfg.categories &^= c }
}

// WithSpaces enables normalisation of non-breaking and other non-standard
// Unicode spaces to a plain ASCII space. Off by default because changing
// spaces can affect line-breaking intent.
func WithSpaces() Option {
	return func(cfg *buildConfig) { cfg.categories |= Spaces }
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
	pairs [][2]string // precomputed at New(); immutable thereafter
}

// defaultPairs caches the compiled pair list for the Default category set.
var (
	cachedDefault     [][2]string
	cachedDefaultOnce sync.Once
)

func getDefaultPairs() [][2]string {
	cachedDefaultOnce.Do(func() {
		cachedDefault = compilePairs(Default, nil)
	})
	return cachedDefault
}

// New creates the extension. With no options all Default categories are active.
func New(opts ...Option) *Extension {
	if len(opts) == 0 {
		return &Extension{pairs: getDefaultPairs()}
	}
	cfg := &buildConfig{categories: Default}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.categories == Default && len(cfg.overrides) == 0 {
		return &Extension{pairs: getDefaultPairs()}
	}
	return &Extension{pairs: compilePairs(cfg.categories, cfg.overrides)}
}

// compilePairs builds the sorted replacement table from active categories
// and any per-entry overrides. Overrides take precedence over builtins;
// an empty override value means "exclude this entry".
func compilePairs(cats Category, overrides map[string]string) [][2]string {
	// Reserve capacity: rough estimate of active entries.
	pairs := make([][2]string, 0, len(builtinMappings))

	// Track which sources are already decided (by override or previous builtin).
	seen := make(map[string]bool, len(overrides)+len(builtinMappings))

	// Overrides first so they shadow builtins.
	for from, to := range overrides {
		seen[from] = true
		if to != "" {
			pairs = append(pairs, [2]string{from, to})
		}
	}

	// Add builtins for active categories.
	for _, m := range builtinMappings {
		if cats&m.cat == 0 || seen[m.from] {
			continue
		}
		seen[m.from] = true
		pairs = append(pairs, [2]string{m.from, m.to})
	}

	// Sort longest source first. Most sources are single codepoints but a few
	// (like ligatures stored as multi-byte UTF-8) differ in byte length.
	sort.Slice(pairs, func(i, j int) bool {
		li, lj := len(pairs[i][0]), len(pairs[j][0])
		if li != lj {
			return li > lj
		}
		return pairs[i][0] < pairs[j][0]
	})
	return pairs
}

// Extend implements goldmark.Extender.
func (e *Extension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&transformer{pairs: e.pairs}, 100),
		),
	)
}
