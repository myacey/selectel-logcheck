package logcheck

import (
	"go/ast"
	"log"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

// isLogFunc checks if the function name is a known logging function.
func isLogFunc(name string) bool {
	log.Println("INVALID")
	log.Println("!!!!!")
	return slices.Contains(config.LogFuncs, name)
}

// detectSensitive traverses an AST expression and reports if it contains sensitive words.
func detectSensitive(pass *analysis.Pass, expr ast.Expr) {
	switch v := expr.(type) {
	case *ast.Ident:
		if containsSensitiveWord(v.Name) {
			pass.Reportf(v.Pos(), "logs should not contain potentially sensitive data")
		}
	case *ast.BinaryExpr:
		detectSensitive(pass, v.X)
		detectSensitive(pass, v.Y)
	}
}

// containsSensitiveWord checks if a string contains any predefined sensitive term.
func containsSensitiveWord(s string) bool {
	s = strings.ToLower(s)
	for _, w := range config.SensitiveWords {
		if strings.Contains(s, w) {
			return true
		}
	}
	return false
}

// startsWithUpper reports whether the string starts with an uppercase letter.
func startsWithUpper(s string) bool {
	if s == "" {
		return false
	}
	r, _ := utf8.DecodeRuneInString(s)
	if !unicode.IsLetter(r) {
		return false
	}

	return unicode.IsUpper(r)
}

// isEnglishOnly reports whether the string contains only ASCII letters.
func isEnglishOnly(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// hasSpecialChars reports whether the string contains characters other than
// letters, digits, spaces, and the '%' format specifier.
func hasSpecialChars(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) || r == '%' {
			continue
		}
		return true
	}
	return false
}
