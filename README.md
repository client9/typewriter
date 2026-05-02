# typewriter
[![Go Reference](https://pkg.go.dev/badge/github.com/client9/typewriter.svg)](https://pkg.go.dev/github.com/client9/typewriter)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://github.com/client9/typewriter/actions/workflows/go.yml/badge.svg)](https://github.com/client9/typewriter/actions)

Converts typographic ("smart") Unicode characters back to their plain ASCII equivalents,
and normalises Unicode style variants (bold, italic, monospace, superscript, subscript)
to plain text — optionally wrapping runs with configurable markup.

Zero dependencies. For a goldmark extension see
[goldmark-typewriter](https://github.com/client9/goldmark-typewriter).

## Quick start

```go
import "github.com/client9/typewriter"

// Package-level convenience — all Default categories.
clean := typewriter.Replace(s)
clean := typewriter.ReplaceBytes(b)

// Configured instance.
r := typewriter.New(typewriter.Config{
    Categories: typewriter.Default,
    Runs: []typewriter.RunStyle{
        {Style: typewriter.Bold, Prefix: "**", Suffix: "**"},
    },
})
clean = r.Replace(s)
```

## What it converts

### Character substitutions

All categories are active by default.

| Category | Examples | Result |
|----------|---------|--------|
| Quotes | `"` `"` `'` `'` `«` `»` `„` | `"` `'` `<<` `>>` |
| Dashes | em dash `—`, en dash `–`, minus `−` | `---` `--` `-` |
| Ellipsis | `…` | `...` |
| Fractions | `½` `¼` `¾` `⅓` `⅛` | `1/2` `1/4` `3/4` `1/3` `1/8` |
| Symbols | `©` `®` `™` | `(c)` `(r)` `(tm)` |
| Math | `×` `÷` `≠` `≤` `≥` `→` | `x` `/` `!=` `<=` `>=` `->` |
| Ligatures | `ﬁ` `ﬂ` `ﬀ` `ﬃ` | `fi` `fl` `ff` `ffi` |
| Bullets | `•` `†` `‡` | `*` `*` `**` |
| Spaces | NBSP, thin, en, em, figure, hair, U+2028, U+2029 | plain space |

### Unicode style variants (run-based)

Runs of styled characters are detected and converted together, so the whole run can be
wrapped with a prefix and suffix.

| Style | Example | Default (strip) | Markdown | HTML |
|-------|---------|-----------------|----------|------|
| `Bold` | `𝗛𝗲𝗹𝗹𝗼` | `Hello` | `**Hello**` | `<b>Hello</b>` |
| `Italic` | `𝘸𝘰𝘳𝘭𝘥` | `world` | `_world_` | `<i>world</i>` |
| `BoldItalic` | `𝙃𝙚𝙡𝙡𝙤` | `Hello` | `***Hello***` | |
| `Monospace` | `𝙷𝚎𝚕𝚕𝚘` | `Hello` | `` `Hello` `` | |
| `Superscript` | `E=mc²` | `E=mc2` | `E=mc^2` | |
| `Subscript` | `H₂O` | `H2O` | | |

Style variants are not active by default — configure with `Config.Runs`.

## Configuration

### Config struct

```go
type Config struct {
    Categories Category
    Overrides  map[string]string // from → to; empty to = pass through unchanged
    Runs       []RunStyle
}

type RunStyle struct {
    Style  UnicodeStyle
    Prefix string
    Suffix string
}
```

### Enable only specific categories

`Categories` is a bitfield — set it to exactly the categories you want:

```go
// Only convert dashes and ellipses.
r := typewriter.New(typewriter.Config{
    Categories: typewriter.Dashes | typewriter.Ellipsis,
})
```

### Disable specific categories

Use bit-clear to remove from the default set:

```go
r := typewriter.New(typewriter.Config{
    Categories: typewriter.Default &^ typewriter.Math,
})
```

### Override or exclude individual characters

```go
r := typewriter.New(typewriter.Config{
    Categories: typewriter.Default,
    Overrides: map[string]string{
        "—":   "--",   // prefer -- over --- for em dash
        "×":   "",     // leave × unchanged (empty = pass through)
        "°":   "deg",  // add a mapping not in builtins
    },
})
```

### Convert Unicode bold/italic to markdown

```go
r := typewriter.New(typewriter.Config{
    Categories: typewriter.Default,
    Runs: []typewriter.RunStyle{
        {Style: typewriter.Bold,   Prefix: "**", Suffix: "**"},
        {Style: typewriter.Italic, Prefix: "_",  Suffix: "_"},
    },
})
r.Replace("𝗛𝗲𝗹𝗹𝗼 𝘸𝘰𝘳𝘭𝘥")  // → "**Hello** _world_"
```

### Convert Unicode bold/italic to HTML

```go
r := typewriter.New(typewriter.Config{
    Categories: typewriter.Default,
    Runs: []typewriter.RunStyle{
        {Style: typewriter.Bold,   Prefix: "<b>",  Suffix: "</b>"},
        {Style: typewriter.Italic, Prefix: "<i>",  Suffix: "</i>"},
    },
})
```

### Superscripts and subscripts

```go
r := typewriter.New(typewriter.Config{
    Categories: typewriter.Default,
    Runs: []typewriter.RunStyle{
        {Style: typewriter.Superscript, Prefix: "^"},   // E=mc² → E=mc^2
        {Style: typewriter.Subscript},                  // H₂O  → H2O
    },
})
```

## Normalising before smart typography

Markdown from mixed sources (hand-authored, Word, AI-generated) arrives with
inconsistent typography. A typographer pass on mixed input produces inconsistent output:
content already containing `"Hello"` passes through unchanged while `"Hello"` gets
converted.

The fix is to strip to a clean ASCII baseline first:

```go
import (
    "github.com/client9/typewriter"
    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark/extension"
)

clean := typewriter.ReplaceBytes(src)

md := goldmark.New(goldmark.WithExtensions(extension.Typographer))
md.Convert(clean, &buf)
```

For goldmark integration (converting typographic characters inside an AST walk) see
[goldmark-typewriter](https://github.com/client9/goldmark-typewriter).

## License

[MIT](/LICENSE.md)

