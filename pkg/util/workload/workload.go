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
	builder.WriteString(str.Repeat("-", 78))
	output.Println(builder.String())
}

func (w *Workload) PrintFooter(result string, url string, queued string, finished string) {
	builder := &strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s\n", str.Repeat("-", 78)))
	builder.WriteString(fmt.Sprintf("Result    | %s\n", result))
	builder.WriteString(fmt.Sprintf("Build     | %s\n", url))
	builder.WriteString(fmt.Sprintf("Queued    | %s\n", queued))
	builder.WriteString(fmt.Sprintf("Finished  | %s\n", finished))
	builder.WriteString(fmt.Sprintf("Elapsed   | %s\n", ""))
	builder.WriteString(str.Repeat("=", 78))
	output.Println(builder.String())
}
