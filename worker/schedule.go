package worker

import (
	"scullion/config"
	"scullion/ctx"
	"scullion/fn"
	"scullion/log"
	"scullion/option"
	"time"
)

func Schedule(num int, ruleDef config.RuleDef, registrars []fn.Registrar, stateChan chan *ctx.State, runOptions option.RunOptions) {
	logger, err := log.NewLogger(ruleDef.Name, runOptions.Level, runOptions.NoDate)
	if err != nil {
		panic(err)
	}
	logger.Info("Launched")

	// should something unrecoverable occur, close the channel (this shuts down everything)
	defer close(stateChan)

	dur, err := time.ParseDuration(ruleDef.Schedule.Frequency)
	if err != nil {
		logger.Errorf("halting: %w", err)
		return
	}

	subpgm, err := ctx.NewSubprogram(ruleDef.Name, ruleDef.Pipeline, ruleDef.Actions, logger)
	if err != nil {
		logger.Errorf("halting: %w", err)
		return
	}

	for t := range time.Tick(dur) {
		logger.Infof("Tick at %s", t)
		state := ctx.NewState(subpgm.Dup(), stateChan)
		for _, registerFuncs := range registrars {
			registerFuncs(state)
		}
		stateChan <- state
	}

	logger.Error("exiting")
}
