package typewriter_test

import (
	"fmt"

	"github.com/client9/typewriter"
)

func Example() {
	fmt.Println(typewriter.Replace(`"Hello" — wait…`))
	// Output:
	// "Hello" --- wait...
}

func ExampleReplace() {
	fmt.Println(typewriter.Replace("© 2024 — all rights reserved™"))
	// Output:
	// (c) 2024 --- all rights reserved(tm)
}

func ExampleReplaceBytes() {
	fmt.Println(string(typewriter.ReplaceBytes([]byte("ﬁle ½ done…"))))
	// Output:
	// file 1/2 done...
}

func ExampleNew() {
	r := typewriter.New(typewriter.Config{
		Categories: typewriter.Default,
	})
	fmt.Println(r.Replace("½ price — today only…"))
	// Output:
	// 1/2 price --- today only...
}

func ExampleNew_categoryWhitelist() {
	// Convert only ellipsis; leave everything else untouched.
	r := typewriter.New(typewriter.Config{
		Categories: typewriter.Ellipsis,
	})
	fmt.Println(r.Replace(`"wait…"`))
	// Output:
	// "wait..."
}

func ExampleNew_categoryBlacklist() {
	// Convert everything except math symbols.
	r := typewriter.New(typewriter.Config{
		Categories: typewriter.Default &^ typewriter.Math,
	})
	fmt.Println(r.Replace("10× better…"))
	// Output:
	// 10× better...
}

func ExampleNew_overrideValue() {
	// Use double-hyphen for em dash instead of the default triple-hyphen.
	r := typewriter.New(typewriter.Config{
		Categories: typewriter.Default,
		Overrides:  map[string]string{"—": "--"},
	})
	fmt.Println(r.Replace("one—two"))
	// Output:
	// one--two
}

func ExampleNew_overrideExclude() {
	// Leave the multiplication sign unchanged.
	r := typewriter.New(typewriter.Config{
		Categories: typewriter.Default,
		Overrides:  map[string]string{"×": ""},
	})
	fmt.Println(r.Replace("10×"))
	// Output:
	// 10×
}

func ExampleNew_overrideAdd() {
	// Add a mapping not covered by any built-in category.
	r := typewriter.New(typewriter.Config{
		Categories: typewriter.Default,
		Overrides:  map[string]string{"°": "deg"},
	})
	fmt.Println(r.Replace("90°"))
	// Output:
	// 90deg
}

func ExampleNew_runsStrip() {
	// Strip styled Unicode variants to plain ASCII.
	// 𝗛𝗲𝗹𝗹𝗼 = sans-serif bold "Hello" (U+1D5DB … U+1D5FC)
	// 𝘸𝘰𝘳𝘭𝘥 = sans-serif italic "world" (U+1D638 … U+1D625)
	bold := "\U0001d5db\U0001d5f2\U0001d5f9\U0001d5f9\U0001d5fc"
	italic := "\U0001d638\U0001d630\U0001d633\U0001d62d\U0001d625"

	r := typewriter.New(typewriter.Config{
		Categories: typewriter.Default,
		Runs: []typewriter.RunStyle{
			{Style: typewriter.Bold},
			{Style: typewriter.Italic},
		},
	})
	fmt.Println(r.Replace(bold + " " + italic))
	// Output:
	// Hello world
}

func ExampleNew_runsMarkdown() {
	// Convert styled Unicode variants to Markdown.
	// 𝗛𝗲𝗹𝗹𝗼 = sans-serif bold "Hello" (U+1D5DB … U+1D5FC)
	// 𝘸𝘰𝘳𝘭𝘥 = sans-serif italic "world" (U+1D638 … U+1D625)
	bold := "\U0001d5db\U0001d5f2\U0001d5f9\U0001d5f9\U0001d5fc"
	italic := "\U0001d638\U0001d630\U0001d633\U0001d62d\U0001d625"

	r := typewriter.New(typewriter.Config{
		Categories: typewriter.Default,
		Runs: []typewriter.RunStyle{
			{Style: typewriter.Bold, Prefix: "**", Suffix: "**"},
			{Style: typewriter.Italic, Prefix: "_", Suffix: "_"},
		},
	})
	fmt.Println(r.Replace(bold + " " + italic))
	// Output:
	// **Hello** _world_
}

func ExampleNew_runsHTML() {
	// Convert styled Unicode variants to HTML.
	// 𝗛𝗲𝗹𝗹𝗼 = sans-serif bold "Hello" (U+1D5DB … U+1D5FC)
	bold := "\U0001d5db\U0001d5f2\U0001d5f9\U0001d5f9\U0001d5fc"

	r := typewriter.New(typewriter.Config{
		Categories: typewriter.Default,
		Runs: []typewriter.RunStyle{
			{Style: typewriter.Bold, Prefix: "<b>", Suffix: "</b>"},
		},
	})
	fmt.Println(r.Replace(bold + " world"))
	// Output:
	// <b>Hello</b> world
}

func ExampleNew_runsSuperscript() {
	// Render superscript digits with a caret (common in plain-text math).
	r := typewriter.New(typewriter.Config{
		Categories: typewriter.Default,
		Runs:       []typewriter.RunStyle{{Style: typewriter.Superscript, Prefix: "^"}},
	})
	fmt.Println(r.Replace("E=mc²"))
	// Output:
	// E=mc^2
}

func ExampleNew_runsAndSubstitutions() {
	// Run detection and character substitutions compose in a single pass.
	// 𝗛𝗲𝗹𝗹𝗼 = sans-serif bold "Hello" (U+1D5DB … U+1D5FC)
	// 𝘄𝗼𝗿𝗹𝗱 = sans-serif bold "world" (U+1D604 … U+1D5F1)
	boldHello := "\U0001d5db\U0001d5f2\U0001d5f9\U0001d5f9\U0001d5fc"
	boldWorld := "\U0001d604\U0001d5fc\U0001d5ff\U0001d5f9\U0001d5f1"

	r := typewriter.New(typewriter.Config{
		Categories: typewriter.Default,
		Runs:       []typewriter.RunStyle{{Style: typewriter.Bold, Prefix: "**", Suffix: "**"}},
	})
	fmt.Println(r.Replace(boldHello + " © " + boldWorld))
	// Output:
	// **Hello** (c) **world**
}
