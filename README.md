# Maritime Ports Service

## Overview

This service is responsible for storing and querying maritime ports data.

## Getting Started

**Step 1.** Run HTTP server:

```shell script
make run-server
```

or run with `docker-compose`:

```shell script
docker-compose up -d
```

**Step 2.**: Click on http://0.0.0.0:8080/swagger/index.html to open Swagger UI.

**Step 3.**: Run following script to run some `curl` queries or just explore the API via Swagger UI:

```shell
./<project_root>/scripts/manual_api_test.sh
```

## Development Setup

**Step 0.** Install [pre-commit](https://pre-commit.com/):

```shell
pip install pre-commit

# For macOS users.
brew install pre-commit
```

Then run `pre-commit install` to setup git hook scripts.
Used hooks can be found [here](.pre-commit-config.yaml).

______________________________________________________________________

NOTE

> `pre-commit` aids in running checks (end of file fixing,
> markdown linting, go linting, runs go tests, json validation, etc.)
> before you perform your git commits.

______________________________________________________________________

**Step 1.** Install external tooling (golangci-lint, swag etc.):

```shell script
make install
```

**Step 2.** Setup project for local testing (code lint, runs tests, builds all needed binaries):

```shell script
make all
```

______________________________________________________________________

NOTE

> All binaries can be found in `<project_root>/bin` directory.
> Use `make clean` to delete old binaries.

______________________________________________________________________

**Step 3.** Run server:

```shell
make run-server
```

and open http://0.0.0.0:8080/swagger/index.html to open Swagger UI.

______________________________________________________________________

NOTE

> Check [Makefile](Makefile) for other useful commands.

______________________________________________________________________

______________________________________________________________________

NOTE

> After modification of the API remember to run `make docs` is order to re-generate
> the API documentation so that Swagger UI visualizes the up-to-date changes.

______________________________________________________________________

## Future Improvements

- Add more tests for negative scenarios.
- Introduce request and response logging and a `request_id` or `correlation_id` property to allow for easy tracing.
- Swap in-memory storage for a persistent storage type (RDBMS, Document Store, etc.).

## License

[MIT](LICENSE)
