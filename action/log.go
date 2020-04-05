package action

import (
	"fmt"
	"scullion/task"
)

func Log(taskItem task.TaskItem) {
	fmt.Println(taskItem.Variables)
}
