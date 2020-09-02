package analyzer

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/fatih/camelcase"
	"github.com/trustmaster/go-aspell"
	"golang.org/x/tools/go/analysis"
)

var (
	exemptions *string
	minLength  *int
)

func init() {
	exemptions = Analyzer.Flags.String("exemptions", "", "comma separated list of additional words to accept")
	minLength = Analyzer.Flags.Int("min-length", 4, "minimum variable length")
}

// Analyzer is the core.
var Analyzer = &analysis.Analyzer{
	Name: "dv",
	Doc:  "Checks that all variables meet a minimum length or include words.",
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
	for _, exemption := range strings.Split(",", *exemptions) {
		speller.AddToPersonal(exemption)
	}

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
			if len([]byte(variable.Name)) < *minLength && !checkVariableSpelling(variable.Name, speller) {
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

func checkVariableSpelling(varName string, speller aspell.Speller) bool {
	for _, word := range camelcase.Split(varName) {
		if speller.Check(word) {
			return true
		}
	}
	return false
}
