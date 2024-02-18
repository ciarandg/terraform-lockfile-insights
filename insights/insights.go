package insights

import (
	"github.com/ciarandg/provider-finder/lockfile"
)

type Insights struct {
	versions map[string][]string
}

func GetInsights(lockfiles map[string]lockfile.Lockfile) map[string]Insights {
	out := map[string]Insights{}

	for filePath, lockfile := range lockfiles {
		for name, provider := range lockfile.ProviderBlocks {
			val, ok := out[name]
			if !ok {
				val = Insights{map[string][]string{}}
			}

			val2, ok2 := val.versions[provider.Version]
			if !ok2 {
				val2 = []string{}
			}
			val2 = append(val2, filePath)
			val.versions[provider.Version] = val2
			
			out[name] = val
		}
	}

	return out
}