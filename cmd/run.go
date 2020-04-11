package cmd

import (
	"os"
	"os/signal"
	"scullion/option"
	"scullion/task"
	"scullion/worker"
	"sync"
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

	var wg sync.WaitGroup
	for i := 0; i < cmd.OrgPool; i++ {
		wg.Add(1)
		go worker.Org(i, orgChan, spaceChan, &wg)
	}
	for i := 0; i < cmd.SpacePool; i++ {
		wg.Add(1)
		go worker.Space(i, spaceChan, appChan, &wg)
	}
	for i := 0; i < cmd.AppPool; i++ {
		wg.Add(1)
		go worker.App(i, appChan, actionChan, &wg)
	}
	for i := 0; i < cmd.ActionPool; i++ {
		wg.Add(1)
		go worker.Action(i, actionChan, &wg)
	}

	for i, task := range tasks {
		go worker.Task(i, task, client, orgChan, cmd.DryRun)
	}

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	// Begin cascade of shutting down...
	close(orgChan)
	wg.Wait()

	return nil
}
