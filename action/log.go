package action

import (
	"fmt"
	"scullion/task"
)

func Log(taskItem task.Item) {
	fmt.Println(taskItem.Variables)
}
