package worker

import (
	"fmt"
	"net/url"
	"scullion/log"
	"scullion/task"
	"sync"
)

func Space(num int, spaceChan <-chan task.Item, appChan chan<- task.Item, wg *sync.WaitGroup, logLevel string) {
	logger, err := log.NewLogger(fmt.Sprintf("space worker %d", num), logLevel)
	if err != nil {
		panic(err)
	}
	logger.Info("Launched")

	// if space channel closes, let's close the app channel as well!
	defer close(appChan)
	defer wg.Done()

	for taskItem := range spaceChan {
		q := url.Values{}
		q.Add("q", fmt.Sprintf("organization_guid:%s", taskItem.Variables.Org.Guid))
		taskItem.Metadata.Logger.Debugf("space query = %s", q)
		spaces, err := taskItem.Metadata.Client.ListSpacesByQuery(q)
		if err != nil {
			taskItem.Metadata.Logger.Errorf("error querying for org spaces '%s' (%s): %v",
				taskItem.Variables.Org.Name, taskItem.Variables.Org.Guid, err)
			continue
		}
		for _, space := range spaces {
			variables := task.Variables{
				Org:   taskItem.Variables.Org,
				Space: space,
			}
			isTrue, err := taskItem.Metadata.IsSpaceMatch(variables)
			if err != nil {
				taskItem.Metadata.Logger.Errorf("unable to determine true/false: %v", err)
				continue
			}
			if isTrue {
				taskItem.Metadata.Logger.Infof("Matched space '%s' in org '%s'", space.Name, variables.Org.Name)
				newTask := task.Item{
					Variables: variables,
					Metadata:  taskItem.Metadata,
				}
				appChan <- newTask
			} else {
				taskItem.Metadata.Logger.Infof("Skipping space '%s' in org '%s'", space.Name, variables.Org.Name)
			}
		}
	}

	logger.Error("exiting")
}
