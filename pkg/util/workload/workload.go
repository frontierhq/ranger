package workload

import (
	"fmt"
	"strings"

	"github.com/frontierdigital/utils/output"
	"github.com/frontierdigital/utils/str"
)

type Workload struct {
	Name    string `yaml:"name"`
	Source  string `yaml:"source"`
	Version string `yaml:"version"`
}

func (w *Workload) GetSourceProjectAndRepositoryNames() (string, string) {
	sourceParts := strings.Split(w.Source, "/")
	return sourceParts[0], sourceParts[1]
}

func (w *Workload) PrintHeader() {
	builder := &strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s\n", str.Repeat("=", 78)))
	builder.WriteString(fmt.Sprintf("Name     | %s\n", w.Name))
	builder.WriteString(fmt.Sprintf("Source   | %s\n", w.Source))
	builder.WriteString(fmt.Sprintf("Version  | %s\n", w.Version))
	builder.WriteString(strings.Repeat("-", 78))
	output.Println(builder.String())
}
