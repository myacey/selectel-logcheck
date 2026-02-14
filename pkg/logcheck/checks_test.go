package logcheck

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"
)

func TestIsLogFunc(t *testing.T) {
	testCases := []struct {
		name string
		val  string
		exp  bool
	}{
		{
			name: "print-functions",
			val:  "Print",
			exp:  true,
		},
		{
			name: "println-functions",
			val:  "Println",
			exp:  true,
		},
		{
			name: "printf-functions",
			val:  "Printf",
			exp:  true,
		},
		{
			name: "sprint-functions",
			val:  "Sprint",
			exp:  true,
		},
		{
			name: "sprintln-functions",
			val:  "Sprintln",
			exp:  true,
		},
		{
			name: "sprintf-functions",
			val:  "Sprintf",
			exp:  true,
		},
		{
			name: "info-functions",
			val:  "Info",
			exp:  true,
		},
		{
			name: "infof-functions",
			val:  "Infof",
			exp:  true,
		},
		{
			name: "warn-functions",
			val:  "Warn",
			exp:  true,
		},
		{
			name: "warnf-functions",
			val:  "Warnf",
			exp:  true,
		},
		{
			name: "error-functions",
			val:  "Error",
			exp:  true,
		},
		{
			name: "errorf-functions",
			val:  "Errorf",
			exp:  true,
		},
		{
			name: "debug-functions",
			val:  "Debug",
			exp:  true,
		},
		{
			name: "debugf-functions",
			val:  "Debugf",
			exp:  true,
		},
		{
			name: "fatal-functions",
			val:  "Fatal",
			exp:  true,
		},
		{
			name: "fatalln-functions",
			val:  "Fatalln",
			exp:  true,
		},
		{
			name: "fatalf-functions",
			val:  "Fatalf",
			exp:  true,
		},
		{
			name: "panic-functions",
			val:  "Panic",
			exp:  true,
		},
		{
			name: "panicf-functions",
			val:  "Panicf",
			exp:  true,
		},
		{
			name: "panicln-functions",
			val:  "Panicln",
			exp:  true,
		},
		{
			name: "invalid-function",
			val:  "Errorln", // нет такой функции в списке
			exp:  false,
		},
		{
			name: "wrong-case",
			val:  "print", // с маленькой буквы
			exp:  false,
		},
		{
			name: "empty-string",
			val:  "",
			exp:  false,
		},
		{
			name: "similar-name",
			val:  "Printer",
			exp:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := isLogFunc(tc.val)
			require.Equal(t, tc.exp, res)
		})
	}
}

