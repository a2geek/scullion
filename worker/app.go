package worker

import (
	"fmt"
	"net/url"
	"scullion/log"
	"scullion/task"
	"sync"
)

func App(num int, appChan <-chan task.Item, actionChan chan<- task.Item, wg *sync.WaitGroup, logLevel string) {
	logger, err := log.NewLogger(fmt.Sprintf("app worker %d", num), logLevel)
	if err != nil {
		panic(err)
	}
	logger.Info("Launched")

	// if app channel closes, let's close the action channel as well!
	defer close(actionChan)
	defer wg.Done()

	for taskItem := range appChan {
		q := url.Values{}
		q.Add("q", fmt.Sprintf("space_guid:%s", taskItem.Variables.Space.Guid))
		apps, err := taskItem.Metadata.Client.ListAppsByQuery(q)
		if err != nil {
			taskItem.Metadata.Logger.Errorf("error querying for space applications '%s' (%s): %v",
				taskItem.Variables.Space.Name, taskItem.Variables.Space.Guid, err)
			continue
		}
		for _, app := range apps {
			variables := task.Variables{
				Org:   taskItem.Variables.Org,
				Space: taskItem.Variables.Space,
				App:   app,
			}
			isTrue, err := taskItem.Metadata.IsAppMatch(variables)
			if err != nil {
				taskItem.Metadata.Logger.Errorf("unable to determine true/false: %v", err)
				continue
			}
			if isTrue {
				taskItem.Metadata.Logger.Infof("Matched app '%s' in space '%s' of org '%s'",
					app.Name, variables.Space.Name, variables.Org.Name)
				newTask := task.Item{
					Variables: variables,
					Metadata:  taskItem.Metadata,
				}
				actionChan <- newTask
			} else {
				taskItem.Metadata.Logger.Infof("Skipped app '%s' in space '%s' of org '%s'",
					app.Name, variables.Space.Name, variables.Org.Name)
			}
		}
	}

	logger.Error("exiting")
}
