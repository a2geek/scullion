package cmd

import (
	"os"
	"os/signal"
	"scullion/ctx"
	"scullion/fn"
	"scullion/log"
	"scullion/option"
	"scullion/web"
	"scullion/worker"
	"sync"
	"syscall"
)

type Run struct {
	Port int `long:"port" default:"8080" env:"PORT" description:"Set the port number for the web server (0=off)"`

	option.RunOptions          `group:"Run Options"`
	option.TaskOptions         `group:"Task Options"`
	option.CloudFoundryOptions `group:"Cloud Foundry Configuration" namespace:"cf" env-namespace:"CF" reqired:"yes"`
}

func (cmd *Run) Execute(args []string) error {
	cfg, err := cmd.ReadConfiguration()
	if err != nil {
		return err
	}

	client, err := cmd.Client()
	if err != nil {
		panic(err)
	}

	registrars := []fn.Registrar{
		fn.NewCfCurlRegistrar(client),
		fn.NewDatetimeRegistrar(),
		fn.NewFiltererRegistrar(),
		fn.NewLibraryRegistrar(cfg.Library),
		fn.NewTemplateRegistrar(cfg.Templates),
	}

	logger, err := log.NewLogger("main", cmd.Level, cmd.NoDate)
	if err != nil {
		panic(err)
	}
	logger.Info("Started")

	stateChan := make(chan *ctx.State)

	var wg sync.WaitGroup
	for i := 0; i < cmd.WorkerPool; i++ {
		wg.Add(1)
		go worker.Task(i, stateChan, &wg, cmd.RunOptions)
		logger.Debugf("started task worker %d", i)
	}

	for i, rule := range cfg.Rules {
		go worker.Schedule(i, rule, registrars, stateChan, cmd.RunOptions)
		logger.Debugf("started schedule worker '%s'", rule.Name)
	}

	if cmd.Port > 0 {
		go web.Serve(cmd.Port, cfg)
	}

	logger.Info("Running...")
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	// Begin cascade of shutting down...
	logger.Info("Terminate request received. Shutting down...")
	close(stateChan)
	wg.Wait()
	logger.Info("Terminating. Bye!")

	return nil
}
