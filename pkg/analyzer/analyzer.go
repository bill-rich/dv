package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/trustmaster/go-aspell"
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
	speller, err := aspell.NewSpeller(map[string]string{
		"lang": "en_US",
	})
	if err != nil {
		return nil, fmt.Errorf("Error: %s", err.Error())
	}
	defer speller.Delete()

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
		speller.AddToPersonal("ok")
		speller.AddToPersonal("ip")
		speller.AddToPersonal("ctx")
		speller.AddToPersonal("vm")
		speller.AddToPersonal("url")

		if _, ok := variable.Obj.Decl.(*ast.AssignStmt); ok {
			if len([]byte(variable.Name)) < 4 && !speller.Check(variable.Name) {
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
