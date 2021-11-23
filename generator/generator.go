package generator

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/z7zmey/php-parser/node/name"

	"github.com/z7zmey/php-parser/node"
	"github.com/z7zmey/php-parser/node/expr/binary"
	"github.com/z7zmey/php-parser/node/scalar"
	"github.com/z7zmey/php-parser/node/stmt"
	"github.com/z7zmey/php-parser/php7"
	"github.com/z7zmey/php-parser/visitor"
	"github.com/z7zmey/php-parser/walker"
)

type Generator struct {
	Buffers []string
	Imports map[string]struct{}
}

func Generate(src []byte, version string, debug bool) error {
	parser := php7.NewParser(src, version)
	parser.Parse()

	for _, e := range parser.GetErrors() {
		return errors.New(e.String())
	}

	rootNode := parser.GetRootNode()
	if debug {
		v := visitor.Dumper{
			Writer: os.Stdout,
			Indent: "",
		}
		rootNode.Walk(&v)
	}

	gen := NewGenerator()
	rootNode.Walk(gen)
	fmt.Println(gen.Buffer())

	return nil
}

func NewGenerator() *Generator {
	return &Generator{
		Imports: map[string]struct{}{},
	}
}

func (d *Generator) EnterNode(w walker.Walkable) bool {
	switch n := w.(type) {
	case *node.Root:
	case *stmt.Echo:
		d.Imports["fmt"] = struct{}{}
		d.pushBuffer("fmt.Println(")
		n.Exprs[0].Walk(d)
		d.pushBuffer(")\n")
		return false
	case *scalar.Lnumber:
		d.pushBuffer(n.Value)
		return false
	case *scalar.String:
		d.pushBuffer(n.Value)
		return false
	case *binary.Plus:
		n.Left.Walk(d)
		d.pushBuffer("+")
		n.Right.Walk(d)
		return false
	case *binary.Mul:
		n.Left.Walk(d)
		d.pushBuffer("*")
		n.Right.Walk(d)
		return false
	case *binary.Concat:
		n.Left.Walk(d)
		d.pushBuffer("+")
		n.Right.Walk(d)
		return false
	case *name.NamePart:
		if n.Value == "PHP_EOL" {
			d.pushBuffer(`"\n"`)
		}
		return false
	}
	return true
}

func (d *Generator) Buffer() string {
	importStmt := ""
	if len(d.Imports) == 1 {
		for k, _ := range d.Imports {
			importStmt = fmt.Sprintf("import \"%s\"\n\n", k)
		}
	} else if len(d.Imports) > 1 {
		importStmt += fmt.Sprintf("import (\n")
		for k, _ := range d.Imports {
			importStmt += fmt.Sprintf("\"%s\"\n", k)
		}
		importStmt += fmt.Sprintf(")")
	}
	buffers := d.Buffers
	buffers = append([]string{
		"package hoge\n\n",
		importStmt,
		"func main() {\n",
	}, d.Buffers...)

	return strings.Join(buffers, "")
}

func (d *Generator) addImport(i string) {
	d.Imports[i] = struct{}{}
}

func (d *Generator) pushBuffer(b string) {
	d.Buffers = append(d.Buffers, b)
}

func (d *Generator) LeaveNode(w walker.Walkable) {
	switch w.(type) {
	case *node.Root:
		d.pushBuffer("}")
	}
}

func (d *Generator) EnterChildNode(key string, w walker.Walkable) {}
func (d *Generator) LeaveChildNode(key string, w walker.Walkable) {}
func (d *Generator) EnterChildList(key string, w walker.Walkable) {}
func (d *Generator) LeaveChildList(key string, w walker.Walkable) {}
