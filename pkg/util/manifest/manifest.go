package manifest

import (
	"os"

	"github.com/frontierdigital/ranger/pkg/cmd/app"

	"gopkg.in/yaml.v2"
)

func LoadManifest(path string) (app.Manifest, error) {
	manifest := app.Manifest{}

	data, err := os.ReadFile(path)
	if err != nil {
		return manifest, err
	}

	err = yaml.Unmarshal(data, &manifest)
	if err != nil {
		return manifest, err
	}

	return manifest, nil
}
