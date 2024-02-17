package main

import (
	"context"
	"fmt"
	"os"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/hcl"
)

func childByType(parent *sitter.Node, childType string) (*sitter.Node, error) {
	childCount := int(parent.ChildCount())
	for i := 0; i < childCount; i++ {
		child := parent.NamedChild(i)
		if child.Type() == childType {
			return child, nil
		}
	}
	return nil, fmt.Errorf("could not find child of type %s", childType)
}

func main() {
	filePath := "example.lock.hcl"
	sourceCode, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	parser := sitter.NewParser()
	parser.SetLanguage(hcl.GetLanguage())

	tree, _ := parser.ParseCtx(context.Background(), nil, sourceCode)

	n := tree.RootNode()

	// body := n.NamedChild(2)
	// providerBlock := body.NamedChild(2)
	// // providerName := providerBlock.NamedChild(1)
	// providerContents := providerBlock.NamedChild(3)
	// providerVersionStatement := providerContents.NamedChild(0)
	// providerVersion := providerVersionStatement.NamedChild(1)
	// fmt.Println(providerVersion.Type())
	// fmt.Println(providerVersion.Content(sourceCode))
	// fmt.Println(providerVersion)

	// find body block
	body, err := childByType(n, "body")
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	fmt.Println(body.Type())
	fmt.Println(body)
	   // find provider blocks
}
