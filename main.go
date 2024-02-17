package main

import (
	"context"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/hcl"
)

func main() {
	parser := sitter.NewParser()
	parser.SetLanguage(hcl.GetLanguage())

	sourceCode := []byte("provider \"foo\" \"bar\" {}")
	tree, _ := parser.ParseCtx(context.Background(), nil, sourceCode)

	n := tree.RootNode()

	fmt.Println(n) // (program (lexical_declaration (variable_declarator (identifier) (number))))
	
	child := n.NamedChild(0)
	fmt.Println(child.Type()) // lexical_declaration
	fmt.Println(child.StartByte()) // 0
	fmt.Println(child.EndByte()) // 9
}
