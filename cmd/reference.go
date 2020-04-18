package cmd

import (
	"fmt"
	"reflect"
	"scullion/task"
	"strings"
)

type Reference struct {
}

func (cmd *Reference) Execute(args []string) error {
	varType := reflect.TypeOf(task.RunEnv{})

	fmt.Println("Note that methods prefixed with an operation usually have that operator overloaded.")
	fmt.Println("Thereform, 'Add' for a Time and Duration can be expressed 'time + duration'.")
	fmt.Println("This list is dynamically generated at run time, so should be accurate for your version.")
	fmt.Println()

	fmt.Println("Operations:")
	for i := 0; i < varType.NumMethod(); i++ {
		method := varType.Method(i)
		fmt.Printf("  %s(%s) %s\n", method.Name, strings.Join(inFrom(method), ", "), outFrom(method))
	}

	return nil
}

func inFrom(method reflect.Method) []string {
	types := make([]string, 0)
	// Skip 0 as that is the type...
	for j := 1; j <= method.Type.NumIn()-1; j++ {
		field := method.Type.In(j)
		types = append(types, field.String())
	}
	return types
}

func outFrom(method reflect.Method) string {
	return method.Type.Out(0).String()
}
