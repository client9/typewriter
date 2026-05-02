// Package typewriter converts typographic ("smart") Unicode characters back to
// their plain ASCII typewriter equivalents, and normalises Unicode style
// variants (bold, italic, monospace, superscript, subscript) back to plain
// letters.
//
// It is designed for cleaning up text that has passed through a word processor,
// rich-text editor, or AI-generated content pipeline: curly quotes, em dashes,
// ligatures, non-breaking spaces, styled Unicode alphabets, and similar
// characters that look fine on screen but cause problems in plain-text contexts
// such as source code, configuration files, command-line arguments, and
// Markdown documents.
//
// # Quick start
//
// The package-level functions use all [Default] conversions and require no
// configuration:
//
//	clean := typewriter.Replace(s)
//	clean := typewriter.ReplaceBytes(b)
//
// # Creating a Replacer
//
// For custom behaviour, create a [Replacer] with [New]:
//
//	r := typewriter.New(typewriter.Config{
//	    Categories: typewriter.Default,
//	})
//	clean := r.Replace(s)
//
// A [Replacer] is safe for concurrent use and should be created once and
// reused.
//
// # Categories
//
// Built-in conversions are grouped into [Category] bitfields. [Default]
// (equivalently [CategoryAll]) enables all of them. Use bitwise operations to
// select a subset:
//
//	// Only ellipsis and dashes:
//	typewriter.Config{Categories: typewriter.Ellipsis | typewriter.Dashes}
//
//	// Everything except math symbols:
//	typewriter.Config{Categories: typewriter.Default &^ typewriter.Math}
//
// The defined categories are [Quotes], [Dashes], [Ellipsis], [Fractions],
// [Symbols], [Math], [Ligatures], [Bullets], and [Spaces].
//
// # Overrides
//
// [Config.Overrides] customises or extends the built-in table on a
// character-by-character basis. Overrides are applied before built-ins.
// An empty target string excludes the source from conversion:
//
//	typewriter.Config{
//	    Categories: typewriter.Default,
//	    Overrides: map[string]string{
//	        "—": "--",  // prefer -- over the default ---
//	        "×": "",    // leave × unchanged
//	        "°": "deg", // add a mapping not in any built-in category
//	    },
//	}
//
// # Unicode style runs
//
// Social-media and AI-generated text frequently contains styled Unicode
// variants: 𝗯𝗼𝗹𝗱, 𝘪𝘵𝘢𝘭𝘪𝘤, 𝚖𝚘𝚗𝚘𝚜𝚙𝚊𝚌𝚎, ˢᵘᵖᵉʳˢᶜʳⁱᵖᵗ, ₛᵤᵦₛ꜀ᵣᵢₚₜ.
// [Config.Runs] maps contiguous runs of styled characters to plain ASCII,
// optionally wrapping the run with a configurable prefix and suffix:
//
//	r := typewriter.New(typewriter.Config{
//	    Categories: typewriter.Default,
//	    Runs: []typewriter.RunStyle{
//	        {Style: typewriter.Bold,   Prefix: "**", Suffix: "**"},
//	        {Style: typewriter.Italic, Prefix: "_",  Suffix: "_"},
//	    },
//	})
//
// With empty [RunStyle.Prefix] and [RunStyle.Suffix] (the zero value), styled
// runs are stripped to plain ASCII with no added markup. Character
// substitutions (quotes, dashes, etc.) and run detection compose correctly in
// a single pass.
package typewriter
