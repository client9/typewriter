package typewriter_test

import (
	"bytes"
	"testing"

	"github.com/client9/typewriter"
	"github.com/yuin/goldmark"
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

func TestSpacesOptIn(t *testing.T) {
	// NBSP should NOT be converted without WithSpaces.
	got := render(t, typewriter.New(), "a b")
	if got != "<p>a b</p>\n" {
		t.Errorf("default: NBSP should pass through, got %q", got)
	}

	// With WithSpaces it should become a regular space.
	got = render(t, typewriter.New(typewriter.WithSpaces()), "a b")
	if got != "<p>a b</p>\n" {
		t.Errorf("WithSpaces: got %q", got)
	}
}

func TestSkipsCodeSpan(t *testing.T) {
	ext := typewriter.New()
	got := render(t, ext, "outside … but `inside … code`")
	want := "<p>outside ... but <code>inside … code</code></p>\n"
	if got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}

func TestSkipsFencedCodeBlock(t *testing.T) {
	ext := typewriter.New()
	got := render(t, ext, "```\n… — \"\n```")
	want := "<pre><code>… — &quot;\n</code></pre>\n"
	if got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
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
