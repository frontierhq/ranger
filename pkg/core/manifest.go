package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/frontierdigital/utils/output"
	"gopkg.in/yaml.v2"
)

func (m *Manifest) PrintHeader() {
	builder := &strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s\n", strings.Repeat("~", 78)))
	builder.WriteString(fmt.Sprintf("Environment  | %s\n", m.Environment))
	builder.WriteString(fmt.Sprintf("Set          | %s\n", m.Set))
	builder.WriteString(fmt.Sprintf("Version      | %d\n", m.Version))
	builder.WriteString(fmt.Sprintf("%s\n", strings.Repeat("~", 78)))
	output.Println(builder.String())
}

func (m *Manifest) PrintWorkloadsSummary() {
	builder := &strings.Builder{}
	builder.WriteString("Workloads:\n")
	for _, w := range m.Workloads {
		builder.WriteString(fmt.Sprintf(" * %s (type: %s, version: %s)\n", w.Name, w.Type, w.Version))
	}
	output.Println(builder.String())
}

func (m *Manifest) Save() error {
	data, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	builder := &strings.Builder{}
	builder.WriteString("---\n")
	builder.WriteString(string(data[:]))

	file, err := os.Create(m.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprint(file, builder.String())
	if err != nil {
		return err
	}

	return nil
}

func LoadManifest(filePath string) (Manifest, error) {
	manifest := Manifest{}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return manifest, err
	}

	err = yaml.Unmarshal(data, &manifest)
	if err != nil {
		return manifest, err
	}

	manifest.FilePath = filePath

	return manifest, nil
}

func LoadManifestFromString(manifestContent string) (Manifest, error) {
	manifest := Manifest{}

	err := yaml.Unmarshal([]byte(manifestContent), &manifest)
	if err != nil {
		return manifest, err
	}

	return manifest, nil
}
