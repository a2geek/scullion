package worker

import (
	"fmt"
	"scullion/action"
	"scullion/config"
	"scullion/task"
	"time"

	"github.com/cloudfoundry-community/go-cfclient"
)

func Task(num int, taskDef config.TaskDef, client *cfclient.Client, orgChan chan<- task.Item) {
	fmt.Printf("Launched task worker %d\n", num)
	metadata, err := task.NewMetadata(taskDef, client, action.Log)
	if err != nil {
		panic(err)
	}
	dur, err := time.ParseDuration(taskDef.Schedule.Frequency)
	if err != nil {
		panic(err)
	}
	for t := range time.Tick(dur) {
		fmt.Printf("[%s] Tick at %s\n", taskDef.Name, t)
		taskItem := task.Item{
			Metadata: metadata,
		}
		orgChan <- taskItem
	}
}
