package core

import "fmt"

func (e *WorkloadDestroyPreventedError) Error() string {
	return fmt.Sprintf("destroy prevented for workload with name '%s'", e.Workload.Name)
}
