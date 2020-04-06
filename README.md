# Scullion

> Currently in development; things may move around. PR's more than welcome!

Cleans up after your Cloud Foundry development activities, so you don't have to.

# Usage

Scullion has multiple subcommands:
* `run` is intended to run the rules continuously against a foundation. This is the most used use case.
* `validate` allows validation of the task configuration; both the expression syntax as well as the object model (by way of sample data pulled from a running foundation).
* TODO: `one-time` executes the rules once.
* TOOD: `reference` dump out the object model to assist in scripting.

Both `run` and `one-time` have a dry-run mode to allow observation prior to taking action.

# Configuration

For each rule, it consists of:
* A name to identify what is triggering an action.
* A schedule (for run mode).
* Filters to select organizations, spaces, and applications.
* An action to take on applications that match.

## Actions

The following actions are allowed:
* `log`
* TODO: `stop`
* TODO: `delete`

Note that `stop` and `delete` are modified by the dry-run flag.

## Rules (tasks)

Scullion configuration consists of a number of tasks.  The general layout is done in JSON like this:
```
[
    {
        "name": "a sample",
        "schedule": {
            "frequency": "1h"
        },
        "filters": {
            "organization": "name ne 'system'",
            "space": "name eq 'test'",
            "application": "state eq 'STARTED' and age gt 3h",
            "action": "stop"
        }
    },
    {
      <snip>
    }
]
```

Configuration can be specified via an environment variable or a file. Most likely, with Cloud Foundry, this will be an via the manifest like this:
```
  env:
    SCULLION_CONFIG: >
      {
        "...": "..."
      }
```

# CLI snapshots

(From development version)

```
$ go run main.go --help
Usage:
  main [OPTIONS] <run | validate>

Application Options:
  -v, --verbose  Enable verbose output

Help Options:
  -h, --help     Show this help message

Available commands:
  run
  validate
```

```$ go run main.go validate --help
Usage:
  main [OPTIONS] validate [validate-OPTIONS]

Application Options:
  -v, --verbose   Enable verbose output

Help Options:
  -h, --help      Show this help message

[validate command options]

    Task Options:
      -e, --env=  Load configuration from environment variable (default: SCULLION_TASKS)
      -f, --file= Read configuration from given file
```

```
$ go run main.go r --help
Usage:
  main [OPTIONS] run [run-OPTIONS]

Application Options:
  -v, --verbose                     Enable verbose output

Help Options:
  -h, --help                        Show this help message

[run command options]

    Task Options:
      -e, --env=                    Load configuration from environment variable (default: SCULLION_TASKS)
      -f, --file=                   Read configuration from given file

    Worker Pools:
          --worker-org-pool=        Set the number of organization workers in the pool (default: 1) [$WORKER_ORG_POOL]
          --worker-space-pool=      Set the number of space workers in the pool (default: 1) [$WORKER_SPACE_POOL]
          --worker-app-pool=        Set the number of application workers in the pool (default: 1) [$WORKER_APP_POOL]
          --worker-action-pool=     Set the number of action (stop/start) workers in the pool (default: 1) [$WORKER_ACTION_POOL]

    Cloud Foundry Configuration:
      -a, --cf-api=                 API URL [$CF_API]
      -u, --cf-username=            Username [$CF_USERNAME]
      -p, --cf-password=            Password [$CF_PASSWORD]
      -k, --cf-skip-ssl-validation  Skip SSL validation of Cloud Foundry endpoint; not recommended [$CF_SKIP_SSL_VALIDATION]
```