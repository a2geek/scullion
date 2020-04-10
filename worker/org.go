package worker

import (
	"fmt"
	"scullion/task"
)

func Org(num int, orgChan <-chan task.Item, spaceChan chan<- task.Item) {
	fmt.Printf("Launched org worker %d\n", num)
	for {
		taskItem := <-orgChan
		orgs, err := taskItem.Metadata.Client.ListOrgs()
		if err != nil {
			panic(err)
		}
		for _, org := range orgs {
			variables := task.Variables{
				Org: org,
			}
			isTrue, err := taskItem.Metadata.IsOrgMatch(variables)
			if err != nil {
				panic(err)
			}
			if isTrue {
				fmt.Printf("[%s] Matched org '%s'\n", taskItem.Metadata.Name, org.Name)
				newTask := task.Item{
					Variables: variables,
					Metadata:  taskItem.Metadata,
				}
				spaceChan <- newTask
			} else {
				fmt.Printf("[%s] Skipping org '%s'\n", taskItem.Metadata.Name, org.Name)
			}
		}
	}
}
