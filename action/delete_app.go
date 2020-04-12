package action

import (
	"scullion/task"
)

func DeleteApp(taskItem task.Item) {
	taskItem.Metadata.Logger.Infof("Deleting app '%s' in space '%s' of org '%s'",
		taskItem.Variables.App.Name,
		taskItem.Variables.Space.Name,
		taskItem.Variables.Org.Name)

	if err := taskItem.Metadata.Client.DeleteApp(taskItem.Variables.App.Guid); err != nil {
		taskItem.Metadata.Logger.Errorf("Unable to delete application '%s' in space '%s' of org '%s': %s",
			taskItem.Variables.App.Name,
			taskItem.Variables.Space.Name,
			taskItem.Variables.Org.Name,
			err)
	}
}

func DryRunDeleteApp(taskItem task.Item) {
	taskItem.Metadata.Logger.Infof("(DRY-RUN) Deleting app '%s' in space '%s' of org '%s'",
		taskItem.Variables.App.Name,
		taskItem.Variables.Space.Name,
		taskItem.Variables.Org.Name)
}
