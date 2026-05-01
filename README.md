# typewriter

Converts typographic ("smart") Unicode characters back to their plain ASCII typewriter
equivalents. Zero dependencies.

Use it as a preprocessor to normalise mixed-source markdown before a smart-typography
pass, or as a standalone sanitiser for copy-paste corruption in prose or code.

For a goldmark extension that applies these conversions at the AST level, see
[goldmark-typewriter](https://github.com/client9/goldmark-typewriter).

## Quick start

```go
import "github.com/client9/typewriter"

// Package-level convenience — all Default categories active.
clean := typewriter.Replace(s)
clean := typewriter.ReplaceBytes(b)

// Configured instance.
r := typewriter.New(typewriter.WithoutCategory(typewriter.Math))
clean = r.Replace(s)
```

## What it converts

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
| Spaces | NBSP, thin, en, em, figure, hair spaces | plain space |

## Configuration

### Enable only specific categories

`WithCategory` sets the active categories to exactly what you pass, replacing the default:

```go
// Only convert dashes and ellipses; everything else passes through.
r := typewriter.New(typewriter.WithCategory(typewriter.Dashes | typewriter.Ellipsis))
```

### Disable specific categories

`WithoutCategory` removes categories from the active set:

```go
r := typewriter.New(typewriter.WithoutCategory(typewriter.Math))
```

Options compose left-to-right:

```go
// Everything except Math and Bullets.
r := typewriter.New(
    typewriter.WithCategory(typewriter.CategoryAll),
    typewriter.WithoutCategory(typewriter.Math | typewriter.Bullets),
)
```

### Override or exclude individual mappings

```go
r := typewriter.New(
    typewriter.WithMapping("—", "--"),  // prefer -- over --- for em dash
    typewriter.WithMapping("×", ""),    // leave × unchanged (empty = pass through)
    typewriter.WithMapping("°", "deg"), // add a mapping not in builtins
)
```

## Normalising before smart typography

Markdown arrives from multiple sources — hand-authored files, Word, AI-generated text,
web scrapers — each with different typographic conventions. A typographer pass on mixed
input produces inconsistent output: content already containing `"Hello"` passes through
unchanged while `"Hello"` gets converted.

The fix is to strip to a clean ASCII baseline first, then apply smart typography:

```go
import (
    "github.com/client9/typewriter"
    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark/extension"
)

// Step 1: strip typographic characters from the raw markdown source.
clean := typewriter.ReplaceBytes(src)

// Step 2: render with the typographer — consistent output regardless of input source.
md := goldmark.New(goldmark.WithExtensions(extension.Typographer))
md.Convert(clean, &buf)
```

The [goldmark-typewriter](https://github.com/client9/goldmark-typewriter) extension
operates at the AST level and is useful for standalone normalisation, but it cannot
participate in this two-pass pipeline: goldmark's typographer is an inline parser (runs
during tokenisation) while the AST transformer runs after, so the typographer always
fires first in a single goldmark instance.
