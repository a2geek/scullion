package action

import (
	"scullion/task"
)

func StopApp(taskItem task.Item) {
	taskItem.Metadata.Logger.Infof("Stopping app '%s' in space '%s' of org '%s'",
		taskItem.Variables.App.Name,
		taskItem.Variables.Space.Name,
		taskItem.Variables.Org.Name)

	if err := taskItem.Metadata.Client.StopApp(taskItem.Variables.App.Guid); err != nil {
		taskItem.Metadata.Logger.Errorf("Unable to stop application '%s' in space '%s' of org '%s': %s",
			taskItem.Variables.App.Name,
			taskItem.Variables.Space.Name,
			taskItem.Variables.Org.Name,
			err)
	}
}

func DryRunStopApp(taskItem task.Item) {
	taskItem.Metadata.Logger.Infof("(DRY-RUN) Stopping app '%s' in space '%s' of org '%s'",
		taskItem.Variables.App.Name,
		taskItem.Variables.Space.Name,
		taskItem.Variables.Org.Name)
}
