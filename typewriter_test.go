package typewriter_test

import (
	"testing"

	"github.com/client9/typewriter"
)

func TestQuotes(t *testing.T) {
	tests := []struct{ in, want string }{
		{"“Hello”", `"Hello"`},   // "Hello" (curly doubles)
		{"'it’s'", `'it's'`},     // 'it's' (curly singles)
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
		{"em—dash", "em---dash"},     // — EM DASH
		{"en–dash", "en--dash"},      // – EN DASH
		{"fig‒dash", "fig-dash"},     // ‒ FIGURE DASH
		{"nb‑hyphen", "nb-hyphen"},   // ‑ NON-BREAKING HYPHEN
		{"a‐b", "a-b"},               // ‐ HYPHEN
		{"minus−sign", "minus-sign"}, // − MINUS SIGN
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
	// All 9 Spaces entries should convert to a plain ASCII space.
	tests := []struct {
		name string
		in   string
	}{
		{"nbsp", "a b"},        // NO-BREAK SPACE
		{"narrow_nbsp", "a b"}, // NARROW NO-BREAK SPACE
		{"figure", "a b"},      // FIGURE SPACE
		{"en", "a b"},          // EN SPACE
		{"em", "a b"},          // EM SPACE
		{"thin", "a b"},        // THIN SPACE
		{"hair", "a b"},        // HAIR SPACE
		{"line_sep", "a b"},    // LINE SEPARATOR
		{"para_sep", "a b"},    // PARAGRAPH SEPARATOR
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := typewriter.Replace(tt.in)
			if got != "a b" {
				t.Errorf("got %q, want %q", got, "a b")
			}
		})
	}

	// Excluding Spaces preserves the characters.
	t.Run("opt_out", func(t *testing.T) {
		r := typewriter.New(typewriter.Config{
			Categories: typewriter.Default &^ typewriter.Spaces,
		})
		got := r.Replace("a b")
		if got != "a b" {
			t.Errorf("without Spaces: got %q", got)
		}
	})
}

func TestCategoryWhitelist(t *testing.T) {
	r := typewriter.New(typewriter.Config{Categories: typewriter.Ellipsis})
	if got := r.Replace("wait…"); got != "wait..." {
		t.Errorf("ellipsis: got %q", got)
	}
	if got := r.Replace("em—dash"); got != "em—dash" {
		t.Errorf("dash should pass through: got %q", got)
	}
}

func TestCategoryExclude(t *testing.T) {
	r := typewriter.New(typewriter.Config{
		Categories: typewriter.Default &^ typewriter.Math,
	})
	if got := r.Replace("10×"); got != "10×" {
		t.Errorf("× should pass through: got %q", got)
	}
	if got := r.Replace("wait…"); got != "wait..." {
		t.Errorf("ellipsis should still convert: got %q", got)
	}
}

