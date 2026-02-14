// // Package logcheck defines an Analyzer that reports problematic log usage.
package logcheck

import (
	"go/ast"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

// Analyzer is the main entry point for the analysis tool.
var Analyzer = &analysis.Analyzer{
	Name: "logcheck",
	Doc:  "reports invalid logs",
	Run:  run,
}

var config = defaultConfig

// run inspects the AST of every file in the package.
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			// check whether the call expression matches function call
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			analyzeCall(pass, call)
			return true
		})
	}

	return nil, nil
}

// analyzeCall checks if the call expression is a log function and validates it.
func analyzeCall(pass *analysis.Pass, call *ast.CallExpr) {
	fnc, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || !isLogFunc(fnc.Sel.Name) {
		return
	}

	if len(call.Args) == 0 {
		return
	}

	// Check the first argument as a log message.
	if lit, ok := call.Args[0].(*ast.BasicLit); ok {
		validateLogMessage(pass, lit)
	}

	// Check all arguments for sensitive data.
	for _, arg := range call.Args {
		detectSensitive(pass, arg)
	}
}

// validateLogMessage checks the log message's format rules.
func validateLogMessage(pass *analysis.Pass, lit *ast.BasicLit) {
	msg, err := strconv.Unquote(lit.Value)
	if err != nil || msg == "" {
		return
	}

	if config.CheckLowercase && startsWithUpper(msg) {
		fixedMsg := strings.ToLower(string(msg[0])) + msg[1:]
		pass.Report(analysis.Diagnostic{
			Pos:     lit.Pos(),
			End:     lit.End(),
			Message: "log message should start with lowercase letter",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "convert first letter to lowercase",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     lit.Pos() + 1,
							End:     lit.End() - 1,
							NewText: []byte(fixedMsg),
						},
					},
				},
			},
		})
	}

	if config.CheckEnglish && !isEnglishOnly(msg) {
		pass.Reportf(lit.Pos(), "log message should contain only english letters")
	}

	if config.CheckSpecial && hasSpecialChars(msg) {
		fixedMsg := make([]byte, 0, len(msg))
		for _, r := range msg {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) || r == '%' {
				fixedMsg = append(fixedMsg, byte(r))
			}
		}
		safeMsg := strconv.Quote(string(fixedMsg))
		pass.Report(analysis.Diagnostic{
			Pos:     lit.Pos(),
			End:     lit.Pos(),
			Message: "log message should not contain special characters",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "remove special characters",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     lit.Pos() + 1,
							End:     lit.End() - 1,
							NewText: []byte(safeMsg[1 : len(safeMsg)-1]),
						},
					},
				},
			},
		})
	}
}
