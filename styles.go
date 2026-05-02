package typewriter

// styledRune is the internal lookup entry for a styled Unicode rune.
type styledRune struct {
	ascii rune
	style UnicodeStyle
}

// styleRange describes a contiguous block of styled Unicode characters
// mapping to the corresponding ASCII characters starting at asciiStart.
type styleRange struct {
	style        UnicodeStyle
	unicodeStart rune
	asciiStart   rune
	count        int
}

// styleRanges covers the sans-serif mathematical variants, which are the most
// common in copy-pasted social-media and AI-generated content.
var styleRanges = []styleRange{
	// Bold sans-serif
	{Bold, 0x1D5D4, 'A', 26}, // 𝗔–𝗭
	{Bold, 0x1D5EE, 'a', 26}, // 𝗮–𝘇
	{Bold, 0x1D7EC, '0', 10}, // 𝟬–𝟵

	// Italic sans-serif
	{Italic, 0x1D608, 'A', 26}, // 𝘈–𝘡
	{Italic, 0x1D622, 'a', 26}, // 𝘢–𝘻

	// Bold Italic sans-serif
	{BoldItalic, 0x1D63C, 'A', 26}, // 𝘼–𝙕
	{BoldItalic, 0x1D656, 'a', 26}, // 𝙖–𝙯

	// Monospace
	{Monospace, 0x1D670, 'A', 26}, // 𝙰–𝚉
	{Monospace, 0x1D68A, 'a', 26}, // 𝚊–𝚣
	{Monospace, 0x1D7F6, '0', 10}, // 𝟶–𝟿

	// Subscript digits (contiguous)
	{Subscript, 0x2080, '0', 10}, // ₀–₉
}

// superscriptRunes maps superscript codepoints to ASCII.
// Listed explicitly because the codepoints are non-contiguous.
var superscriptRunes = [...]struct {
	unicode rune
	ascii   rune
}{
	{0x2070, '0'}, // ⁰
	{0x00B9, '1'}, // ¹
	{0x00B2, '2'}, // ²
	{0x00B3, '3'}, // ³
	{0x2074, '4'}, // ⁴
	{0x2075, '5'}, // ⁵
	{0x2076, '6'}, // ⁶
	{0x2077, '7'}, // ⁷
	{0x2078, '8'}, // ⁸
	{0x2079, '9'}, // ⁹
	{0x207F, 'n'}, // ⁿ
	{0x2071, 'i'}, // ⁱ
}

// allStyledRunes is the complete lookup for all defined Unicode styles,
// built once at package initialisation.
var allStyledRunes = func() map[rune]styledRune {
	total := len(superscriptRunes)
	for _, sr := range styleRanges {
		total += sr.count
	}
	m := make(map[rune]styledRune, total)
	for _, sr := range styleRanges {
		for i := range sr.count {
			m[sr.unicodeStart+rune(i)] = styledRune{
				ascii: sr.asciiStart + rune(i),
				style: sr.style,
			}
		}
	}
	for _, sr := range superscriptRunes {
		m[sr.unicode] = styledRune{ascii: sr.ascii, style: Superscript}
	}
	return m
}()

// buildStyleLookup returns a lookup map filtered to only the styles present in
// runs. Returns nil when runs is empty (no run detection needed).
func buildStyleLookup(runs []RunStyle) map[rune]styledRune {
	if len(runs) == 0 {
		return nil
	}
	active := make(map[UnicodeStyle]struct{}, len(runs))
	for _, rs := range runs {
		active[rs.Style] = struct{}{}
	}
	m := make(map[rune]styledRune, len(allStyledRunes))
	for r, sr := range allStyledRunes {
		if _, ok := active[sr.style]; ok {
			m[r] = sr
		}
	}
	return m
}
