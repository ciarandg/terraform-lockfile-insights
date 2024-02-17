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

	body := n.NamedChild(2)
	providerBlock := body.NamedChild(2)
	// providerName := providerBlock.NamedChild(1)
	providerContents := providerBlock.NamedChild(3)
	providerVersionStatement := providerContents.NamedChild(0)
	providerVersion := providerVersionStatement.NamedChild(1)
	fmt.Println(providerVersion.Type())
	fmt.Println(providerVersion.Content(sourceCode))
	fmt.Println(providerVersion)

	// find body block
	   // find provider blocks
}
