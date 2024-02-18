package lockfile

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var contentsSingleProvider = []byte(`
provider "example.com/provider" {
  version     = "foo"
  constraints = "bar"
  hashes = [
    "cat",
    "dog",
    "frog",
  ]
}
`)

var contentsMultipleProviders = []byte(`
provider "google.com/apple/banana" {
  version = "pear"
  hashes = [
    "mercury",
    "earth",
    "venus",
    "neptune",
  ]
}

provider "example.com/provider" {
  version     = "foo"
  constraints = "bar"
  hashes = [
    "cat",
    "dog",
    "frog",
  ]
}

provider "example.com/other-provider" {
  version = "tundra"
  hashes = [
    "plains",
    "wetland",
    "prairies",
  ]
}
`)

var contentsInconsistentWhitespace = []byte(`
provider "google.com/apple/banana" {
  version = "pear"
  hashes = [ "mercury", "earth", "venus", "neptune", ]
}
  provider "example.com/provider" {
  
  
    	version     = "foo"
  
    constraints = "bar"
  
    hashes = [ "cat",
      "dog",
  
      "frog",
    ]
  
  }

provider "example.com/other-provider" {
version = "tundra"
hashes = [
"plains",
"wetland",
"prairies",
]
}
`)

var contentsNoConstraints = []byte(`
provider "example.com/provider" {
  version     = "foo"
  hashes = [
    "cat",
    "dog",
    "frog",
  ]
}
`)

var contentsInvalidEmpty = []byte("")

var contentsInvalidNotHcl = []byte(`
{
  "name": "John Doe",
  "age": 30,
  "city": "New York",
  "isStudent": false,
  "hobbies": ["reading", "hiking", "cooking"]
}
`)

var contentsInvalidDuplicateProvider = []byte(`
provider "google.com/apple/banana" {
  version = "pear"
  hashes = [
    "mercury",
    "earth",
    "venus",
    "neptune",
  ]
}

provider "example.com/provider" {
  version     = "foo"
  constraints = "bar"
  hashes = [
    "cat",
    "dog",
    "frog",
  ]
}

provider "google.com/apple/banana" {
  version = "pear"
  hashes = [
    "mercury",
    "earth",
    "venus",
    "neptune",
  ]
}
`)

var contentsInvalidNoVersion = []byte(`
provider "example.com/provider" {
  constraints = "bar"
  hashes = [
    "cat",
    "dog",
    "frog",
  ]
}
`)

var contentsInvalidNoHashes = []byte(`
provider "example.com/provider" {
  version     = "foo"
  constraints = "bar"
}
`)

var contentsInvalidSingleQuoteStrings = []byte(`
provider 'example.com/provider' {
  version     = "foo"
  constraints = 'bar'
  hashes = [
    'cat',
    "dog",
    'frog',
  ]
}
`)

func TestNewLockfileSingleProvider(t *testing.T) {
	l, err := NewLockfile(contentsSingleProvider)
	assert.Nil(t, err)
	assert.Equal(t, len(l.ProviderBlocks), 1)
	block, ok := l.ProviderBlocks["example.com/provider"]
	assert.True(t, ok)
	assert.Equal(t, block.Version, "foo")
	assert.Equal(t, block.Constraints, "bar")
	assert.Equal(t, block.Hashes, []string{"cat", "dog", "frog"})
}

func TestNewLockfileMultipleProviders(t *testing.T) {
	providerNames := []string{"google.com/apple/banana", "example.com/provider", "example.com/other-provider"}
	providerVersions := []string{"pear", "foo", "tundra"}
	providerConstraints := []string{"", "bar", ""}
	providerHashes := [][]string{
		{"mercury", "earth", "venus", "neptune"},
		{"cat", "dog", "frog"},
		{"plains", "wetland", "prairies"},
	}

	l, err := NewLockfile(contentsMultipleProviders)
	assert.Nil(t, err)
	assert.Equal(t, len(l.ProviderBlocks), len(providerNames))
	for i := 0; i < len(providerNames); i++ {
		name := providerNames[i]
		block := l.ProviderBlocks[name]
		assert.Equal(t, block.Version, providerVersions[i])
		assert.Equal(t, block.Constraints, providerConstraints[i])
		assert.Equal(t, block.Hashes, providerHashes[i])
	}
}

