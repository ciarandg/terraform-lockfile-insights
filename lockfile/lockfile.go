package lockfile

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/hcl"
)

type LockfileSummary map[string]ProviderSummary

type ProviderSummary struct {
	version string
	// constraints string
}

func childByType(parent *sitter.Node, childType string) (*sitter.Node, error) {
	childCount := int(parent.NamedChildCount())
	for i := 0; i < childCount; i++ {
		child := parent.NamedChild(i)
		if child.Type() == childType {
			return child, nil
		}
	}
	return nil, fmt.Errorf("could not find child of type %s", childType)
}

func providerBlocks(bodyBlock *sitter.Node, sourceCode []byte) []*sitter.Node {
	var out []*sitter.Node
	childCount := int(bodyBlock.NamedChildCount())
	for i := 0; i < childCount; i++ {
		child := bodyBlock.NamedChild(i)
		if child.NamedChildCount() > 0 {
			identifier := child.NamedChild(0)
			if identifier.Content(sourceCode) == "provider" {
				out = append(out, child)
			}
		}
	}
	return out
}

func providerName(providerBlock *sitter.Node, sourceCode []byte) (string, error) {
	if int(providerBlock.NamedChildCount()) < 2 {
		return "", errors.New("expected at least 2 named children in provider block")
	}
	nameInQuotes := providerBlock.NamedChild(1).Content(sourceCode)
	return strings.Trim(nameInQuotes, `"`), nil
}

func providerVersion(providerBlock *sitter.Node, sourceCode []byte) (string, error) {
	if int(providerBlock.NamedChildCount()) < 4 {
		return "", errors.New("expected at least 4 named children in provider block")
	}
	blockBody := providerBlock.NamedChild(3)
	if int(blockBody.NamedChildCount()) < 1 {
		return "", errors.New("expected at least 1 named child in provider block body")
	}
	versionStatement := blockBody.NamedChild(0)
	if int(versionStatement.NamedChildCount()) < 2 {
		return "", errors.New("expected at least 2 named children in provider version statement")
	}
	version := versionStatement.NamedChild(1).Content(sourceCode)
	return strings.Trim(version, `"`), nil
}


func GenerateLockfileSummary(filePath string) (LockfileSummary, error) {
	out := LockfileSummary{}

	sourceCode, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	parser := sitter.NewParser()
	parser.SetLanguage(hcl.GetLanguage())

	tree, _ := parser.ParseCtx(context.Background(), nil, sourceCode)

	n := tree.RootNode()

	body, err := childByType(n, "body")
	if err != nil {
		return nil, err
	}

	providerBlocks := providerBlocks(body, sourceCode)
	for i := 0; i < len(providerBlocks); i++ {
		block := providerBlocks[i]

		name, err := providerName(block, sourceCode)
		if err != nil {
			return nil, err
		}

		version, err := providerVersion(block, sourceCode)
		if err != nil {
			return nil, err
		}

		out[name] = ProviderSummary{version}
	}

	return out, nil
}