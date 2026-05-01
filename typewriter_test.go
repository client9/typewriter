package typewriter_test

import (
	"bytes"
	"testing"

	"github.com/client9/typewriter"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

func render(t *testing.T, ext *typewriter.Extension, src string) string {
	t.Helper()
	md := goldmark.New(goldmark.WithExtensions(ext))
	var buf bytes.Buffer
	if err := md.Convert([]byte(src), &buf); err != nil {
		t.Fatalf("Convert: %v", err)
	}
	return buf.String()
}

func TestQuotes(t *testing.T) {
	ext := typewriter.New()
	tests := []struct{ in, want string }{
		{"“Hello”", "<p>&quot;Hello&quot;</p>\n"},     // “Hello” (curly doubles)
		{"’it’s’", "<p>'it's'</p>\n"},                 // ‘it’s’ (curly singles)
		{"«hello»", "<p>&lt;&lt;hello&gt;&gt;</p>\n"}, // «hello» (angle)
		{"„low”", "<p>&quot;low&quot;</p>\n"},         // „low” (low-9 open)
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := render(t, ext, tt.in)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

func TestDashes(t *testing.T) {
	ext := typewriter.New()
	tests := []struct{ in, want string }{
		{"em—dash", "<p>em---dash</p>\n"},
		{"en–dash", "<p>en--dash</p>\n"},
		{"minus−sign", "<p>minus-sign</p>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := render(t, ext, tt.in)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

func TestEllipsis(t *testing.T) {
	got := render(t, typewriter.New(), "wait…")
	want := "<p>wait...</p>\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFractions(t *testing.T) {
	ext := typewriter.New()
	tests := []struct{ in, want string }{
		{"½", "<p>1/2</p>\n"},
		{"¼", "<p>1/4</p>\n"},
		{"¾", "<p>3/4</p>\n"},
		{"⅓", "<p>1/3</p>\n"},
		{"⅛", "<p>1/8</p>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := render(t, ext, tt.in)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

func TestSymbols(t *testing.T) {
	ext := typewriter.New()
	tests := []struct{ in, want string }{
		{"©", "<p>(c)</p>\n"},
		{"®", "<p>(r)</p>\n"},
		{"™", "<p>(tm)</p>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := render(t, ext, tt.in)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

func TestMath(t *testing.T) {
	ext := typewriter.New()
	tests := []struct{ in, want string }{
		{"10×", "<p>10x</p>\n"},
		{"10÷2", "<p>10/2</p>\n"},
		{"a≠b", "<p>a!=b</p>\n"},
		{"a≤b", "<p>a&lt;=b</p>\n"},
		{"a≥b", "<p>a&gt;=b</p>\n"},
		{"a→b", "<p>a-&gt;b</p>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := render(t, ext, tt.in)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

func TestLigatures(t *testing.T) {
	ext := typewriter.New()
	tests := []struct{ in, want string }{
		{"ﬁle", "<p>file</p>\n"},
		{"ﬂow", "<p>flow</p>\n"},
		{"ﬀect", "<p>ffect</p>\n"},
		{"ﬃcient", "<p>fficient</p>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := render(t, ext, tt.in)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

func TestBullets(t *testing.T) {
	ext := typewriter.New()
	tests := []struct{ in, want string }{
		{"item•one", "<p>item*one</p>\n"},
		{"note†ref", "<p>note*ref</p>\n"},
		{"note‡ref", "<p>note**ref</p>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := render(t, ext, tt.in)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

func TestSpaces(t *testing.T) {
	// NBSP should be converted to a plain space by default.
	got := render(t, typewriter.New(), "a b")
	if got != "<p>a b</p>\n" {
		t.Errorf("default: got %q", got)
	}

	// WithoutCategory(Spaces) preserves NBSP.
	got = render(t, typewriter.New(typewriter.WithoutCategory(typewriter.Spaces)), "a b")
	if got != "<p>a b</p>\n" {
		t.Errorf("WithoutCategory(Spaces): got %q", got)
	}
}

func TestCodeContent(t *testing.T) {
	// The goldmark extension form cannot modify code span or fenced code block
	// content: the HTML renderer reads that content directly from the original
	// source bytes (segment.Value(source)) and expects ast.Text children, not
	// ast.String. Replacing those nodes would panic the renderer.
	//
	// StripBytes operates before parsing and normalises everything, including
	// code content — which is the right behaviour for copy-paste corruption
	// (smart quotes inside a shell command, for example).
	src := "outside … but `inside … code`"

	// Extension form: prose converted, code preserved (architecture constraint).
	got := render(t, typewriter.New(), src)
	want := "<p>outside ... but <code>inside … code</code></p>\n"
	if got != want {
		t.Errorf("extension: got %q\nwant %q", got, want)
	}

	// StripBytes: everything converted, including inside code spans.
	stripped := string(typewriter.StripBytes([]byte(src)))
	got = render(t, typewriter.New(), stripped)
	want = "<p>outside ... but <code>inside ... code</code></p>\n"
	if got != want {
		t.Errorf("StripBytes: got %q\nwant %q", got, want)
	}
}

func TestWithCategory(t *testing.T) {
	t.Run("whitelist", func(t *testing.T) {
		// Only Ellipsis active; dashes and quotes should pass through
		ext := typewriter.New(typewriter.WithCategory(typewriter.Ellipsis))
		got := render(t, ext, "wait…")
		if got != "<p>wait...</p>\n" {
			t.Errorf("ellipsis: got %q", got)
		}
		got = render(t, ext, "em—dash")
		if got != "<p>em—dash</p>\n" {
			t.Errorf("dash should pass through: got %q", got)
		}
	})
	t.Run("all", func(t *testing.T) {
		// CategoryAll is an alias for Default; verify it compiles and works
		ext := typewriter.New(typewriter.WithCategory(typewriter.CategoryAll))
		got := render(t, ext, "wait…")
		if got != "<p>wait...</p>\n" {
			t.Errorf("got %q", got)
		}
	})
}

func TestWithoutCategory(t *testing.T) {
	ext := typewriter.New(typewriter.WithoutCategory(typewriter.Math))
	// × should pass through unchanged
	got := render(t, ext, "10×")
	want := "<p>10×</p>\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	// ellipsis should still work
	got = render(t, ext, "wait…")
	want = "<p>wait...</p>\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestWithMapping(t *testing.T) {
	t.Run("override", func(t *testing.T) {
		// Prefer -- over --- for em dash
		ext := typewriter.New(typewriter.WithMapping("—", "--"))
		got := render(t, ext, "em—dash")
		want := "<p>em--dash</p>\n"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	t.Run("exclude", func(t *testing.T) {
		// Empty target means pass through
		ext := typewriter.New(typewriter.WithMapping("×", ""))
		got := render(t, ext, "10×")
		want := "<p>10×</p>\n"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	t.Run("add_custom", func(t *testing.T) {
		// Add a mapping not in builtins
		ext := typewriter.New(typewriter.WithMapping("°", "deg"))
		got := render(t, ext, "90°")
		want := "<p>90deg</p>\n"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func BenchmarkNew(b *testing.B) {
	b.Run("no_options", func(b *testing.B) {
		for range b.N {
			_ = typewriter.New()
		}
	})
	b.Run("with_option", func(b *testing.B) {
		for range b.N {
			_ = typewriter.New(typewriter.WithoutCategory(typewriter.Math))
		}
	})
}

func renderWith(t *testing.T, src string, exts ...goldmark.Extender) string {
	t.Helper()
	md := goldmark.New(goldmark.WithExtensions(exts...))
	var buf bytes.Buffer
	if err := md.Convert([]byte(src), &buf); err != nil {
		t.Fatalf("Convert: %v", err)
	}
	return buf.String()
}

// TestChaining explores what happens when typewriter and goldmark's typographer
// are used together on mixed input (some chars already typographic, some ASCII).
//
// The desired outcome is consistent typographic output regardless of what the
// source contained. This test documents why that requires source-level stripping,
// not goldmark extension chaining.
func TestChaining(t *testing.T) {
	// Genuinely mixed input: first pair is Unicode curly quotes (U+201C/U+201D),
	// second pair is plain ASCII quotes. A consistent typographer pass should
	// produce the same quote style for both.
	mixed := "“already curly” and \"straight\""

	t.Logf("input: %s", mixed)

	// Case 1: typographer only, no typewriter.
	// The typographer is an inline parser — it fires during tokenisation and only
	// acts on ASCII trigger sequences. U+201C/U+201D are not ASCII, so it leaves
	// the first pair alone and converts only the second. Output is inconsistent.
	t.Logf("typographer only: %s", renderWith(t, mixed, extension.Typographer))

	// Case 2: single goldmark instance with both extensions.
	// Typographer (inline parser) fires first: converts "straight" → &ldquo;&rdquo;
	// stored in ast.String nodes. Typewriter (AST transformer) fires second: strips
	// U+201C/U+201D Text nodes → ASCII ". It cannot touch the ast.String nodes the
	// typographer already produced. Output is still inconsistent.
	t.Logf("single pass (typewriter+typographer): %s",
		renderWith(t, mixed, typewriter.New(), extension.Typographer))

	// Case 3: source-level strip, then typographer.
	// Apply typewriter's replacements directly to the raw markdown bytes before
	// any goldmark parsing. This is the correct two-pass approach: the intermediate
	// form is clean ASCII markdown, not HTML, so the typographer pass works on it.
	stripped := typewriter.StripBytes([]byte(mixed))
	t.Logf("after StripBytes: %s", stripped)
	t.Logf("two pass (strip→typographer): %s",
		renderWith(t, string(stripped), extension.Typographer))
}
