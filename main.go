package main

import (
	"context"
	"fmt"
	"os"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/hcl"
)

func main() {
	filePath := "example.lock.hcl"
	sourceCode, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}

	parser := sitter.NewParser()
	parser.SetLanguage(hcl.GetLanguage())

	tree, _ := parser.ParseCtx(context.Background(), nil, sourceCode)

	n := tree.RootNode()

	fmt.Println(n)
	
	child := n.NamedChild(0)
	fmt.Println(child.Type())
	fmt.Println(child.StartByte())
	fmt.Println(child.EndByte())
}
