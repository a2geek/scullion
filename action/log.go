package action

import (
	"fmt"
	"scullion/task"
)

func Log(taskItem task.Item) {
	fmt.Printf("[%s] Action would be taken on app '%s' in space '%s' of org '%s'\n",
		taskItem.Metadata.Name,
		taskItem.Variables.Org.Name,
		taskItem.Variables.Space.Name,
		taskItem.Variables.App.Name)
}
