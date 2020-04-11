package worker

import (
	"fmt"
	"net/url"
	"scullion/task"
	"sync"
)

func Space(num int, spaceChan <-chan task.Item, appChan chan<- task.Item, wg *sync.WaitGroup) {
	fmt.Printf("Launched space worker %d\n", num)

	// if space channel closes, let's close the app channel as well!
	defer close(appChan)
	defer wg.Done()

	for taskItem := range spaceChan {
		q := url.Values{}
		q.Add("q", fmt.Sprintf("organization_guid:%s", taskItem.Variables.Org.Guid))
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
