package worker

import (
	"fmt"
	"net/url"
	"scullion/task"
)

func App(num int, appChan <-chan task.Item, actionChan chan<- task.Item) {
	fmt.Printf("Launched app worker %d\n", num)
	for {
		taskItem := <-appChan
		q := url.Values{
			"space_guid": []string{taskItem.Variables.Space.Guid},
		}
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
