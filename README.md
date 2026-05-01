# typewriter

A [goldmark](https://github.com/yuin/goldmark) extension that converts typographic
("smart") Unicode characters back to their plain ASCII typewriter equivalents.

It is the complement of goldmark's built-in
[typographer](https://github.com/yuin/goldmark#typographer) extension and tools like
Smarty Pants. Content inside code spans and fenced code blocks is left untouched.

## Quick start

```go
import (
    "github.com/client9/typewriter"
    "github.com/yuin/goldmark"
)

md := goldmark.New(goldmark.WithExtensions(typewriter.New()))
```

## What it converts

All categories are active by default.

| Category | Examples | Result |
|----------|---------|--------|
| Quotes | `"hello"` `'it's'` `<<hi>>` `"low"` | `"hello"` `'it's'` `<<hi>>` `"low"` |
| Dashes | em dash, en dash, minus sign | `---` `--` `-` |
| Ellipsis | horizontal ellipsis | `...` |
| Fractions | vulgar fraction characters | `1/2` `1/4` `3/4` `1/3` `1/8` ... |
| Symbols | copyright, registered, trademark | `(c)` `(r)` `(tm)` |
| Math | multiply, divide, not-equal, arrows | `x` `/` `!=` `<=` `>=` `->` |
| Ligatures | fi, fl, ff, ffi, ffl | `fi` `fl` `ff` `ffi` `ffl` |
| Bullets | bullet, dagger, double dagger | `*` `*` `**` |
| Spaces* | NBSP, thin, en, em, figure spaces | plain space |

\* Spaces are **opt-in** — see `WithSpaces()` below.

## Configuration

### Enable only specific categories

`WithCategory` sets the active categories to exactly what you pass, replacing the
default. Use it when you want a strict whitelist:

```go
// Only convert dashes and ellipses; everything else passes through
md := goldmark.New(goldmark.WithExtensions(
    typewriter.New(typewriter.WithCategory(typewriter.Dashes | typewriter.Ellipsis)),
))
```

### Enable all categories including spaces

```go
md := goldmark.New(goldmark.WithExtensions(
    typewriter.New(typewriter.WithCategory(typewriter.CategoryAll)),
))
```

### Disable specific categories

`WithoutCategory` removes categories from the active set (default: `Default`):

```go
md := goldmark.New(goldmark.WithExtensions(
    typewriter.New(typewriter.WithoutCategory(typewriter.Math)),
))
```

Multiple categories can be combined with `|`:

```go
typewriter.WithoutCategory(typewriter.Math | typewriter.Bullets)
```

The two options compose left-to-right — set a base, then subtract:

```go
// Everything including Spaces, except Math
typewriter.New(
    typewriter.WithCategory(typewriter.CategoryAll),
    typewriter.WithoutCategory(typewriter.Math),
)
```

### Enable space normalisation

Every other category is a pure representation change: `"Hello"` and `"Hello"` mean the
same thing; `—` and `---` mean the same thing. The converted document is semantically
identical to the original.

Non-breaking space (U+00A0) is different. It is not a stylistic variant of a regular
space — it carries an instruction: *do not break the line here*. An author who writes
`100 km` with a non-breaking space between the number and unit is telling the renderer
to keep them together. Converting it to a plain space silently removes that intent and
the text may reflow differently. The other Unicode spaces (thin, en, em, figure) also
encode specific widths that a plain space does not.

Because space conversion changes what the document *does* rather than how it is encoded,
it is opt-in:

```go
md := goldmark.New(goldmark.WithExtensions(
    typewriter.New(typewriter.WithSpaces()),
))
```

### Override or exclude individual mappings

```go
md := goldmark.New(goldmark.WithExtensions(
    typewriter.New(
        typewriter.WithMapping("—", "--"),  // prefer -- over --- for em dash
        typewriter.WithMapping("x", ""),    // leave x unchanged (empty = pass through)
        typewriter.WithMapping("°", "deg"), // add a mapping not in builtins
    ),
))
```

## How it works

The extension registers a goldmark AST transformer that runs after parsing. It walks
`ast.Text` nodes and applies `bytes.ReplaceAll` for each active mapping. Nodes of type
`CodeBlock`, `FencedCodeBlock`, `CodeSpan`, `HTMLBlock`, and `RawHTML` are skipped
entirely. When a text node is modified, it is replaced with an `ast.String` node
containing the rewritten bytes.

Replacement pairs are sorted longest-source-first so that multi-byte ligatures (`ﬃ`
= 3 bytes of UTF-8) are matched before shorter overlapping prefixes.

## Relationship to goldmark typographer

goldmark's typographer extension converts ASCII sequences to typographic Unicode:

```
"Hello" -> "Hello"
---     -> —
...     -> …
```

typewriter does the reverse — useful when you receive Markdown from an editor that
inserted smart quotes or other typographic characters and you want predictable ASCII
output for downstream tooling or diffing.
