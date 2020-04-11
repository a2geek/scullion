package action

import (
	"fmt"
	"scullion/task"
)

const logAction = "log"
const deleteAppAction = "delete-app"
const stopAppAction = "stop-app"

var actualActionMap = map[string]func(task.Item){
	logAction:       Log,
	deleteAppAction: DeleteApp,
	stopAppAction:   StopApp,
}
var dryRunActionMap = map[string]func(task.Item){
	logAction:       Log,
	deleteAppAction: DryRunDeleteApp,
	stopAppAction:   DryRunStopApp,
}

func NewActionFunc(name string, dryRun bool) (func(task.Item), error) {
	actionMap := actualActionMap
	if dryRun {
		actionMap = dryRunActionMap
	}
	if fn, ok := actionMap[name]; ok {
		return fn, nil
	}
	return nil, fmt.Errorf("action '%s' not found", name)
}
