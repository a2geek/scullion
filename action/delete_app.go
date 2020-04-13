package action

import (
	"fmt"
	"net/url"
	"scullion/task"
)

func DeleteApp(taskItem task.Item) {
	taskItem.Metadata.Logger.Infof("Deleting app '%s' in space '%s' of org '%s'",
		taskItem.Variables.App.Name,
		taskItem.Variables.Space.Name,
		taskItem.Variables.Org.Name)

	q := url.Values{}
	q.Add("q", fmt.Sprintf("app_guid:%s", taskItem.Variables.App.Guid))
	taskItem.Metadata.Logger.Debugf("service bindings query: %s", q)
	bindings, err := taskItem.Metadata.Client.ListServiceBindingsByQuery(q)
	taskItem.Metadata.Logger.Debugf("found %d service bindings", len(bindings))
	if err != nil {
		taskItem.Metadata.Logger.Errorf("Unable to identify bindings for application '%s' (%s): %v",
			taskItem.Variables.App.Name,
			taskItem.Variables.App.Guid)
	}

	for _, binding := range bindings {
		taskItem.Metadata.Logger.Infof("Removing service binding '%s' (%s) for app '%s' (%s)",
			binding.Name, binding.Guid,
			taskItem.Variables.App.Name, taskItem.Variables.App.Guid)
		err = taskItem.Metadata.Client.DeleteServiceBinding(binding.Guid)
		if err != nil {
			taskItem.Metadata.Logger.Errorf("Unable to delete service binding '%s' (%s) for app '%s' (%s)",
				binding.Name, binding.Guid,
				taskItem.Variables.App.Name, taskItem.Variables.App.Guid)
		}
	}

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
