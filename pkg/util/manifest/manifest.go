package manifest

import (
	"os"

	"github.com/frontierdigital/ranger/pkg/cmd/app/type/manifest"

	"gopkg.in/yaml.v2"
)

func LoadManifest(path string) (manifest.Manifest, error) {
	manifest := manifest.Manifest{}

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
