package main

import (
	"fmt"
	"os"

	"github.com/ciarandg/provider-finder/filesystem"
	"github.com/ciarandg/provider-finder/lockfile"
)

func main() {
	var dirPath string
	if len(os.Args) > 1 {
		dirPath = os.Args[1]
	} else {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error: could not determine current working directory")
		}
		dirPath = wd
	}

	lockfilePaths, err := filesystem.GetLockfilePaths(dirPath)
	if err != nil {
		fmt.Printf("Encountered an error while looking for lockfiles: %s\n", err)
		os.Exit(1)
	}

	for i := 0; i < len(lockfilePaths); i++ {
		filePath := lockfilePaths[i]
		lockfile, err := lockfile.NewLockfileFromPath(filePath)
		if err != nil {
			fmt.Printf("Encountered an error while initializing lockfile %s: %s\n", filePath, err)
			os.Exit(1)
		}
		fmt.Println(filePath)
		fmt.Println(lockfile.ProviderBlocks)
	}
}
