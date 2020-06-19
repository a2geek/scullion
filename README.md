# Scullion

[![Go Report Card](https://goreportcard.com/badge/github.com/a2geek/scullion)](https://goreportcard.com/report/github.com/a2geek/scullion)

> Currently in development; things may move around. PR's more than welcome!

Cleans up after your Cloud Foundry development activities, so you don't have to.

## Processing hypothesis/overview

* Cloud Foundry makes deploying applications really easy.
* Developers don't like to clean up after themselves and CF tends to get cluttered.
* Go has interesting concurrency capabilities with goroutines and channels.
* By defining processing stages that emit results to the next stage dynamic rules can be defined.

## TODO

Beyond what is noted elsewhere, additional items:

* Add web component to include stats?

## Usage

Scullion has multiple sub-commands:

* `run` is intended to run the rules continuously against a foundation. This is the most used use case.
* `validate` allows validation of the task configuration; both the expression syntax as well as the object model (by way of sample data pulled from a running foundation).
* `one-time` executes the rules once.
* `reference` dump out the object model to assist in scripting.

Both `run` and `one-time` have a dry-run mode to allow observation prior to taking action.

## Web endpoints

`/health` always returns with success.

`/config` displays currently running configuration.

## Configuration

Scullion consists of a number of components:
* Rules,
* the Library (optional),
* Templates (optional).

### Rules

For each rule, it consists of a number of components.

* A name to identify what is triggering an action.
* A schedule (for run mode).
* A pipeline for data flow. This is the main set of processing. Each step _must_ emit results for processing to continue. This both enables multi-result capability and filtering.
* Actions. All actions are acted upon once a pipeline completes.

Yaml structure:

```
rules:
- name: the name used in logging goes here
  schedule:
    frequency: 1h   # See Go's ParseDuration for formats (`s`, `m`, `h` are the most useful)
  pipeline:         # Any failure stops the pipeline
  - # First stage
  - # Second stage
  - # Etc...
  actions:          # All actions are executed independently and likely concurrently
  - # Action 1
  - # Action 2
```

#### Pipelines

A pipeline originates all data and processes all data. It can be structured as needed. The base language is from [Expr](https://github.com/antonmedv/expr) with additional capabilities added to extend the language.

The following capabilities are exposed:
* Cloud Foundry API:
  * `GET(path, name)`: Retrieve one value. The entire response is stored in the variable identified by `name`.
  * `GetResources(path, name)`: Pages through a resource that contains a `resources` array. Each item in the `resources` array is emitted independently for the next step.
  * `POST(path, body)`: POST `body` to the given API endpoint.
  * `PUT(path, body)`: PUT `body` to the given API endpoint.
* Dates: (based on [`dates_test.go`](https://github.com/antonmedv/expr/blob/master/docs/examples/dates_test.go))
  * `Add(date.Time, date.Duration) date.Time` also mapped to `+`.
  * `After(date.Time, date.Time) bool` also mapped to `>`.
  * `AfterDuration(a, b time.Duration) bool` also mapped to `>`.
  * `AfterOrEqual(a, b time.Time) bool` also mapped to `>=`.
  * `AfterOrEqualDuration(a, b time.Duration) bool` also mapped to `>=`.
  * `Before(a, b time.Time) bool` also mapped to `<`.
  * `BeforeDuration(a, b time.Duration) bool` also mapped to `<`.
  * `BeforeOrEqual(a, b time.Time) bool` also mapped to `<=`.
  * `BeforeOrEqualDuration(a, b time.Duration) bool` also mapped to `<=`.
  * `Date(s string) time.Time` can be used to create a date based on a typical Cloud Foundry date format.
  * `Duration(s string) time.Duration` can be used to parse a duration (see Go [`ParseDuration`](https://golang.org/pkg/time/#ParseDuration) for formats supported).
  * `Equal(a, b time.Time) bool` also mapped to `=`.
  * `EqualDuration(a, b time.Duration) bool` also mapped to `=`.
  * `Now() time.Time` returns the current time.
  * `Since(s string) time.Duration` uses `Date` to parse the string and return `Duration` of that has passed.
  * `Sub(a, b time.Time) time.Duration` also mapped to `-`.
* Filters:
  * `Filter(expression)`: If true, emit current state. If false, processing stops.
* Libraries:
  * `Call(name)`: Call a subprogram in the library with current state.
* Templates:
  * `Template(name, parameters...)`: Allow templating, particularly for `POST(...)` or `PUT(...)` API calls. This is a very light wrapper around Go's `fmt.Sprintf(...)` and capabilities are described in [the package overview](https://golang.org/pkg/fmt/#pkg-overview).

### the Library

The library is reusable code. As an idea, the library could define common tasks (stop application or identify applications to process) or just give a name to a specific set of code. Use the `Call('library-name')` method to use it.  `Call` can be used both in the pipeline and in the library.

Yaml structure:

```
library:
- name: apps-to-consider
  pipeline:
  - # First stage
  - # Second stage
- name: stop-app
  pipeline:
  - # Etc...
```

### Templates

Templates are useful to extract strings, particularly large ones or strings that are easier to manage when formatted (like JSON structures). See the _Templates_ function reference for lookup capabilities.

Sample:

```
templates:
  StopApp: >
    {
      "state": "STOPPED"
    }
  SpaceDeveloper: >
    {
      "type": "space_developer", 
      "relationships": { 
        "user": {
          "data": {
            "guid": "%s"
          }
        },
        "space": {
          "data": {
            "guid": "%s"
          }
        }
      }
    }
```

### Putting it all together

TBD

(see `sample.json` for working sample)

Configuration can be specified via an environment variable or a file. Most likely, with Cloud Foundry, this will be an via the manifest like this:

```yaml
  env:
    SCULLION_CONFIG: >
      {
        "...": "..."
      }
```

(see `manifest.yml` for working sample)

## CLI snapshots

(From development version)

```shell
$ go run main.go --help
Usage:
  main [OPTIONS] <command>

Help Options:
  -h, --help  Show this help message

Available commands:
  disassemble
  one-time
  run
  validate
```

```shell
$ go run main.go disasm --help
Usage:
  main [OPTIONS] disassemble [disassemble-OPTIONS]

Help Options:
  -h, --help      Show this help message

[disassemble command options]

    Task Options:
      -e, --env=  Load configuration from environment variable (default: SCULLION_TASKS)
      -f, --file= Read configuration from given file
```

```shell
$ go run main.go validate --help
Usage:
  main [OPTIONS] validate [validate-OPTIONS]

Help Options:
  -h, --help      Show this help message

[validate command options]

    Task Options:
      -e, --env=  Load configuration from environment variable (default: SCULLION_TASKS)
      -f, --file= Read configuration from given file
```

```shell
$ go run main.go one-time --help
Usage:
  main [OPTIONS] one-time [one-time-OPTIONS]

Help Options:
  -h, --help                                     Show this help message

[one-time command options]

    Run Options:
          --dry-run                              Perform a dry run and log actions that would be taken [$SCULLION_DRY_RUN]
      -l, --log-level=[ERROR|WARNING|INFO|DEBUG] Set the logging level (default: INFO) [$SCULLION_LOG_LEVEL]
          --no-timestamp                         Suppress timestamp from logs (useful if other components add date) [$SCULLION_NO_TIMESTAMP]

    Task Options:
      -e, --env=                                 Load configuration from environment variable (default: SCULLION_TASKS)
      -f, --file=                                Read configuration from given file

    Worker Pools:
          --worker-org-pool=                     Set the number of organization workers in the pool (default: 1) [$WORKER_ORG_POOL]
          --worker-space-pool=                   Set the number of space workers in the pool (default: 1) [$WORKER_SPACE_POOL]
          --worker-app-pool=                     Set the number of application workers in the pool (default: 1) [$WORKER_APP_POOL]
          --worker-action-pool=                  Set the number of action (stop/start) workers in the pool (default: 1) [$WORKER_ACTION_POOL]

    Cloud Foundry Configuration:
      -a, --cf-api=                              API URL [$CF_API]
      -u, --cf-username=                         Username [$CF_USERNAME]
      -p, --cf-password=                         Password [$CF_PASSWORD]
      -k, --cf-skip-ssl-validation               Skip SSL validation of Cloud Foundry endpoint; not recommended [$CF_SKIP_SSL_VALIDATION]
```

```shell
$ go run main.go run --help
Usage:
  main [OPTIONS] run [run-OPTIONS]

Help Options:
  -h, --help                                     Show this help message

[run command options]
          --port=                                Set the port number for the web server (0=off) (default: 8080) [$PORT]

    Run Options:
          --dry-run                              Perform a dry run and log actions that would be taken [$SCULLION_DRY_RUN]
      -l, --log-level=[ERROR|WARNING|INFO|DEBUG] Set the logging level (default: INFO) [$SCULLION_LOG_LEVEL]
          --no-timestamp                         Suppress timestamp from logs (useful if other components add date) [$SCULLION_NO_TIMESTAMP]

    Task Options:
      -e, --env=                                 Load configuration from environment variable (default: SCULLION_TASKS)
      -f, --file=                                Read configuration from given file

    Worker Pools:
          --worker-org-pool=                     Set the number of organization workers in the pool (default: 1) [$WORKER_ORG_POOL]
          --worker-space-pool=                   Set the number of space workers in the pool (default: 1) [$WORKER_SPACE_POOL]
          --worker-app-pool=                     Set the number of application workers in the pool (default: 1) [$WORKER_APP_POOL]
          --worker-action-pool=                  Set the number of action (stop/start) workers in the pool (default: 1) [$WORKER_ACTION_POOL]

    Cloud Foundry Configuration:
      -a, --cf-api=                              API URL [$CF_API]
      -u, --cf-username=                         Username [$CF_USERNAME]
      -p, --cf-password=                         Password [$CF_PASSWORD]
      -k, --cf-skip-ssl-validation               Skip SSL validation of Cloud Foundry endpoint; not recommended [$CF_SKIP_SSL_VALIDATION]
```

```script
$ go run main.go reference
Note that methods prefixed with an operation usually have that operator overloaded.
Thereform, 'Add' for a Time and Duration can be expressed 'time + duration'.
This list is dynamically generated at run time, so should be accurate for your version.

Operations:
  Add(time.Time, time.Duration) time.Time
  After(time.Time, time.Time) bool
  AfterDuration(time.Duration, time.Duration) bool
  AfterOrEqual(time.Time, time.Time) bool
  AfterOrEqualDuration(time.Duration, time.Duration) bool
  Before(time.Time, time.Time) bool
  BeforeDuration(time.Duration, time.Duration) bool
  BeforeOrEqual(time.Time, time.Time) bool
  BeforeOrEqualDuration(time.Duration, time.Duration) bool
  Date(string) time.Time
  Duration(string) time.Duration
  Equal(time.Time, time.Time) bool
  EqualDuration(time.Duration, time.Duration) bool
  Now() time.Time
  Since(string) time.Duration
  Sub(time.Time, time.Time) time.Duration
```