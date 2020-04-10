package worker

import (
	"fmt"
	"net/url"
	"scullion/task"
)

func Space(num int, spaceChan <-chan task.Item, appChan chan<- task.Item) {
	fmt.Printf("Launched space worker %d\n", num)
	for {
		taskItem := <-spaceChan
		q := url.Values{
			"organization_guid": []string{taskItem.Variables.Org.Guid},
		}
		spaces, err := taskItem.Metadata.Client.ListSpacesByQuery(q)
		if err != nil {
			panic(err)
		}
		for _, space := range spaces {
			variables := task.Variables{
				Org:   taskItem.Variables.Org,
				Space: space,
			}
			isTrue, err := taskItem.Metadata.IsSpaceMatch(variables)
			if err != nil {
				panic(err)
			}
			if isTrue {
				fmt.Printf("[%s] Matched space '%s' in org '%s'\n", taskItem.Metadata.Name, space.Name, variables.Org.Name)
				newTask := task.Item{
					Variables: variables,
					Metadata:  taskItem.Metadata,
				}
				appChan <- newTask
			} else {
				fmt.Printf("[%s] Skipping space '%s' in org '%s'\n", taskItem.Metadata.Name, space.Name, variables.Org.Name)
			}
		}
	}
}
