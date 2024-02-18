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
	contraints string // optional, will be an empty string if not present
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
				constraints, _ := providerConstraints(child, sourceCode)
				hashes, _ := providerHashes(child, sourceCode)
				out[name] = ProviderBlock{version, constraints, hashes}
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

func childrenByType(parent *sitter.Node, childType string) []*sitter.Node {
	out := []*sitter.Node{}
	childCount := int(parent.NamedChildCount())
	for i := 0; i < childCount; i++ {
		child := parent.NamedChild(i)
		if child.Type() == childType {
			out = append(out, child)
		}
	}
	return out
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

func providerConstraints(providerBlock *sitter.Node, sourceCode []byte) (string, error) {
	isConstraintsStatement := func (block *sitter.Node) bool {
		return block.NamedChildCount() == 2 && block.NamedChild(0).Content(sourceCode) == "constraints"
	}

	blockBody := childByType(providerBlock, "body")
	if blockBody == nil {
		return "", errors.New("failed to find block body in provider block")
	}
	constraintsStatement := childByPredicate(blockBody, isConstraintsStatement)
	if constraintsStatement != nil {
		constraintsLiteral := childByTypeRec(constraintsStatement.NamedChild(1), "template_literal")
		content := constraintsLiteral.Content(sourceCode)
		return content, nil
	}
	return "", nil
}

func providerHashes(providerBlock *sitter.Node, sourceCode []byte) ([]string, error) {
	isHashesStatement := func (block *sitter.Node) bool {
		return block.NamedChildCount() == 2 && block.NamedChild(0).Content(sourceCode) == "hashes"
	}

	blockBody := childByType(providerBlock, "body")
	if blockBody == nil {
		return []string{}, errors.New("failed to find block body in provider block")
	}
	hashesStatement := childByPredicate(blockBody, isHashesStatement)
	hashesTuple := hashesStatement.NamedChild(1).NamedChild(0).NamedChild(0)
	hashExpressions := childrenByType(hashesTuple, "expression")
	hashLiterals := []string{}
	for i := 0; i < len(hashExpressions); i++ {
		literal := childByTypeRec(hashExpressions[i], "template_literal").Content(sourceCode)
		hashLiterals = append(hashLiterals, literal)
	}
	return hashLiterals, nil
}