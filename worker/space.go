package worker

import (
	"fmt"
	"net/url"
	"scullion/task"
	"scullion/util"

	"github.com/antonmedv/expr"
)

func Space(num int, spaceChan <-chan task.TaskItem, appChan chan<- task.TaskItem) {
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
			variables := task.TaskVariables{
				Org:   taskItem.Variables.Org,
				Space: space,
			}
			result, err := expr.Run(taskItem.Metadata.SpaceExpr, variables)
			if err != nil {
				panic(err)
			}
			isTrue, err := util.IsTrue(result)
			if err != nil {
				panic(err)
			}
			if isTrue {
				fmt.Printf("[%s] Matched space '%s' in org '%s'\n", taskItem.Metadata.Name, space.Name, variables.Org.Name)
				newTask := task.TaskItem{
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
