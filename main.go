package main

import (
	"fmt"
	"os"

	"github.com/ciarandg/provider-finder/lockfile"
)

func main() {
	filePath := "example.lock.hcl"
	summary, err := lockfile.GenerateLockfileSummary(filePath)
	if err != nil {
		fmt.Printf("Encountered an error while processing %s: %s\n", filePath, err)
		os.Exit(1)
	}
	fmt.Println(summary)
}
