package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrontier-com/go-utils/output"
)

func (d *WorkloadResult) PrintResult() {
	var elasped time.Duration
	if d.QueueTime != nil && d.FinishTime != nil {
		elasped = d.FinishTime.Sub(*d.QueueTime)
	}

	var queueTime string
	if d.QueueTime != nil {
		queueTime = d.QueueTime.Format(time.RFC1123)
	}

	var finishTime string
	if d.FinishTime != nil {
		finishTime = d.FinishTime.Format(time.RFC1123)
	}

	builder := &strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s\n", strings.Repeat("-", 78)))
	builder.WriteString(fmt.Sprintf("Result    | %s\n", d.Status))
	builder.WriteString(fmt.Sprintf("Link      | %s\n", d.Link))
	builder.WriteString(fmt.Sprintf("Queued    | %s\n", queueTime))
	builder.WriteString(fmt.Sprintf("Finished  | %s\n", finishTime))
	builder.WriteString(fmt.Sprintf("Elapsed   | %s\n", elasped))
	builder.WriteString(fmt.Sprintf("%s\n", strings.Repeat("=", 78)))
	output.Println(builder.String())
}
