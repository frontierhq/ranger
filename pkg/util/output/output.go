package output

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	// Println prints in the default colour.
	Println = printlnInColor(nil)
	// PrintlnLog is used to print log messages for debugging. It prints in a faded colour.
	PrintlnLog = printlnInColor(color.New(color.FgBlack, color.Faint).PrintfFunc())
	// PrintlnInfo is used to print standard output always visible to the user. It prints
	// in the default console colour.
	PrintlnInfo = printlnInColor(color.White)
	// PrintlnWarn is used to print warning messages and are highlighted in yellow.
	PrintlnWarn = printlnInColor(color.Yellow)
	// PrintlnError is used to print errors in red and indicate a fatal problem.
	PrintlnError = printlnInColor(color.Red)

	// Printlnf prints a formatted string in the default colour.
	Printlnf = printlnfInColor(nil)
	// PrintlnfLog is used to print a formatted string as a log message for debugging.
	// It prints in a faded colour.
	PrintlnfLog = printlnfInColor(color.New(color.FgBlack, color.Faint).PrintfFunc())
	// PrintlnfInfo is used to print a formatted string as standard output always visible
	// to the user. It prints in the default console colour.
	PrintlnfInfo = printlnfInColor(color.White)
	// PrintlnfWarn is used to print a formatted string as a warning messages and are
	// highlighted in yellow.
	PrintlnfWarn = printlnfInColor(color.Yellow)
	// PrintlnfError is used to print error messages as formatted strings and indicate
	// a fatal problem.
	PrintlnfError = printlnfInColor(color.Red)

	// Println prints a formatted string in the default colour.
	Printf = printfInColor(nil)
	// PrintlnfLog is used to print a formatted string as a log message for debugging.
	// It prints in a faded colour.
	PrintfLog = printfInColor(color.New(color.FgBlack, color.Faint).PrintfFunc())
	// PrintlnfInfo is used to print a formatted string as standard output always visible
	// to the user. It prints in the default console colour.
	PrintfInfo = printfInColor(color.White)
	// PrintlnfWarn is used to print a formatted string as a warning messages and are
	// highlighted in yellow.
	PrintfWarn = printfInColor(color.Yellow)
	// PrintlnfError is used to print error messages as formatted strings and indicate
	// a fatal problem.
	PrintfError = printfInColor(color.Red)
)

func printlnInColor(color func(string, ...interface{})) func(...interface{}) {
	printer := func(args ...interface{}) {
		if color != nil {
			color(fmt.Sprintln(args...))
		} else {
			fmt.Println(args...)
		}
	}
	return printer
}

func printlnfInColor(color func(string, ...interface{})) func(string, ...interface{}) {
	printer := func(format string, args ...interface{}) {
		if color != nil {
			color(fmt.Sprintln(fmt.Sprintf(format, args...)))
		} else {
			fmt.Println(fmt.Sprintf(format, args...))
		}
	}
	return printer
}

func printfInColor(color func(string, ...interface{})) func(string, ...interface{}) {
	printer := func(format string, args ...interface{}) {
		if color != nil {
			color(fmt.Sprintf(format, args...))
		} else {
			fmt.Printf(format, args...)
		}
	}
	return printer
}
