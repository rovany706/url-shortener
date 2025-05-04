package analyzers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// OsExitAnalyzer анализатор использования os.Exit() в функции main().
var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check for os.Exit() in main()",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		// package main
		if pass.Pkg.Name() != "main" {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			// func main()
			if f, ok := node.(*ast.FuncDecl); ok {
				if f.Name.Name != "main" {
					return false
				}
			}

			// os.Exit() call
			if c, ok := node.(*ast.CallExpr); ok {
				if s, ok := c.Fun.(*ast.SelectorExpr); ok {
					if i, ok := s.X.(*ast.Ident); ok {
						if i.Name == "os" && s.Sel.Name == "Exit" {
							pass.Reportf(c.Pos(), `os.Exit is not allowed`)
						}
					}
				}
			}

			return true
		})
	}

	return nil, nil
}
