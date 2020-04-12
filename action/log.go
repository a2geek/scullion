package action

import (
	"scullion/task"
)

func Log(taskItem task.Item) {
	taskItem.Metadata.Logger.Infof("Action would be taken on app '%s' in space '%s' of org '%s'",
		taskItem.Variables.App.Name,
		taskItem.Variables.Space.Name,
		taskItem.Variables.Org.Name)
}
