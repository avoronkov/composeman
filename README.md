# composeman - run docker-compose.yml using podman

## Installation

Using go:
```
$ go get github.com/avoronkov/composeman
```

## Usage

- Start all services (add `-d` to run in background):
```
$ composeman up
```

- Stop all services (add `-v` to remove anonymous volumes):
```
$ composeman down
```

- Specify list of services to start:
```
$ composeman up -d my-service
```
(Note: starting dependent services is not supported yet.)

## Supported docker-compose directives

- `image`

- `environment`

- `env_file`

- `ports`

- `volumes`

- `command`

- `build`

	* `context`

	* `target`

## TODO

- Handle "`build: args`" directive.

## Development notes

### Running tests

```
$ go test -coverprofile=coverage.out -coverpkg=github.com/avoronkov/composeman/... ./...
...
ok      github.com/avoronkov/composeman/tests   0.858s  coverage: 48.7% of statements in github.com/avoronkov/composeman/...
$ go tool cover -html=coverage.out
```
