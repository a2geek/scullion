package cmd

import (
	"os"
	"os/signal"
	"scullion/option"
	"scullion/task"
	"scullion/worker"
	"syscall"
)

type Run struct {
	option.TaskOptions         `group:"Task Options"`
	option.WorkerPools         `group:"Worker Pools" namespace:"worker" env-namespace:"WORKER"`
	option.CloudFoundryOptions `group:"Cloud Foundry Configuration" namespace:"cf" env-namespace:"CF" reqired:"yes"`
}

func (cmd *Run) Execute(args []string) error {
	tasks, err := cmd.ReadConfiguration()
	if err != nil {
		return err
	}

	client, err := cmd.Client()
	if err != nil {
		panic(err)
	}

	orgChan := make(chan task.Item)
	spaceChan := make(chan task.Item)
	appChan := make(chan task.Item)
	actionChan := make(chan task.Item)

	for i := 0; i < cmd.OrgPool; i++ {
		go worker.Org(i, orgChan, spaceChan)
	}
	for i := 0; i < cmd.SpacePool; i++ {
		go worker.Space(i, spaceChan, appChan)
	}
	for i := 0; i < cmd.AppPool; i++ {
		go worker.App(i, appChan, actionChan)
	}
	for i := 0; i < cmd.ActionPool; i++ {
		go worker.Action(i, actionChan)
	}

	for i, task := range tasks {
		go worker.Task(i, task, client, orgChan)
	}

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	return nil
}
