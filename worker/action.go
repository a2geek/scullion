package worker

import (
	"fmt"
	"scullion/task"
)

func Action(num int, actionChan <-chan task.Item) {
	fmt.Printf("Launched action worker %d\n", num)
	for {
		taskItem := <-actionChan
		fmt.Printf("[%s] Performing action on app '%s' in space '%s' of org '%s'\n", taskItem.Metadata.Name,
			taskItem.Variables.App.Name, taskItem.Variables.Space.Name, taskItem.Variables.Org.Name)
		taskItem.Metadata.Action(taskItem)
	}
}