func TestOverrides(t *testing.T) {
	t.Run("override", func(t *testing.T) {
		r := typewriter.New(typewriter.Config{
			Categories: typewriter.Default,
			Overrides:  map[string]string{"—": "--"},
		})
		if got := r.Replace("em—dash"); got != "em--dash" {
			t.Errorf("got %q", got)
		}
	})
	t.Run("exclude", func(t *testing.T) {
		r := typewriter.New(typewriter.Config{
			Categories: typewriter.Default,
			Overrides:  map[string]string{"×": ""},
		})
		if got := r.Replace("10×"); got != "10×" {
			t.Errorf("got %q", got)
		}
	})
	t.Run("add_custom", func(t *testing.T) {
		r := typewriter.New(typewriter.Config{
			Categories: typewriter.Default,
			Overrides:  map[string]string{"°": "deg"},
		})
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

// boldHello = 𝗛𝗲𝗹𝗹𝗼  (sans-serif bold)
const boldHello = "\U0001d5db\U0001d5f2\U0001d5f9\U0001d5f9\U0001d5fc"

// italicWorld = 𝘸𝘰𝘳𝘭𝘥  (sans-serif italic)
const italicWorld = "\U0001d638\U0001d630\U0001d633\U0001d62d\U0001d625"

// boldWorld = "world" in sans-serif bold (U+1D604..U+1D5F1); may render italic in some fonts
const boldWorld = "\U0001d604\U0001d5fc\U0001d5ff\U0001d5f9\U0001d5f1"

// boldItalicHi = 𝙃𝙞  (sans-serif bold-italic, U+1D643 = sans-serif bold-italic 'H')
const boldItalicHi = "\U0001d643\U0001d666"

func TestRunsBoldStrip(t *testing.T) {
	// No prefix/suffix: strip to plain ASCII.
	r := typewriter.New(typewriter.Config{
		Categories: typewriter.Default,
		Runs:       []typewriter.RunStyle{{Style: typewriter.Bold}},
	})
	if got := r.Replace(boldHello); got != "Hello" {
		t.Errorf("got %q", got)
	}
}

func TestRunsBoldMarkdown(t *testing.T) {
	r := typewriter.New(typewriter.Config{
		Categories: typewriter.Default,
		Runs:       []typewriter.RunStyle{{Style: typewriter.Bold, Prefix: "**", Suffix: "**"}},
	})
	if got := r.Replace(boldHello + " world"); got != "**Hello** world" {
		t.Errorf("got %q", got)
	}
}

func TestRunsMultipleStyles(t *testing.T) {
	r := typewriter.New(typewriter.Config{
		Categories: typewriter.Default,
		Runs: []typewriter.RunStyle{
			{Style: typewriter.Bold, Prefix: "**", Suffix: "**"},
			{Style: typewriter.Italic, Prefix: "_", Suffix: "_"},
		},
	})
	if got := r.Replace(boldHello + " " + italicWorld); got != "**Hello** _world_" {
		t.Errorf("got %q", got)
	}
}

func TestRunsSuperscript(t *testing.T) {
	r := typewriter.New(typewriter.Config{
		Categories: typewriter.Default,
		Runs:       []typewriter.RunStyle{{Style: typewriter.Superscript, Prefix: "^"}},
	})
	got := r.Replace("E=mc²")
	if got != "E=mc^2" {
		t.Errorf("got %q", got)
	}
}

func TestRunsSubscript(t *testing.T) {
	r := typewriter.New(typewriter.Config{
		Categories: typewriter.Default,
		Runs:       []typewriter.RunStyle{{Style: typewriter.Subscript}},
	})
	got := r.Replace("H₂O")
	if got != "H2O" {
		t.Errorf("got %q", got)
	}
}

func TestRunsUnconfiguredStylePassthrough(t *testing.T) {
	// BoldItalic is not in Runs; those runes should survive unchanged.
	r := typewriter.New(typewriter.Config{
		Categories: typewriter.Default,
		Runs:       []typewriter.RunStyle{{Style: typewriter.Bold, Prefix: "**", Suffix: "**"}},
	})
	got := r.Replace(boldItalicHi)
	if got != boldItalicHi {
		t.Errorf("got %q, want original %q", got, boldItalicHi)
	}
}

func TestRunsEmptyAndASCII(t *testing.T) {
	r := typewriter.New(typewriter.Config{
		Categories: typewriter.Default,
		Runs:       []typewriter.RunStyle{{Style: typewriter.Bold, Prefix: "**", Suffix: "**"}},
	})
	if got := r.Replace(""); got != "" {
		t.Errorf("empty: got %q", got)
	}
	if got := r.Replace("hello world"); got != "hello world" {
		t.Errorf("plain ASCII: got %q", got)
	}
}

func TestRunsInterleaved(t *testing.T) {
	// Char substitutions (©) and run detection in the same string.
	r := typewriter.New(typewriter.Config{
		Categories: typewriter.Default,
		Runs:       []typewriter.RunStyle{{Style: typewriter.Bold, Prefix: "**", Suffix: "**"}},
	})
	got := r.Replace(boldHello + " © " + boldWorld)
	if got != "**Hello** (c) **world**" {
		t.Errorf("got %q", got)
	}
}

func BenchmarkNew(b *testing.B) {
	b.Run("no_runs", func(b *testing.B) {
		cfg := typewriter.Config{Categories: typewriter.Default}
		for range b.N {
			_ = typewriter.New(cfg)
		}
	})
	b.Run("with_runs", func(b *testing.B) {
		cfg := typewriter.Config{
			Categories: typewriter.Default,
			Runs:       []typewriter.RunStyle{{Style: typewriter.Bold, Prefix: "**", Suffix: "**"}},
		}
		for range b.N {
			_ = typewriter.New(cfg)
		}
	})
}

func BenchmarkReplace(b *testing.B) {
	r := typewriter.New(typewriter.Config{Categories: typewriter.Default})
	s := "“Hello” — wait… © 2024"
	b.ResetTimer()
	for range b.N {
		_ = r.Replace(s)
	}
}

func BenchmarkReplaceWithRuns(b *testing.B) {
	r := typewriter.New(typewriter.Config{
		Categories: typewriter.Default,
		Runs:       []typewriter.RunStyle{{Style: typewriter.Bold, Prefix: "**", Suffix: "**"}},
	})
	s := boldHello + " — wait…"
	b.ResetTimer()
	for range b.N {
		_ = r.Replace(s)
	}
}
