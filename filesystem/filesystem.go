package filesystem

import (
	"os"
	"path/filepath"
)

func GetLockfilePaths(rootDirPath string) ([]string, error) {
	out := []string{}

	err := filepath.Walk(rootDirPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
		if filepath.Base(path) == ".terraform.lock.hcl" {
			out = append(out, path)
		}
        return nil
    })

	if err != nil {
		return []string{}, err
	}
	return out, nil
}