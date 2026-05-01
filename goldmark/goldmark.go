// Package typewriterext provides a goldmark extension that applies typewriter
// conversions as an AST transformer. For preprocessing raw markdown source
// before parsing, use the parent typewriter package directly.
package typewriterext

import (
	"github.com/client9/typewriter"
	gm "github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
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
