package analyzer

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

type visitor struct {
	fset *token.FileSet
}

var Analyzer = &analysis.Analyzer{
	Name: "goprintffuncname",
	Doc:  "Checks that printf-like functions are named with `f` at the end.",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := func(node ast.Node) bool {
		if node == nil {
			return true
		}

		variable, ok := node.(*ast.Ident)
		if !ok {
			return true
		}

		if variable.Obj == nil {
			return true
		}

		if _, ok := variable.Obj.Decl.(*ast.AssignStmt); ok {
			if len([]byte(variable.Name)) < 3 && variable.Name != "ok" && variable.Name != "_" && variable.Name != "ip" {
				pass.Reportf(node.Pos(), "short variable name '%s' should be more descriptive",
					variable.Name)
			}
		}
		return true
	}
	for _, f := range pass.Files {
		ast.Inspect(f, inspect)
	}
	return nil, nil
}
