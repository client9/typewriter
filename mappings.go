package typewriter

// mapping is a single typographic → typewriter conversion.
type mapping struct {
	cat  Category
	from string // unicode source (typically one codepoint)
	to   string // ASCII target
}

// builtinMappings is the full conversion table, grouped by category.
// Order within a category is preserved but doesn't affect correctness
// (all sources are single codepoints; no prefix-collision risk).
var builtinMappings = []mapping{

	// Quotes — curly/angle/low quotation marks → straight ASCII
	{Quotes, "“", `"`},  // " LEFT DOUBLE QUOTATION MARK
	{Quotes, "”", `"`},  // " RIGHT DOUBLE QUOTATION MARK
	{Quotes, "„", `"`},  // „ DOUBLE LOW-9 QUOTATION MARK
	{Quotes, "‘", `'`},  // ' LEFT SINGLE QUOTATION MARK
	{Quotes, "’", `'`},  // ' RIGHT SINGLE QUOTATION MARK
	{Quotes, "‚", `'`},  // ‚ SINGLE LOW-9 QUOTATION MARK
	{Quotes, "«", `<<`}, // « LEFT-POINTING DOUBLE ANGLE QUOTATION MARK
	{Quotes, "»", `>>`}, // » RIGHT-POINTING DOUBLE ANGLE QUOTATION MARK
	{Quotes, "‹", `<`},  // ‹ SINGLE LEFT-POINTING ANGLE QUOTATION MARK
	{Quotes, "›", `>`},  // › SINGLE RIGHT-POINTING ANGLE QUOTATION MARK

	// Dashes — typographic dashes and hyphens → ASCII hyphens
	// Em dash maps to --- (matching goldmark typographer's --- → —)
	// En dash maps to -- (matching -- → –)
	{Dashes, "—", `---`}, // — EM DASH
	{Dashes, "–", `--`},  // – EN DASH
	{Dashes, "‒", `-`},   // ‒ FIGURE DASH
	{Dashes, "‑", `-`},   // ‑ NON-BREAKING HYPHEN
	{Dashes, "‐", `-`},   // ‐ HYPHEN
	{Dashes, "−", `-`},   // − MINUS SIGN

	// Ellipsis
	{Ellipsis, "…", `...`}, // … HORIZONTAL ELLIPSIS
	{Ellipsis, "⋯", `...`}, // ⋯ MIDLINE HORIZONTAL ELLIPSIS

	// Fractions — Unicode vulgar fractions → n/d
	{Fractions, "½", `1/2`},  // ½
	{Fractions, "¼", `1/4`},  // ¼
	{Fractions, "¾", `3/4`},  // ¾
	{Fractions, "⅓", `1/3`},  // ⅓
	{Fractions, "⅔", `2/3`},  // ⅔
	{Fractions, "⅕", `1/5`},  // ⅕
	{Fractions, "⅖", `2/5`},  // ⅖
	{Fractions, "⅗", `3/5`},  // ⅗
	{Fractions, "⅘", `4/5`},  // ⅘
	{Fractions, "⅙", `1/6`},  // ⅙
	{Fractions, "⅚", `5/6`},  // ⅚
	{Fractions, "⅛", `1/8`},  // ⅛
	{Fractions, "⅜", `3/8`},  // ⅜
	{Fractions, "⅝", `5/8`},  // ⅝
	{Fractions, "⅞", `7/8`},  // ⅞
	{Fractions, "⅐", `1/7`},  // ⅐
	{Fractions, "⅑", `1/9`},  // ⅑
	{Fractions, "⅒", `1/10`}, // ⅒

	// Symbols — common typographic symbols → ASCII equivalents
	// Matches goldmark typographer: (c) → ©, (r) → ®, (tm) → ™
	{Symbols, "©", `(c)`},  // © COPYRIGHT SIGN
	{Symbols, "®", `(r)`},  // ® REGISTERED SIGN
	{Symbols, "™", `(tm)`}, // ™ TRADE MARK SIGN
	{Symbols, "§", `S.`},   // § SECTION SIGN
	{Symbols, "¶", `P.`},   // ¶ PILCROW SIGN

	// Math — mathematical operators → ASCII
	{Math, "×", `x`},   // × MULTIPLICATION SIGN  (e.g. 10x)
	{Math, "÷", `/`},   // ÷ DIVISION SIGN
	{Math, "≠", `!=`},  // ≠ NOT EQUAL TO
	{Math, "≤", `<=`},  // ≤ LESS-THAN OR EQUAL TO
	{Math, "≥", `>=`},  // ≥ GREATER-THAN OR EQUAL TO
	{Math, "≈", `~=`},  // ≈ ALMOST EQUAL TO
	{Math, "±", `+/-`}, // ± PLUS-MINUS SIGN
	{Math, "∞", `inf`}, // ∞ INFINITY
	{Math, "→", `->`},  // → RIGHTWARDS ARROW
	{Math, "←", `<-`},  // ← LEFTWARDS ARROW
	{Math, "⇒", `=>`},  // ⇒ RIGHTWARDS DOUBLE ARROW
	{Math, "⇐", `<==`}, // ⇐ LEFTWARDS DOUBLE ARROW (<= would collide with ≤)

	// Ligatures — typographic ligatures → component letters
	{Ligatures, "ﬃ", `ffi`}, // ﬃ LATIN SMALL LIGATURE FFI (longest first)
	{Ligatures, "ﬄ", `ffl`}, // ﬄ LATIN SMALL LIGATURE FFL
	{Ligatures, "ﬀ", `ff`},  // ﬀ LATIN SMALL LIGATURE FF
	{Ligatures, "ﬁ", `fi`},  // ﬁ LATIN SMALL LIGATURE FI
	{Ligatures, "ﬂ", `fl`},  // ﬂ LATIN SMALL LIGATURE FL
	{Ligatures, "ﬅ", `st`},  // ﬅ LATIN SMALL LIGATURE LONG S T
	{Ligatures, "ﬆ", `st`},  // ﬆ LATIN SMALL LIGATURE ST

	// Bullets — mid-paragraph markers → typewriter equivalents
	{Bullets, "•", `*`},  // • BULLET
	{Bullets, "‣", `-`},  // ‣ TRIANGULAR BULLET
	{Bullets, "·", `.`},  // · MIDDLE DOT
	{Bullets, "․", `.`},  // ․ ONE DOT LEADER
	{Bullets, "‥", `..`}, // ‥ TWO DOT LEADER
	{Bullets, "†", `*`},  // † DAGGER
	{Bullets, "‡", `**`}, // ‡ DOUBLE DAGGER

	// Spaces — normalise non-standard whitespace to ASCII space.
	// Includes Unicode Zs (space separators) and the structural separators
	// U+2028/U+2029 (Zl/Zp). The latter appear in AI-generated text and JSON
	// payloads; converting to plain space is the honest conversion.
	{Spaces, " ", ` `}, // NO-BREAK SPACE
	{Spaces, " ", ` `}, // NARROW NO-BREAK SPACE
	{Spaces, " ", ` `}, // FIGURE SPACE
	{Spaces, " ", ` `}, // EN SPACE
	{Spaces, " ", ` `}, // EM SPACE
	{Spaces, " ", ` `}, // THIN SPACE
	{Spaces, " ", ` `}, // HAIR SPACE
	{Spaces, " ", ` `}, // U+2028  LINE SEPARATOR (Zl)
	{Spaces, " ", ` `}, // U+2029  PARAGRAPH SEPARATOR (Zp)
}
