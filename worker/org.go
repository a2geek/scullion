package worker

import (
	"fmt"
	"scullion/log"
	"scullion/task"
	"sync"
)

func Org(num int, orgChan <-chan task.Item, spaceChan chan<- task.Item, wg *sync.WaitGroup, logLevel string) {
	logger, err := log.NewLogger(fmt.Sprintf("org worker %d", num), logLevel)
	if err != nil {
		panic(err)
	}
	logger.Info("Launched")

	// if org channel closes, let's close the space channel as well!
	defer close(spaceChan)
	defer wg.Done()

	for taskItem := range orgChan {
		orgs, err := taskItem.Metadata.Client.ListOrgs()
		if err != nil {
			taskItem.Metadata.Logger.Errorf("error querying for orgs: %v", err)
			continue
		}
		for _, org := range orgs {
			variables := task.Variables{
				Org: org,
			}
			isTrue, err := taskItem.Metadata.IsOrgMatch(variables)
			if err != nil {
				taskItem.Metadata.Logger.Errorf("unable to determine true/false: %v", err)
				continue
			}
			if isTrue {
				taskItem.Metadata.Logger.Infof("Matched org '%s'", org.Name)
				newTask := task.Item{
					Variables: variables,
					Metadata:  taskItem.Metadata,
				}
				spaceChan <- newTask
			} else {
				taskItem.Metadata.Logger.Infof("Skipping org '%s'", org.Name)
			}
		}
	}

	logger.Error("exiting")
}
