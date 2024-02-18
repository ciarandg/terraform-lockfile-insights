package insights

import (
	"encoding/json"
	"errors"

	"github.com/ciarandg/terraform-lockfile-insights/lockfile"
)

type Insights struct {
	Versions map[string][]string `json:"versions"`
}

func GetInsights(lockfiles map[string]lockfile.Lockfile) map[string]Insights {
	out := map[string]Insights{}

	for filePath, lockfile := range lockfiles {
		for name, provider := range lockfile.ProviderBlocks {
			val, ok := out[name]
			if !ok {
				val = Insights{map[string][]string{}}
			}

			val2, ok2 := val.Versions[provider.Version]
			if !ok2 {
				val2 = []string{}
			}
			val2 = append(val2, filePath)
			val.Versions[provider.Version] = val2
			
			out[name] = val
		}
	}

	return out
}

func GetInsightsJson(lockfiles map[string]lockfile.Lockfile) (string, error) {
	insights := GetInsights(lockfiles)
	jsonData, err := json.Marshal(insights)
	if err != nil {
		return "", errors.New("could not marshal JSON")
	}
	return string(jsonData), nil
}

func GetInsightsJsonPretty(lockfiles map[string]lockfile.Lockfile) (string, error) {
	insights := GetInsights(lockfiles)
	jsonData, err := json.MarshalIndent(insights, "", "  ")
	if err != nil {
		return "", errors.New("could not marshal JSON")
	}
	return string(jsonData), nil
}