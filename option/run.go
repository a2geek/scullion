package option

type RunOptions struct {
	DryRun bool   `long:"dry-run" env:"SCULLION_DRY_RUN" description:"Perform a dry run and log actions that would be taken"`
	Level  string `default:"INFO" short:"l" env:"SCULLION_LOG_LEVEL" long:"log-level" description:"Set the logging level" choice:"ERROR" choice:"WARNING" choice:"INFO" choice:"DEBUG"`
	NoDate bool   `env:"SCULLION_NO_TIMESTAMP" long:"no-timestamp" description:"Suppress timestamp from logs (useful if other components add date)"`
}
