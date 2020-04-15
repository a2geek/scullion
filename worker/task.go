package worker

import (
	"fmt"
	"scullion/action"
	"scullion/config"
	"scullion/log"
	"scullion/option"
	"scullion/task"
	"time"

	"github.com/cloudfoundry-community/go-cfclient"
)

func Task(num int, taskDef config.TaskDef, client *cfclient.Client, orgChan chan<- task.Item, runOptions option.RunOptions) {
	logger, err := log.NewLogger(fmt.Sprintf("task worker %d", num), runOptions.Level, runOptions.NoDate)
	if err != nil {
		panic(err)
	}
	logger.Info("Launched")

	// should something unrecoverable occur, close the channel which cascades ... to make it obvious!?
	defer close(orgChan)

	actionFunc, err := action.NewActionFunc(taskDef.Filters.Action, runOptions.DryRun)
	if err != nil {
		logger.Errorf("halting: %v", err)
		return
	}

	metadata, err := task.NewMetadata(taskDef, client, actionFunc, runOptions)
	if err != nil {
		logger.Errorf("halting: %v", err)
		return
	}

	dur, err := time.ParseDuration(taskDef.Schedule.Frequency)
	if err != nil {
		logger.Errorf("halting: %v", err)
		return
	}

	for t := range time.Tick(dur) {
		metadata.Logger.Infof("Tick at %s", t)
		taskItem := task.Item{
			Metadata: metadata,
		}
		orgChan <- taskItem
	}

	logger.Error("exiting")
}
