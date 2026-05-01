package typewriter_test

import (
	"testing"

	"github.com/client9/typewriter"
)

func TestQuotes(t *testing.T) {
	tests := []struct{ in, want string }{
		{"“Hello”", `"Hello"`},   // "Hello" (curly doubles)
		{"‘it’s’", `'it's'`},     // 'it's' (curly singles)
		{"«hello»", `<<hello>>`}, // «hello» (angle)
		{"„low”", `"low"`},       // „low" (low-9 open)
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := typewriter.Replace(tt.in)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDashes(t *testing.T) {
	tests := []struct{ in, want string }{
		{"em—dash", "em---dash"},
		{"en–dash", "en--dash"},
		{"minus−sign", "minus-sign"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := typewriter.Replace(tt.in)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEllipsis(t *testing.T) {
	got := typewriter.Replace("wait…")
	if got != "wait..." {
		t.Errorf("got %q, want %q", got, "wait...")
	}
}

func TestFractions(t *testing.T) {
	tests := []struct{ in, want string }{
		{"½", "1/2"},
		{"¼", "1/4"},
		{"¾", "3/4"},
		{"⅓", "1/3"},
		{"⅛", "1/8"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := typewriter.Replace(tt.in)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSymbols(t *testing.T) {
	tests := []struct{ in, want string }{
		{"©", "(c)"},
		{"®", "(r)"},
		{"™", "(tm)"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := typewriter.Replace(tt.in)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMath(t *testing.T) {
	tests := []struct{ in, want string }{
		{"10×", "10x"},
		{"10÷2", "10/2"},
		{"a≠b", "a!=b"},
		{"a≤b", "a<=b"},
		{"a≥b", "a>=b"},
		{"a→b", "a->b"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := typewriter.Replace(tt.in)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLigatures(t *testing.T) {
	tests := []struct{ in, want string }{
		{"ﬁle", "file"},
		{"ﬂow", "flow"},
		{"ﬀect", "ffect"},
		{"ﬃcient", "fficient"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := typewriter.Replace(tt.in)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBullets(t *testing.T) {
	tests := []struct{ in, want string }{
		{"item•one", "item*one"},
		{"note†ref", "note*ref"},
		{"note‡ref", "note**ref"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := typewriter.Replace(tt.in)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSpaces(t *testing.T) {
	// NBSP should be converted to a plain space by default.
	got := typewriter.Replace("a b")
	if got != "a b" {
		t.Errorf("default: got %q", got)
	}

	// WithoutCategory(Spaces) preserves NBSP.
	r := typewriter.New(typewriter.WithoutCategory(typewriter.Spaces))
	got = r.Replace("a b")
	if got != "a b" {
		t.Errorf("WithoutCategory(Spaces): got %q", got)
	}
}

func TestWithCategory(t *testing.T) {
	t.Run("whitelist", func(t *testing.T) {
		r := typewriter.New(typewriter.WithCategory(typewriter.Ellipsis))
		if got := r.Replace("wait…"); got != "wait..." {
			t.Errorf("ellipsis: got %q", got)
		}
		if got := r.Replace("em—dash"); got != "em—dash" {
			t.Errorf("dash should pass through: got %q", got)
		}
	})
	t.Run("all", func(t *testing.T) {
		r := typewriter.New(typewriter.WithCategory(typewriter.CategoryAll))
		if got := r.Replace("wait…"); got != "wait..." {
			t.Errorf("got %q", got)
		}
	})
}

func TestWithoutCategory(t *testing.T) {
	r := typewriter.New(typewriter.WithoutCategory(typewriter.Math))
	if got := r.Replace("10×"); got != "10×" {
		t.Errorf("× should pass through: got %q", got)
	}
	if got := r.Replace("wait…"); got != "wait..." {
		t.Errorf("ellipsis should still convert: got %q", got)
	}
}

func TestWithMapping(t *testing.T) {
	t.Run("override", func(t *testing.T) {
		r := typewriter.New(typewriter.WithMapping("—", "--"))
		if got := r.Replace("em—dash"); got != "em--dash" {
			t.Errorf("got %q", got)
		}
	})
	t.Run("exclude", func(t *testing.T) {
		r := typewriter.New(typewriter.WithMapping("×", ""))
		if got := r.Replace("10×"); got != "10×" {
			t.Errorf("got %q", got)
		}
	})
	t.Run("add_custom", func(t *testing.T) {
		r := typewriter.New(typewriter.WithMapping("°", "deg"))
		if got := r.Replace("90°"); got != "90deg" {
			t.Errorf("got %q", got)
		}
	})
}

func TestReplaceBytes(t *testing.T) {
	in := []byte("wait…")
	got := typewriter.ReplaceBytes(in)
	if string(got) != "wait..." {
		t.Errorf("got %q", got)
	}
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

func BenchmarkReplace(b *testing.B) {
	r := typewriter.New()
	s := "“Hello” — wait… © 2024"
	b.ResetTimer()
	for range b.N {
		_ = r.Replace(s)
	}
}
