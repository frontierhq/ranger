package core

import (
	"fmt"
	"strings"

	"github.com/gofrontier-com/go-utils/output"
)

func (w *WorkloadInstance) GetTypeProjectAndRepositoryNames() (string, string) {
	typeParts := strings.Split(w.Type, "/")
	return typeParts[0], typeParts[1]
}

func (w *WorkloadInstance) PrintHeader() {
	builder := &strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s\n", strings.Repeat("=", 78)))
	builder.WriteString(fmt.Sprintf("Name     | %s\n", w.Name))
	builder.WriteString(fmt.Sprintf("Type     | %s\n", w.Type))
	builder.WriteString(fmt.Sprintf("Version  | %s\n", w.Version))
	builder.WriteString(strings.Repeat("-", 78))
	output.Println(builder.String())
}
