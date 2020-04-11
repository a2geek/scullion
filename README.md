# Scullion

[![Go Report Card](https://goreportcard.com/badge/github.com/a2geek/scullion)](https://goreportcard.com/report/github.com/a2geek/scullion)

> Currently in development; things may move around. PR's more than welcome!

Cleans up after your Cloud Foundry development activities, so you don't have to.

# TODO

Beyond what is noted elsewhere, additional items:
* Decide if the Golang model should be kept or switch to the JSON structrues as Cloud Foundry returns
* Need to resolve logging mechanism
* `panic` needs to be removed for robustness

# Usage

Scullion has multiple subcommands:

* `run` is intended to run the rules continuously against a foundation. This is the most used use case.
* `validate` allows validation of the task configuration; both the expression syntax as well as the object model (by way of sample data pulled from a running foundation).
* `one-time` executes the rules once.
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

## Filters

Filters can be applied at any of these levels:

* Organization
* Space
* Application

Each step in the filter hierarchy contains the results from the prior stels. Thus, the organization filter only has the organization. Space filters include both space detail as well as organization (if it was processed), etc.

TODO: If a rule only applies to spaces, the workers will begin processing against a space and skip organizations entirely. Same optimization applies to applications.

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
            "organization": "Org.name != 'system'",
            "space": "Space.name == 'test'",
            "application": "App.state == 'STARTED' && (Now() - Date(App.updated_at)) > Duration('1H')",
            "action": "stop"
        }
    },
    {
      <snip>
    }
]
```

(see `sample.json` for working sample)

Configuration can be specified via an environment variable or a file. Most likely, with Cloud Foundry, this will be an via the manifest like this:

```
  env:
    SCULLION_CONFIG: >
      {
        "...": "..."
      }
```

(see `manifest.yml` for working sample)

# CLI snapshots

(From development version)

```
$ go run main.go --help
Usage:
  main [OPTIONS] <command>

Application Options:
  -v, --verbose  Enable verbose output

Help Options:
  -h, --help     Show this help message

Available commands:
  disassemble
  one-time
  run
  validate
```

```
$ go run main.go disasm --help
Usage:
  main [OPTIONS] disassemble [disassemble-OPTIONS]

Application Options:
  -v, --verbose   Enable verbose output

Help Options:
  -h, --help      Show this help message

[disassemble command options]

    Task Options:
      -e, --env=  Load configuration from environment variable (default: SCULLION_TASKS)
      -f, --file= Read configuration from given file
```

```
$ go run main.go validate --help
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
$ go run main.go one-time --help
Usage:
  main [OPTIONS] one-time [one-time-OPTIONS]

Application Options:
  -v, --verbose                     Enable verbose output

Help Options:
  -h, --help                        Show this help message

[one-time command options]

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

```
$ go run main.go run --help
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