func TestDetectSensitive(t *testing.T) {
	testCases := []struct {
		name string
		expr ast.Expr
		want int
	}{
		{
			name: "safe indent",
			expr: &ast.Ident{Name: "username"},
			want: 0,
		},
		{
			name: "sensitive ident",
			expr: &ast.Ident{Name: "password"},
			want: 1,
		},
		{
			name: "binary expr with sensitive",
			expr: &ast.BinaryExpr{
				X: &ast.BasicLit{},
				Y: &ast.Ident{Name: "token"},
			},
			want: 1,
		},
		{
			name: "binary expr nested",
			expr: &ast.BinaryExpr{
				X: &ast.BinaryExpr{
					X: &ast.BasicLit{},
					Y: &ast.Ident{Name: "password"},
				},
				Y: &ast.Ident{Name: "apiKey"},
			},
			want: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got int

			pass := &analysis.Pass{
				Fset: token.NewFileSet(),
				Report: func(data analysis.Diagnostic) {
					got++
				},
			}

			detectSensitive(pass, tc.expr)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestContainsSensitiveWord(t *testing.T) {
	testCases := []struct {
		name string
		val  string
		exp  bool
	}{
		{
			name: "contains-password",
			val:  "my password",
			exp:  true,
		},
		{
			name: "contains-pass",
			val:  "pass123",
			exp:  true,
		},
		{
			name: "contains-token",
			val:  "auth_token",
			exp:  true,
		},
		{
			name: "contains-api-key",
			val:  "api_key=123",
			exp:  true,
		},
		{
			name: "contains-apikey",
			val:  "apikey=123",
			exp:  true,
		},
		{
			name: "contains-secret",
			val:  "secretValue",
			exp:  true,
		},
		{
			name: "case-insensitive",
			val:  "My PaSsWoRd",
			exp:  true,
		},
		{
			name: "partial-word-password",
			val:  "myPassword123",
			exp:  true,
		},
		{
			name: "partial-word-token",
			val:  "tokenizer",
			exp:  true,
		},
		{
			name: "no-sensitive-words",
			val:  "hello world",
			exp:  false,
		},
		{
			name: "empty-string",
			val:  "",
			exp:  false,
		},
		{
			name: "only-spaces",
			val:  "   ",
			exp:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := containsSensitiveWord(tc.val)
			require.Equal(t, tc.exp, res)
		})
	}
}

func TestStartsWithUpper(t *testing.T) {
	testCases := []struct {
		name string
		val  string
		exp  bool
	}{
		{
			name: "starts-with-upper",
			val:  "Upper",
			exp:  true,
		},
		{
			name: "starts-with-lower",
			val:  "lower",
			exp:  false,
		},
		{
			name: "starts-with-digit",
			val:  "123abc",
			exp:  false,
		},
		{
			name: "starts-with-special",
			val:  "!hello",
			exp:  false,
		},
		{
			name: "starts-with-cyrillic-upper",
			val:  "Привет",
			exp:  true, // Привет начинается с заглавной П
		},
		{
			name: "starts-with-cyrillic-lower",
			val:  "привет",
			exp:  false,
		},
		{
			name: "single-upper-letter",
			val:  "A",
			exp:  true,
		},
		{
			name: "single-lower-letter",
			val:  "a",
			exp:  false,
		},
		{
			name: "empty-string",
			val:  "",
			exp:  false,
		},
		{
			name: "unicode-upper",
			val:  "École",
			exp:  true, // É это заглавная
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := startsWithUpper(tc.val)
			require.Equal(t, tc.exp, res)
		})
	}
}

func TestIsEnglishOnly(t *testing.T) {
	testCases := []struct {
		name string
		val  string
		exp  bool
	}{
		{
			name: "only-english-letters",
			val:  "HelloWorld",
			exp:  true,
		},
		{
			name: "with-cyrillic",
			val:  "HelloПривет",
			exp:  false,
		},
		{
			name: "with-chinese",
			val:  "Hello世界",
			exp:  false,
		},
		{
			name: "with-digits",
			val:  "Hello123",
			exp:  true,
		},
		{
			name: "with-special-chars",
			val:  "Hello!",
			exp:  true,
		},
		{
			name: "empty-string",
			val:  "",
			exp:  true,
		},
		{
			name: "only-spaces",
			val:  "   ",
			exp:  true,
		},
		{
			name: "accented-letters",
			val:  "café",
			exp:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := isEnglishOnly(tc.val)
			require.Equal(t, tc.exp, res)
		})
	}
}

func TestHasSpecialChars(t *testing.T) {
	testCases := []struct {
		name string
		val  string
		exp  bool
	}{
		{
			name: "only-letters",
			val:  "HelloWorld",
			exp:  false,
		},
		{
			name: "only-digits",
			val:  "12345",
			exp:  false,
		},
		{
			name: "only-spaces",
			val:  "   ",
			exp:  false,
		},
		{
			name: "with-percent",
			val:  "Hello %s",
			exp:  false,
		},
		{
			name: "with-exclamation",
			val:  "Hello!",
			exp:  true,
		},
		{
			name: "with-question-mark",
			val:  "Hello?",
			exp:  true,
		},
		{
			name: "with-dot",
			val:  "Hello.",
			exp:  true,
		},
		{
			name: "with-comma",
			val:  "Hello, world",
			exp:  true,
		},
		{
			name: "with-hyphen",
			val:  "hello-world",
			exp:  true,
		},
		{
			name: "with-underscore",
			val:  "hello_world",
			exp:  true,
		},
		{
			name: "with-slash",
			val:  "hello/world",
			exp:  true,
		},
		{
			name: "mixed-allowed",
			val:  "Hello 123 %",
			exp:  false,
		},
		{
			name: "empty-string",
			val:  "",
			exp:  false,
		},
		{
			name: "cyrillic-letters",
			val:  "Привет",
			exp:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := hasSpecialChars(tc.val)
			require.Equal(t, tc.exp, res)
		})
	}
}
