# typewriter

Converts typographic ("smart") Unicode characters back to their plain ASCII typewriter
equivalents. Useful as a goldmark extension for standalone normalisation, or as a
source-level preprocessor before a smart-typography pass.

## Quick start

### As a goldmark extension

```go
import (
    "github.com/client9/typewriter"
    "github.com/yuin/goldmark"
)

md := goldmark.New(goldmark.WithExtensions(typewriter.New()))
```

### As a source preprocessor

```go
clean := typewriter.StripBytes(src)
```

## Two modes, different behaviour

| | Prose | Code spans / fenced blocks |
|---|---|---|
| `New()` goldmark extension | converted | preserved |
| `StripBytes` | converted | converted |

The extension form preserves code content because goldmark's HTML renderer reads code
spans and fenced blocks directly from the original source bytes — AST-level replacement
is not possible there. This is an architecture constraint of goldmark, not a policy
choice.

`StripBytes` operates on raw bytes before any parsing, so it normalises everything
including code content. This is often the right behaviour: typographic characters inside
code are almost always copy-paste corruption (smart quotes in a shell command, for
example) and stripping them is a fix.

## What it converts

All categories are active by default.

| Category | Examples | Result |
|----------|---------|--------|
| Quotes | curly doubles, singles, angle, low-9 | `"` `'` `<<` `>>` |
| Dashes | em dash, en dash, minus sign | `---` `--` `-` |
| Ellipsis | horizontal ellipsis | `...` |
| Fractions | vulgar fraction characters | `1/2` `1/4` `3/4` `1/3` `1/8` ... |
| Symbols | copyright, registered, trademark | `(c)` `(r)` `(tm)` |
| Math | multiply, divide, not-equal, arrows | `x` `/` `!=` `<=` `>=` `->` |
| Ligatures | fi, fl, ff, ffi, ffl | `fi` `fl` `ff` `ffi` `ffl` |
| Bullets | bullet, dagger, double dagger | `*` `*` `**` |
| Spaces | NBSP, thin, en, em, figure spaces | plain space |

## Configuration

### Enable only specific categories

`WithCategory` sets the active categories to exactly what you pass, replacing the
default:

```go
// Only convert dashes and ellipses; everything else passes through
md := goldmark.New(goldmark.WithExtensions(
    typewriter.New(typewriter.WithCategory(typewriter.Dashes | typewriter.Ellipsis)),
))
```

### Disable specific categories

`WithoutCategory` removes categories from the active set:

```go
md := goldmark.New(goldmark.WithExtensions(
    typewriter.New(typewriter.WithoutCategory(typewriter.Math)),
))
```

The two options compose left-to-right:

```go
// Everything except Math
typewriter.New(
    typewriter.WithCategory(typewriter.CategoryAll),
    typewriter.WithoutCategory(typewriter.Math),
)
```

### Preserve non-breaking spaces

Unicode spaces are converted by default. In practice they arrive via copy-paste from
Word or AI-generated text rather than deliberate authoring. To preserve them:

```go
typewriter.New(typewriter.WithoutCategory(typewriter.Spaces))
```

### Override or exclude individual mappings

```go
typewriter.New(
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
// Step 1: strip all typographic characters from the raw markdown source.
clean := typewriter.StripBytes(src)

// Step 2: render with the typographer — consistent output regardless of what
// the source contained.
md := goldmark.New(goldmark.WithExtensions(extension.Typographer))
md.Convert(clean, &buf)
```

The goldmark extension form does not work for this pipeline: goldmark's typographer is
an inline parser (runs during tokenisation) while typewriter is an AST transformer (runs
after). In a single goldmark instance the typographer always fires first, so typewriter
ends up stripping what the typographer just applied. `StripBytes` must be used for the
two-pass approach.
