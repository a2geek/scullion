package action

import (
	"fmt"
	"scullion/task"
)

func StopApp(taskItem task.Item) {
	fmt.Printf("[%s] Stopping app '%s' in space '%s' of org '%s'\n",
		taskItem.Metadata.Name,
		taskItem.Variables.App.Name,
		taskItem.Variables.Space.Name,
		taskItem.Variables.Org.Name)

	if err := taskItem.Metadata.Client.StopApp(taskItem.Variables.App.Guid); err != nil {
		fmt.Printf("Unable to stop application '%s' in space '%s' of org '%s': %s\n",
			taskItem.Variables.App.Name,
			taskItem.Variables.Space.Name,
			taskItem.Variables.Org.Name,
			err)
	}
}

func DryRunStopApp(taskItem task.Item) {
	fmt.Printf("[%s] (DRY-RUN) Stopping app '%s' in space '%s' of org '%s'\n",
		taskItem.Metadata.Name,
		taskItem.Variables.App.Name,
		taskItem.Variables.Space.Name,
		taskItem.Variables.Org.Name)
}
