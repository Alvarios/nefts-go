package utils

import (
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"regexp"
	"strings"
	"unicode"
)

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func HardFormat(str string) string {
	reNoSpaces := regexp.MustCompile(` `)
	noSpaces := reNoSpaces.ReplaceAllString(str, "_")
	toLower := strings.ToLower(noSpaces)

	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	normalized, _, _ :=transform.String(t, toLower)

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(err)
	}

	processedString := reg.ReplaceAllString(normalized, "")

	return processedString
}

func SoftFormat(str string) string {
	reNoBreaks := regexp.MustCompile(`\r?\n`)
	noBreaks := reNoBreaks.ReplaceAllString(str, " ")
	trimmed := strings.TrimSpace(noBreaks)
	reNoDoubleSpaces := regexp.MustCompile(` {2}`)
	noDoubleSpaces := reNoDoubleSpaces.ReplaceAllString(trimmed, " ")

	return noDoubleSpaces
}

func SoftFormatWithLineBreaks(str string) string {
	trimmed := strings.TrimSpace(str)
	reNoDoubleSpaces := regexp.MustCompile(` {2}`)
	noDoubleSpaces := reNoDoubleSpaces.ReplaceAllString(trimmed, " ")

	return noDoubleSpaces
}

func QuickReplace(str string, reg string, replacement string) string {
	return regexp.MustCompile(reg).ReplaceAllString(str, replacement)
}