package cmd

import (
	"scullion/action"
	"scullion/option"
	"scullion/task"
	"scullion/worker"
	"sync"
)

type OneTime struct {
	option.RunOptions          `group:"Run Options"`
	option.TaskOptions         `group:"Task Options"`
	option.WorkerPools         `group:"Worker Pools" namespace:"worker" env-namespace:"WORKER"`
	option.CloudFoundryOptions `group:"Cloud Foundry Configuration" namespace:"cf" env-namespace:"CF" reqired:"yes"`
}

func (cmd *OneTime) Execute(args []string) error {
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
		go worker.Org(i, orgChan, spaceChan, &wg, cmd.RunOptions)
	}
	for i := 0; i < cmd.SpacePool; i++ {
		wg.Add(1)
		go worker.Space(i, spaceChan, appChan, &wg, cmd.RunOptions)
	}
	for i := 0; i < cmd.AppPool; i++ {
		wg.Add(1)
		go worker.App(i, appChan, actionChan, &wg, cmd.RunOptions)
	}
	for i := 0; i < cmd.ActionPool; i++ {
		wg.Add(1)
		go worker.Action(i, actionChan, &wg, cmd.RunOptions)
	}

	// Cannot use Task directly as it has the timer embedded
	for _, taskDef := range tasks {
		actionFunc, err := action.NewActionFunc(taskDef.Filters.Action, cmd.DryRun)
		if err != nil {
			panic(err)
		}
		metadata, err := task.NewMetadata(taskDef, client, actionFunc, cmd.RunOptions)
		if err != nil {
			panic(err)
		}
		taskItem := task.Item{
			Metadata: metadata,
		}
		orgChan <- taskItem
	}

	// Begin cascade of shutting down...
	close(orgChan)
	wg.Wait()

	return nil
}
