package lockfile

import (
	"context"
	"errors"
	"os"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/hcl"
)

type Lockfile struct {
	ProviderBlocks map[string]ProviderBlock
}

type ProviderBlock struct {
	version string
	contraints string
	hashes []string
}

func NewLockfile(filePath string) (Lockfile, error) {
	sourceCode, err := os.ReadFile(filePath)
	if err != nil {
		return Lockfile{}, err
	}

	bodyBlock, err := bodyBlock(sourceCode)
	if err != nil {
		return Lockfile{}, err
	}

	providerBlocks, err := providerBlocks(sourceCode, bodyBlock)
	if err != nil {
		return Lockfile{}, err
	}

	return Lockfile{providerBlocks}, nil
}

func bodyBlock(sourceCode []byte) (*sitter.Node, error) {
	parser := sitter.NewParser()
	parser.SetLanguage(hcl.GetLanguage())

	tree, _ := parser.ParseCtx(context.Background(), nil, sourceCode)

	n := tree.RootNode()

	body := childByType(n, "body")
	if body == nil {
		return nil, errors.New("failed to find body block in lockfile source code")
	}
	return body, nil
}

func providerBlocks(sourceCode []byte, bodyBlock *sitter.Node) (map[string]ProviderBlock, error) {
	out := map[string]ProviderBlock{}

	childCount := int(bodyBlock.NamedChildCount())
	for i := 0; i < childCount; i++ {
		child := bodyBlock.NamedChild(i)
		if child.NamedChildCount() > 0 {
			identifier := child.NamedChild(0)
			if identifier.Content(sourceCode) == "provider" {
				name, _ := providerName(child, sourceCode)
				version, _ := providerVersion(child, sourceCode)
				out[name] = ProviderBlock{version, "fake_constraint", []string{"fake_hash"}} // TODO fix stubs
			}
		}
	}

	return out, nil
}

func childByType(parent *sitter.Node, childType string) *sitter.Node {
	childCount := int(parent.NamedChildCount())
	for i := 0; i < childCount; i++ {
		child := parent.NamedChild(i)
		if child.Type() == childType {
			return child
		}
	}
	return nil
}

func childByTypeRec(parent *sitter.Node, childType string) *sitter.Node {
	childCount := int(parent.NamedChildCount())
	for i := 0; i < childCount; i++ {
		match := childByTypeRec(parent.NamedChild(i), childType)
		if match != nil {
			return match
		}
	}
	match := childByType(parent, childType)
	if match != nil {
		return match
	}
	return nil
}

func childByPredicate(parent *sitter.Node, predicate func (*sitter.Node) bool) *sitter.Node {
	childCount := int(parent.NamedChildCount())
	for i := 0; i < childCount; i++ {
		child := parent.NamedChild(i)
		if predicate(child) {
			return child
		}
	}
	return nil
}

func providerName(providerBlock *sitter.Node, sourceCode []byte) (string, error) {
	if int(providerBlock.NamedChildCount()) < 2 {
		return "", errors.New("expected at least 2 named children in provider block")
	}
	nameInQuotes := providerBlock.NamedChild(1).Content(sourceCode)
	return strings.Trim(nameInQuotes, `"`), nil
}

func providerVersion(providerBlock *sitter.Node, sourceCode []byte) (string, error) {
	isVersionStatement := func (block *sitter.Node) bool {
		return block.NamedChildCount() == 2 && block.NamedChild(0).Content(sourceCode) == "version"
	}

	blockBody := childByType(providerBlock, "body")
	if blockBody == nil {
		return "", errors.New("failed to find block body in provider block")
	}
	versionStatement := childByPredicate(blockBody, isVersionStatement)
	versionLiteral := childByTypeRec(versionStatement.NamedChild(1), "template_literal")
	return versionLiteral.Content(sourceCode), nil
}