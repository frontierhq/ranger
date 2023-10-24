package workload

import (
	"fmt"
	"strings"

	"github.com/frontierdigital/utils/output"
	"github.com/frontierdigital/utils/str"
)

type Workload struct {
	ExtraParameters []ExtraParameter `yaml:"extraParameters"`
	Name            string           `yaml:"name"`
	Type            string           `yaml:"type"`
	Version         string           `yaml:"version"`
}

type ExtraParameter struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func (w *Workload) GetTypeProjectAndRepositoryNames() (string, string) {
	typeParts := strings.Split(w.Type, "/")
	return typeParts[0], typeParts[1]
}

func (w *Workload) PrintHeader() {
	builder := &strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s\n", str.Repeat("=", 78)))
	builder.WriteString(fmt.Sprintf("Name     | %s\n", w.Name))
	builder.WriteString(fmt.Sprintf("Type     | %s\n", w.Type))
	builder.WriteString(fmt.Sprintf("Version  | %s\n", w.Version))
	builder.WriteString(strings.Repeat("-", 78))
	output.Println(builder.String())
}