func TestNewLockfileInconsistentWhitespace(t *testing.T) {
	providerNames := []string{"google.com/apple/banana", "example.com/provider", "example.com/other-provider"}
	providerVersions := []string{"pear", "foo", "tundra"}
	providerConstraints := []string{"", "bar", ""}
	providerHashes := [][]string{
		{"mercury", "earth", "venus", "neptune"},
		{"cat", "dog", "frog"},
		{"plains", "wetland", "prairies"},
	}

	l, err := NewLockfile(contentsInconsistentWhitespace)
	assert.Nil(t, err)
	assert.Equal(t, len(l.ProviderBlocks), len(providerNames))
	for i := 0; i < len(providerNames); i++ {
		name := providerNames[i]
		block := l.ProviderBlocks[name]
		assert.Equal(t, block.Version, providerVersions[i])
		assert.Equal(t, block.Constraints, providerConstraints[i])
		assert.Equal(t, block.Hashes, providerHashes[i])
	}
}

func TestNewLockfileNoConstraints(t *testing.T) {
	l, err := NewLockfile(contentsNoConstraints)
	assert.Nil(t, err)
	assert.Equal(t, len(l.ProviderBlocks), 1)
	block, ok := l.ProviderBlocks["example.com/provider"]
	assert.True(t, ok)
	assert.Equal(t, block.Version, "foo")
	assert.Equal(t, block.Constraints, "")
	assert.Equal(t, block.Hashes, []string{"cat", "dog", "frog"})
}

func TestNewLockfileInvalidEmpty(t *testing.T) {
	l, err := NewLockfile(contentsInvalidEmpty)
	assert.NotNil(t, err)
	assert.Equal(t, len(l.ProviderBlocks), 0)
}

func TestNewLockfileInvalidNotHcl(t *testing.T) {
	l, err := NewLockfile(contentsInvalidNotHcl)
	assert.NotNil(t, err)
	assert.Equal(t, len(l.ProviderBlocks), 0)
}

func TestNewLockfileInvalidDuplicateLockfile(t *testing.T) {
	l, err := NewLockfile(contentsInvalidDuplicateProvider)
	assert.NotNil(t, err)
	assert.Equal(t, len(l.ProviderBlocks), 0)
}

func TestNewLockfileInvalidNoVersion(t *testing.T) {
	l, err := NewLockfile(contentsInvalidNoVersion)
	assert.NotNil(t, err)
	assert.Equal(t, len(l.ProviderBlocks), 0)
}

func TestNewLockfileInvalidNoHashes(t *testing.T) {
	l, err := NewLockfile(contentsInvalidNoHashes)
	assert.NotNil(t, err)
	assert.Equal(t, len(l.ProviderBlocks), 0)
}

func TestNewLockfileInvalidSingleQuoteStrings(t *testing.T) {
	l, err := NewLockfile(contentsInvalidSingleQuoteStrings)
	assert.NotNil(t, err)
	assert.Equal(t, len(l.ProviderBlocks), 0)
}

func TestNewLockfileFromPath(t *testing.T) {
  tmpFile, err := os.CreateTemp("/tmp", "TestNewLockfileFromPath")
  assert.Nil(t, err)
  tmpFile.Write(contentsSingleProvider)
  tmpFile.Close()

  l, err := NewLockfileFromPath(tmpFile.Name())
	assert.Nil(t, err)
	assert.Equal(t, len(l.ProviderBlocks), 1)
	block, ok := l.ProviderBlocks["example.com/provider"]
	assert.True(t, ok)
	assert.Equal(t, block.Version, "foo")
	assert.Equal(t, block.Constraints, "bar")
	assert.Equal(t, block.Hashes, []string{"cat", "dog", "frog"})
}

func TestNewLockfileFromPathInvalidNotRealFile(t *testing.T) {
  filePath := "/tmp/TestNewLockfileFromPathInvalidNotRealFile"
  _, err := os.Stat(filePath)
  assert.Contains(t, err.Error(), "no such file or directory")

	l, err := NewLockfileFromPath(filePath)
	assert.NotNil(t, err)
	assert.Equal(t, len(l.ProviderBlocks), 0)
}