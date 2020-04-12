package worker

import (
	"fmt"
	"scullion/log"
	"scullion/task"
	"sync"
)

func Action(num int, actionChan <-chan task.Item, wg *sync.WaitGroup, logLevel string) {
	logger, err := log.NewLogger(fmt.Sprintf("action worker %d", num), logLevel)
	if err != nil {
		panic(err)
	}
	logger.Info("Launched")

	defer wg.Done()
	for taskItem := range actionChan {
		taskItem.Metadata.Logger.Infof("Performing action on app '%s' in space '%s' of org '%s'\n",
			taskItem.Variables.App.Name, taskItem.Variables.Space.Name, taskItem.Variables.Org.Name)
		taskItem.Metadata.Action(taskItem)
	}

	logger.Error("exited")
}
