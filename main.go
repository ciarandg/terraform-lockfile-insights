package main

import (
	"fmt"
	"os"

	"github.com/ciarandg/provider-finder/lockfile"
)

func main() {
	filePath := "example.lock.hcl"
	lockfile, err := lockfile.NewLockfile(filePath)
	if err != nil {
		fmt.Printf("Encountered an error while initializing lockfile %s: %s\n", filePath, err)
		os.Exit(1)
	}
	fmt.Println(lockfile.ProviderBlocks)
}
