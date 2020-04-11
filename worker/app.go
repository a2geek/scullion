package worker

import (
	"fmt"
	"net/url"
	"scullion/task"
	"sync"
)

func App(num int, appChan <-chan task.Item, actionChan chan<- task.Item, wg *sync.WaitGroup) {
	fmt.Printf("Launched app worker %d\n", num)

	// if app channel closes, let's close the action channel as well!
	defer close(actionChan)
	defer wg.Done()

	for taskItem := range appChan {
		q := url.Values{}
		q.Add("q", fmt.Sprintf("space_guid:%s", taskItem.Variables.Space.Guid))
		apps, err := taskItem.Metadata.Client.ListAppsByQuery(q)
		if err != nil {
			panic(err)
		}
		for _, app := range apps {
			variables := task.Variables{
				Org:   taskItem.Variables.Org,
				Space: taskItem.Variables.Space,
				App:   app,
			}
			isTrue, err := taskItem.Metadata.IsAppMatch(variables)
			if err != nil {
				panic(err)
			}
			if isTrue {
				fmt.Printf("[%s] Matched app '%s' in space '%s' of org '%s'\n", taskItem.Metadata.Name, app.Name,
					variables.Space.Name, variables.Org.Name)
				newTask := task.Item{
					Variables: variables,
					Metadata:  taskItem.Metadata,
				}
				actionChan <- newTask
			} else {
				fmt.Printf("[%s] Skipped app '%s' in space '%s' of org '%s'\n", taskItem.Metadata.Name, app.Name,
					variables.Space.Name, variables.Org.Name)
			}
		}
	}
}
