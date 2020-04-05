package worker

import (
	"fmt"
	"scullion/task"
	"scullion/util"

	"github.com/antonmedv/expr"
)

func Org(num int, orgChan <-chan task.TaskItem, spaceChan chan<- task.TaskItem) {
	fmt.Printf("Launched org worker %d\n", num)
	for {
		taskItem := <-orgChan
		orgs, err := taskItem.Metadata.Client.ListOrgs()
		if err != nil {
			panic(err)
		}
		for _, org := range orgs {
			variables := task.TaskVariables{
				Org: org,
			}
			result, err := expr.Run(taskItem.Metadata.OrgExpr, variables)
			if err != nil {
				panic(err)
			}
			isTrue, err := util.IsTrue(result)
			if err != nil {
				panic(err)
			}
			if isTrue {
				fmt.Printf("[%s] Matched org '%s'\n", taskItem.Metadata.Name, org.Name)
				newTask := task.TaskItem{
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