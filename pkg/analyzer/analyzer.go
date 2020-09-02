package analyzer

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"strings"
)

func main() {
	v := visitor{fset: token.NewFileSet()}
	for _, filePath := range os.Args[1:] {
		if filePath == "--" { // to be able to run this like "go run main.go -- input.go"
			continue
		}

		f, err := parser.ParseFile(v.fset, filePath, nil, 0)
		if err != nil {
			log.Fatalf("Failed to parse file %s: %s", filePath, err)
		}

		ast.Walk(&v, f)
	}
}

type visitor struct {
	fset *token.FileSet
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := func(node ast.Node) bool {
		if node == nil {
			return true
		}

		var buf bytes.Buffer
		printer.Fprint(&buf, v.fset, node)
		variable, ok := node.(*ast.Ident)
		if !ok {
			return true
		}

		if variable.Obj == nil {
			return true
		}

		if _, ok := variable.Obj.Decl.(*ast.AssignStmt); ok {
			if len([]byte(variable.Name)) < 3 && variable.Name != "ok" && variable.Name != "_" {
				fmt.Printf("%s: short variable name '%s' should be more descriptive\n",
					v.fset.Position(node.Pos()), variable.Name)
			}
		}

		for _, f := range pass.Files {
			ast.Inspect(f, inspect)
		}
	}
	return nil, nil
}
