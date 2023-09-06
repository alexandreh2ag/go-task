# go-task

`gtask` can manage scheduled and worker tasks.

## Config file

```yaml
workers:
  - id: "task1"
    command: "echo 'task 1'"
  - id: "task2"
    command: "echo 'task 2'"

scheduled:
  - id: "task1"
    expr: "*/5 * * * *"
    command: "echo 'task 1'"
  - id: "task2"
    expr: "0 12 * * *"
    command: "echo 'task 2'"
```

## Usage

```help
Usage:
  gtask [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  schedule    schedule sub commands
  validate    validate config
  version     Show version info
  worker      worker sub commands

```

### Workers

#### Generate Supervisord config

CLI options:
* user: Define user who run command
* working-dir: Define working dir
* group-name: Define group name (like prefix)
* output: Define path to save config file
* format: Choose format of output file

```shell
# this command will read gtask.yml and generate supervisord config with workers list.
gtask worker generate --config gtask.yml --group-name my-group --format supervisor --output dest/path.conf
```

### schedule

#### Run

CLI options:
* working-dir: Define working dir
* timezone: Choose a specific timezone
* no-result-print: Hide output of command
* result-path: Define path to save output of command

```shell
# this command will read gtask.yml and run scheduled tasks (based on cron expr). 
gtask schedule run --config gtask.yml
# or 
gtask schedule run --config gtask.yml --timezone 'Europe/Paris'
```

#### Start

CLI options:
* working-dir: Define working dir
* timezone: Choose a specific timezone
* no-result-print: Hide output of command
* result-path: Define path to save output of command
* tick: Select duration of each tick


```shell
# this command will read gtask.yml and start a daemon that scheduled tasks (based on cron expr). 
gtask schedule start --config gtask.yml
# or 
gtask schedule start --config gtask.yml --timezone 'Europe/Paris' --tick 10m
```

## Requirements

* golang (1.21+)

## Development


* Install dependencies:

  ```bash
  go mod download
  ```

* Generate mocks:

  ```bash
  ./bin/mock.sh
  ```

* Build the project:

  ```bash
  go build .
  ```

* Run tests:

  ```bash
  go test ./...
  ```

## License

MIT License, see [LICENSE](LICENSE.md).
