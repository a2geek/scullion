package action

import (
	"fmt"
	"scullion/task"
)

func DeleteApp(taskItem task.Item) {
	fmt.Printf("[%s] Deleting app '%s' in space '%s' of org '%s'\n",
		taskItem.Metadata.Name,
		taskItem.Variables.App.Name,
		taskItem.Variables.Space.Name,
		taskItem.Variables.Org.Name)

	if err := taskItem.Metadata.Client.DeleteApp(taskItem.Variables.App.Guid); err != nil {
		fmt.Printf("Unable to delete application: %s\n", err)
	}
}

func DryRunDeleteApp(taskItem task.Item) {
	fmt.Printf("[%s] (DRY-RUN) Deleting app '%s' in space '%s' of org '%s'\n",
		taskItem.Metadata.Name,
		taskItem.Variables.App.Name,
		taskItem.Variables.Space.Name,
		taskItem.Variables.Org.Name)
}
