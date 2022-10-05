package myanalyzers

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const Doc = "exist os.Exit in main function"

// ExitAnalyzer it is forbidden to use a direct call to os.Exit in the main function of the main package
var ExitAnalyzer = &analysis.Analyzer{
	Name: "exit",
	Doc:  Doc,
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		fileName := pass.Fset.Position(file.Pos()).Filename
		if !strings.HasSuffix(fileName, ".go") {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			if f, ok := node.(*ast.FuncDecl); ok {
				if f.Name.Name == "main" {
					for _, s := range f.Body.List {
						if expr, ok := s.(*ast.ExprStmt); ok {
							if call, ok := expr.X.(*ast.CallExpr); ok {
								if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
									i := selector.X.(*ast.Ident)
									if i.Name == "os" && selector.Sel.Name == "Exit" {
										pass.Reportf(selector.Pos(), "Exit method call")
									}
								}
							}
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
