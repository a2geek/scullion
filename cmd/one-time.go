package cmd

import (
	"scullion/ctx"
	"scullion/fn"
	"scullion/log"
	"scullion/option"
	"scullion/worker"
	"sync"
)

type OneTime struct {
	option.RunOptions          `group:"Run Options"`
	option.TaskOptions         `group:"Task Options"`
	option.CloudFoundryOptions `group:"Cloud Foundry Configuration" namespace:"cf" env-namespace:"CF" reqired:"yes"`
}

func (cmd *OneTime) Execute(args []string) error {
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
		// fn.NewLibraryRegistrar(lib),
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

	// Cannot use Task directly as it has the timer embedded
	for _, ruleDef := range cfg.Rules {
		logrule, err := log.NewLogger(ruleDef.Name, cmd.Level, cmd.NoDate)
		if err != nil {
			panic(err)
		}

		subpgm, err := ctx.NewSubprogram(ruleDef.Name, ruleDef.Pipeline, ruleDef.Actions, logrule)
		if err != nil {
			panic(err)
		}

		state := ctx.NewState(subpgm.Dup(), stateChan)
		for _, registerFuncs := range registrars {
			registerFuncs(state)
		}
		stateChan <- state
	}

	// Begin cascade of shutting down...
	close(stateChan)
	wg.Wait()
	logger.Info("Terminating. Bye!")

	return nil
}
