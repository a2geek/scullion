package worker

import (
	"fmt"
	"scullion/ctx"
	"scullion/log"
	"scullion/option"
	"sync"
)

func Task(num int, stateChan chan *ctx.State, wg *sync.WaitGroup, runOpts option.RunOptions) {
	logger, err := log.NewLogger(fmt.Sprintf("task worker %d", num), runOpts.Level, runOpts.NoDate)
	if err != nil {
		panic(err)
	}
	logger.Info("Launched")

	// if ctx channel closes, be certain we clean up the wait group!
	defer wg.Done()

	for state := range stateChan {
		state.Execute()
	}

	logger.Error("exiting")
}
