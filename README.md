# Scullion

> Cleans up after your development activities.

## Configuration

Scullion configuration consists of a number of tasks.  The general layout is done in JSON like this:
```
[
    {
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

Configuration can be specified via an environment variable or a file. Most likely, with Cloud Foundry, this will be an environment file like this:
```
  env:
    SCULLION_CONFIG: >
      {
        "...": "..."
      }
```
