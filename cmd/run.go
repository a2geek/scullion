package cmd

import (
	"os"
	"os/signal"
	"scullion/log"
	"scullion/option"
	"scullion/task"
	"scullion/web"
	"scullion/worker"
	"sync"
	"syscall"
)

type Run struct {
	Port int `long:"port" default:"8080" env:"PORT" description:"Set the port number for the web server (0=off)"`

	option.RunOptions          `group:"Run Options"`
	option.TaskOptions         `group:"Task Options"`
	option.WorkerPools         `group:"Worker Pools" namespace:"worker" env-namespace:"WORKER"`
	option.CloudFoundryOptions `group:"Cloud Foundry Configuration" namespace:"cf" env-namespace:"CF" reqired:"yes"`
}

func (cmd *Run) Execute(args []string) error {
	taskDefs, err := cmd.ReadConfiguration()
	if err != nil {
		return err
	}

	client, err := cmd.Client()
	if err != nil {
		panic(err)
	}

	logger, err := log.NewLogger("main", cmd.Level, cmd.NoDate)
	if err != nil {
		panic(err)
	}
	logger.Info("Started")

	orgChan := make(chan task.Item)
	spaceChan := make(chan task.Item)
	appChan := make(chan task.Item)
	actionChan := make(chan task.Item)

	var wg sync.WaitGroup
	for i := 0; i < cmd.OrgPool; i++ {
		wg.Add(1)
		go worker.Org(i, orgChan, spaceChan, &wg, cmd.RunOptions)
		logger.Debugf("started org worker %d", i)
	}
	for i := 0; i < cmd.SpacePool; i++ {
		wg.Add(1)
		go worker.Space(i, spaceChan, appChan, &wg, cmd.RunOptions)
		logger.Debugf("started space worker %d", i)
	}
	for i := 0; i < cmd.AppPool; i++ {
		wg.Add(1)
		go worker.App(i, appChan, actionChan, &wg, cmd.RunOptions)
		logger.Debugf("started app worker %d", i)
	}
	for i := 0; i < cmd.ActionPool; i++ {
		wg.Add(1)
		go worker.Action(i, actionChan, &wg, cmd.RunOptions)
		logger.Debugf("started action worker %d", i)
	}

	for i, task := range taskDefs {
		go worker.Task(i, task, client, orgChan, cmd.RunOptions)
		logger.Debugf("started task worker '%s'", task.Name)
	}

	if cmd.Port > 0 {
		go web.Serve(cmd.Port, taskDefs)
	}

	logger.Info("Running...")
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	// Begin cascade of shutting down...
	logger.Info("Terminate request received. Shutting down...")
	close(orgChan)
	wg.Wait()
	logger.Info("Terminating. Bye!")

	return nil
}
