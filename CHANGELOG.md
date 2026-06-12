# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2026-06-12

### Added

- **Character substitutions** via `strings.Replacer` trie — single-pass, longest-match
  replacement for all built-in categories:
  - `Quotes` — curly and angle quotation marks → straight ASCII (`"` `'` `«»` `‹›` `„` `‚`)
  - `Dashes` — em dash, en dash, figure dash, non-breaking hyphen, hyphen, minus sign → `-`/`--`/`---`
  - `Ellipsis` — `…` `⋯` → `...`
  - `Fractions` — 19 Unicode vulgar fractions → `n/d` form (½ → `1/2`, etc.)
  - `Symbols` — `©` `®` `™` `§` `¶` → `(c)` `(r)` `(tm)` `S.` `P.`
  - `Math` — `×` `÷` `≠` `≤` `≥` `≈` `±` `∞` `→` `←` `⇒` `⇐` → ASCII equivalents
  - `Ligatures` — `ﬁ` `ﬂ` `ﬀ` `ﬃ` `ﬄ` `ﬅ` `ﬆ` → component letters
  - `Bullets` — `•` `‣` `·` `․` `‥` `†` `‡` → `*` `.` `**`
  - `Spaces` — 9 non-standard Unicode spaces including `U+2028` LINE SEPARATOR and
    `U+2029` PARAGRAPH SEPARATOR → plain ASCII space
- **`Default`** and **`CategoryAll`** constants enabling all categories at once;
  combine with `|` or exclude with `&^`
- **`Overrides`** map for per-codepoint customisation: remap a built-in entry to a
  custom value, or suppress it with an empty string
- **Unicode style run detection** — contiguous runs of sans-serif bold, italic,
  bold-italic, monospace, superscript, and subscript variants decoded to plain ASCII;
  each style configurable with optional `Prefix`/`Suffix` for markdown or HTML wrapping
  (`**bold**`, `_italic_`, `<b>…</b>`, etc.)
- **`UnicodeStyle`** constants: `Bold`, `Italic`, `BoldItalic`, `Monospace`,
  `Superscript`, `Subscript`
- **`Config`** struct with `Categories`, `Overrides`, and `Runs` fields
- **`New(Config) *Replacer`** constructor; `Replacer` is safe for concurrent use
- **`Replacer.Replace(string) string`** and **`Replacer.ReplaceBytes([]byte) []byte`**
  methods; a nil `*Replacer` falls back to the `Default` replacer
- **Package-level `Replace` and `ReplaceBytes`** convenience functions backed by a
  lazily-initialised default `Replacer` (all `Default` categories, no run detection)
- Zero external dependencies; minimum Go version 1.22

