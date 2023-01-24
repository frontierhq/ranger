package util

import (
	"os"

	"github.com/frontierdigital/ranger/core"
	"gopkg.in/yaml.v2"
)

func LoadManifest(path string) (core.Manifest, error) {
	manifest := core.Manifest{}

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
