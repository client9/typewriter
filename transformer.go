package typewriter

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type transformer struct {
	pairs [][2]string
}

func (t *transformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	if len(t.pairs) == 0 {
		return
	}
	source := reader.Source()
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n.Kind() {
		case ast.KindCodeBlock, ast.KindFencedCodeBlock, ast.KindCodeSpan,
			ast.KindHTMLBlock, ast.KindRawHTML:
			// The HTML renderer reads code content directly from source[]
			// via segment.Value(source) and expects ast.Text children, not
			// ast.String. Replacing nodes here would panic the renderer.
			// Use StripBytes on the raw source before parsing to normalise
			// code content.
			return ast.WalkSkipChildren, nil
		case ast.KindText:
			replaceText(n.(*ast.Text), source, t.pairs)
		}
		return ast.WalkContinue, nil
	})
}

func replaceText(node *ast.Text, source []byte, pairs [][2]string) {
	content := node.Segment.Value(source)
	result := content
	for _, p := range pairs {
		result = bytes.ReplaceAll(result, []byte(p[0]), []byte(p[1]))
	}
	if bytes.Equal(result, content) {
		return
	}
	newNode := ast.NewString(result)
	node.Parent().ReplaceChild(node.Parent(), node, newNode)
}
