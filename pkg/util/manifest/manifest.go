package manifest

import (
	"fmt"
	"os"
	"strings"

	"github.com/frontierdigital/ranger/pkg/util/workload"
	"github.com/frontierdigital/utils/output"

	"gopkg.in/yaml.v2"
)

type Manifest struct {
	Environment string               `yaml:"environment"`
	Layer       string               `yaml:"layer"`
	Version     int64                `yaml:"version"`
	Workloads   []*workload.Workload `yaml:"workloads"`
}

func (m *Manifest) PrintHeader() {
	builder := &strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s\n", strings.Repeat("~", 78)))
	builder.WriteString(fmt.Sprintf("Manifest     | %s-%s\n", m.Environment, m.Layer))
	builder.WriteString(fmt.Sprintf("Layer        | %s\n", m.Layer))
	builder.WriteString(fmt.Sprintf("Environment  | %s\n", m.Environment))
	builder.WriteString(fmt.Sprintf("Version      | %d\n", m.Version))
	builder.WriteString(fmt.Sprintf("%s\n", strings.Repeat("~", 78)))
	output.Println(builder.String())
}

func (m *Manifest) PrintWorkloadsSummary() {
	builder := &strings.Builder{}
	builder.WriteString("Workloads:\n")
	for _, w := range m.Workloads {
		builder.WriteString(fmt.Sprintf(" * %s (%s @ %s)\n", w.Name, w.Source, w.Version))
	}
	output.Println(builder.String())
}

func LoadManifest(path string) (Manifest, error) {
	manifest := Manifest{}

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
