// Package typewriterext provides a goldmark extension that applies typewriter
// conversions as an AST transformer. For preprocessing raw markdown source
// before parsing, use the parent typewriter package directly.
//
// Type aliases and re-exported identifiers let callers use this package as
// their sole import — no need to also import github.com/client9/typewriter.
package typewriterext

import (
	"github.com/client9/typewriter"
	gm "github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

// Type aliases — identical to the core types, no conversion required.
type (
	Category = typewriter.Category
	Option   = typewriter.Option
)

// Category constants re-exported by value.
const (
	Quotes      = typewriter.Quotes
	Dashes      = typewriter.Dashes
	Ellipsis    = typewriter.Ellipsis
	Fractions   = typewriter.Fractions
	Symbols     = typewriter.Symbols
	Math        = typewriter.Math
	Ligatures   = typewriter.Ligatures
	Bullets     = typewriter.Bullets
	Spaces      = typewriter.Spaces
	Default     = typewriter.Default
	CategoryAll = typewriter.CategoryAll
)

// Option constructors re-exported as package-level vars.
var (
	WithCategory    = typewriter.WithCategory
	WithoutCategory = typewriter.WithoutCategory
	WithMapping     = typewriter.WithMapping
)

// Extension is the goldmark extension. Create with New.
type Extension struct {
	r *typewriter.Replacer
}

// New creates the goldmark extension. With no options all Default categories
// are active. Accepts the same options as typewriter.New.
func New(opts ...typewriter.Option) *Extension {
	return &Extension{r: typewriter.New(opts...)}
}

// Extend implements goldmark.Extender.
func (e *Extension) Extend(m gm.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&transformer{r: e.r}, 100),
		),
	)
}
