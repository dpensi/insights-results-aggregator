# Insights Results Aggregator

[![GoDoc](https://godoc.org/github.com/RedHatInsights/insights-results-aggregator?status.svg)](https://godoc.org/github.com/RedHatInsights/insights-results-aggregator)
[![Go Report Card](https://goreportcard.com/badge/github.com/RedHatInsights/insights-results-aggregator)](https://goreportcard.com/report/github.com/RedHatInsights/insights-results-aggregator)
[![Build Status](https://travis-ci.org/RedHatInsights/insights-results-aggregator.svg?branch=master)](https://travis-ci.org/RedHatInsights/insights-results-aggregator)
[![codecov](https://codecov.io/gh/RedHatInsights/insights-results-aggregator/branch/master/graph/badge.svg)](https://codecov.io/gh/RedHatInsights/insights-results-aggregator)

Aggregator service for insights results

## Description

Insights Results Aggregator is a service that provides Insight OCP data that are being consumed by OpenShift Cluster Manager. That data contain information about clusters status (especially health, security, performance) based on results generated by Insights rules engine. Insights OCP data are consumed from selected broker, stored in a storage (that basically works as a cache) and exposed via REST API endpoints.

## Architecture

Aggregator service consists of three main parts:

1. Consumer that reads (consumes) Insights OCP messages from specified message broker. Usually Kafka broker is used but it might be possible to develop a interface for different broker. Insights OCP messages are basically encoded in JSON and contain results generated by rule engine.

1. HTTP or HTTPS server that exposes REST API endpoints that can be used to read list of organizations, list of clusters, read rules results for selected cluster etc. Additionally, basic metrics are exposed as well. Those metrics is configured to be consumed by Prometheus and visualized by Grafana.

1. Storage backend which is some instance of SQL database. Currently SQLite3 and PostgreSQL are fully supported, but more SQL databases might be added later.

### Whole data flow

![data_flow](./doc/customer-facing-services-architecture.png)

1. Event about new data from insights operator is consumed from Kafka. That event contains (among other things) URL to S3 Bucket
2. Insights operator data is read from S3 Bucket and insigts rules are applied to that data
3. Results (basically organization ID + cluster name + insights results JSON) are stored back into Kafka, but into different topic
4. That results are consumed by Insights rules aggregator service that caches them
5. The service provides such data via REST API to other tools, like OpenShift Cluster Manager web UI, OpenShift console, etc.

Please note that results are filtered - only results for organizations listed in `org_whitelist.csv` are processed and cached in aggregator.

### DB structure

#### Table report

This table is used as a cache for reports consumed from broker. Size of this
table (i.e. number of records) scales linearly with the number of clusters,
because only latest report for given cluster is stored (it is guarantied by DB
constraints). That table has defined compound key `org_id+cluster`,
additionally `cluster` name needs to be unique across all organizations.

```sql
CREATE TABLE report (
    org_id          INTEGER NOT NULL,
    cluster         VARCHAR NOT NULL UNIQUE,
    report          VARCHAR NOT NULL,
    reported_at     TIMESTAMP,
    last_checked_at TIMESTAMP,
    PRIMARY KEY(org_id, cluster)
)
```

#### Table cluster_rule_user_feedback

```sql
-- user_vote is user's vote, 
-- 0 is none,
-- 1 is like,
-- -1 is dislike
CREATE TABLE cluster_rule_user_feedback (
    cluster_id VARCHAR NOT NULL,
    rule_id INTEGER  NOT NULL,
    user_id VARCHAR NOT NULL,
    user_vote SMALLINT NOT NULL,
    added_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    message VARCHAR NOT NULL,

    PRIMARY KEY(cluster_id, rule_id, user_id)
)
```

## Documentation for developers

All packages developed in this project have documentation available on [GoDoc server](https://godoc.org/):

* [entry point to the service](https://godoc.org/github.com/RedHatInsights/insights-results-aggregator)
* [package `broker`](https://godoc.org/github.com/RedHatInsights/insights-results-aggregator/broker)
* [package `consumer`](https://godoc.org/github.com/RedHatInsights/insights-results-aggregator/consumer)
* [package `content`](https://godoc.org/github.com/RedHatInsights/insights-results-aggregator/content)
* [package `metrics`](https://godoc.org/github.com/RedHatInsights/insights-results-aggregator/metrics)
* [package `migration`](https://godoc.org/github.com/RedHatInsights/insights-results-aggregator/migration)
* [package `producer`](https://godoc.org/github.com/RedHatInsights/insights-results-aggregator/producer)
* [package `server`](https://godoc.org/github.com/RedHatInsights/insights-results-aggregator/server)
* [package `storage`](https://godoc.org/github.com/RedHatInsights/insights-results-aggregator/storage)
* [package `types`](https://godoc.org/github.com/RedHatInsights/insights-results-aggregator/types)

## Configuration

Configuration is done by toml config, default one is `config.toml` in working directory,
but it can be overwritten by `INSIGHTS_RESULTS_AGGREGATOR_CONFIG_FILE` env var.

Also each key in config can be overwritten by corresponding env var. For example if you have config

```toml
...
[storage]
db_driver = "sqlite3"
sqlite_datasource = "./aggregator.db"
pg_username = "user"
pg_password = "password"
pg_host = "localhost"
pg_port = 5432
pg_db_name = "aggregator"
pg_params = ""
...
```

and environment variables

```shell
INSIGHTS_RESULTS_AGGREGATOR__STORAGE__DB_DRIVER="postgres"
INSIGHTS_RESULTS_AGGREGATOR__STORAGE__PG_PASSWORD="your secret password"
```

the actual driver will be postgres with password "your secret password"

It's very useful for deploying docker containers and keeping some of your configuration
outside of main config file(like passwords).

## Server configuration

Server configuration is in section `[server]` in config file.

```toml
[server]
address = ":8080"
api_prefix = "/api/v1/"
api_spec_file = "openapi.json"
debug = true
auth = true
auth_type = "xrh"
```

* `address` is host and port which server should listen to
* `api_prefix` is prefix for RestAPI path
* `api_spec_file` is the location of a required OpenAPI specifications file
* `debug` is developer mode that enables some special API endpoints not used on production
* `auth` turns on or turns authentication
* `auth_type` set type of auth, it means which header to use for auth `x-rh-identity` or `Authorization`. Can be used only with `auth = true`. Possible options: `jwt`, `xrh`

## Local setup

There is a `docker-compose` configuration that provisions a minimal stack of Insight Platform and
a postgres database.
You can download it here <https://gitlab.cee.redhat.com/insights-qe/iqe-ccx-plugin/blob/master/docker-compose.yml>

### Prerequisites

* minio requires `../minio/data/` and `../minio/config` directories to be created
* edit localhost line in your `/etc/hosts`:  `127.0.0.1       localhost kafka minio`
* `ingress` image should present on your machine. You can build it locally from this repo <https://github.com/RedHatInsights/insights-ingress-go>

### Usage

1. Start the stack `podman-compose up` or `docker-compose up`
2. Wait until kafka will be up.
3. Start `ccx-data-pipeline`: `python3 -m insights_messaging config-devel.yaml`
4. Build `insights-results-aggregator`: `make build`
5. Start `insights-results-aggregator`: `INSIGHTS_RESULTS_AGGREGATOR_CONFIG_FILE=config-devel.toml ./insights-results-aggregator`

Stop Minimal Insights Platform stack `podman-compose down` or `docker-compose down`

In order to upload an insights archive, you can use `curl`:

```shell
curl -k -vvvv -F "upload=@/path/to/your/archive.zip;type=application/vnd.redhat.testareno.archive+zip" http://localhost:3000/api/ingress/v1/upload -H "x-rh-identity: eyJpZGVudGl0eSI6IHsiYWNjb3VudF9udW1iZXIiOiAiMDAwMDAwMSIsICJpbnRlcm5hbCI6IHsib3JnX2lkIjogIjEifX19Cg=="
```

or you can use integration tests suite. More details are [here](https://gitlab.cee.redhat.com/insights-qe/iqe-ccx-plugin).

### Kafka producer

It is possible to use the script `produce_insights_results` from `utils` to produce several Insights results into Kafka topic. Its dependency is Kafkacat that needs to be installed on the same machine. You can find installation instructions [on this page](https://github.com/edenhill/kafkacat).

## Database

Aggregator is configured to use SQLite3 DB by default, but it also supports PostgreSQL.
In CI and QA environments, the configuration is overridden by environment variables to use PostgreSQL.

To establish connection to PostgreSQL, the following configuration options need to be changed in `storage` section of `config.toml`:

```toml
[storage]
db_driver = "postgres"
pg_username = "postgres"
pg_password = "postgres"
pg_host = "localhost"
pg_port = 5432
pg_db_name = "controller"
pg_params = "sslmode=disable"
```

### Migration mechanism

This service contains an implementation of a simple database migration mechanism that allows semi-automatic transitions between various database versions as well as building the latest version of the database from scratch.

Before using the migration mechanism, it is first necessary to initialize the migration information table `migration_info`. This can be done using the `migration.InitInfoTable(*sql.DB)` function. Any attempt to get or set the database version without initializing this table first will result in a `no such table: migration_info` error from the SQL driver.

New migrations must be added manually into the code, because it was decided that modifying the list of migrations at runtime is undesirable.

To migrate the database to a certain version, in either direction (both upgrade and downgrade), use the `migration.SetDBVersion(*sql.DB, migration.Version)` function.

**To upgrade the database to the highest available version, use `migration.SetDBVersion(db, migration.GetMaxVersion())`.** This will automatically perform all the necessary steps to migrate the database from its current version to the highest defined version.

See `/migration/migration.go` documentation for an overview of all available DB migration functionality.

## REST API schema based on OpenAPI 3.0

Aggregator service provides information about its REST API schema via endpoint `api/v1/openapi.json`. OpenAPI 3.0 is used to describe the schema; it can be read by human and consumed by computers.

For example, if aggregator is started locally, it is possible to read schema based on OpenAPI 3.0 specification by using the following command:

```
curl localhost:8080/api/v1/openapi.json
```

Please note that OpenAPI schema is accessible w/o the need to provide authorization tokens.

## Prometheus API

It is possible to use `/api/v1/metrics` REST API endpoint to read all metrics exposed to Prometheus or to any tool that is compatible with it.
Currently, the following metrics are exposed:

1. `api_endpoints_requests` the total number of requests per endpoint
1. `api_endpoints_response_time` API endpoints response time
1. `consumed_messages` the total number of messages consumed from Kafka
1. `feedback_on_rules` the total number of left feedback
1. `produced_messages` the total number of produced messages
1. `written_reports` the total number of reports written to the storage

Additionally it is possible to consume all metrics provided by Go runtime. There metrics start with `go_` and `process_` prefixes.

## Contribution

Please look into document [CONTRIBUTING.md](CONTRIBUTING.md) that contains all information about how to contribute to this project.

Please look also at [Definitiot of Done](DoD.md) document with further informations.

## Testing

The following tests can be run to test your code in `insights-results-aggregator`. Detailed information about each type of test is included in the corresponding subsection:

1. Unit tests: checks behaviour of all units in source code (methods, functions)
1. REST API Tests: test the real REST API of locally deployed application with database initialized with test data only
1. Integration tests: the integration tests for `insights-results-aggregator` service
1. Metrics tests: test whether Prometheus metrics are exposed as expected

### Unit tests

Set of unit tests checks all units of source code. Additionally the code coverage is computed and displayed.
Code coverage is stored in a file `coverage.out` and can be checked by a script named `check_coverage.sh`.

To run unit tests use the following command:

`make test`

### All integration tests

`make integration_tests`

#### Only REST API tests

Set of tests to check REST API of locally deployed application with database initialized with test data only.

To run REST API tests use the following command:

`make rest_api_tests`

#### Only metrics tests

`make metrics_tests`

## CI

[Travis CI](https://travis-ci.com/) is configured for this repository. Several tests and checks are started for all pull requests:

* Unit tests that use the standard tool `go test`.
* `go fmt` tool to check code formatting. That tool is run with `-s` flag to perform [following transformations](https://golang.org/cmd/gofmt/#hdr-The_simplify_command)
* `go vet` to report likely mistakes in source code, for example suspicious constructs, such as Printf calls whose arguments do not align with the format string.
* `golint` as a linter for all Go sources stored in this repository
* `gocyclo` to report all functions and methods with too high cyclomatic complexity. The cyclomatic complexity of a function is calculated according to the following rules: 1 is the base complexity of a function +1 for each 'if', 'for', 'case', '&&' or '||' Go Report Card warns on functions with cyclomatic complexity > 9
* `goconst` to find repeated strings that could be replaced by a constant
* `gosec` to inspect source code for security problems by scanning the Go AST
* `ineffassign` to detect and print all ineffectual assignments in Go code
* `errcheck` for checking for all unchecked errors in go programs
* `shellcheck` to perform static analysis for all shell scripts used in this repository
* `abcgo` to measure ABC metrics for Go source code and check if the metrics does not exceed specified threshold

Please note that all checks mentioned above have to pass for the change to be merged into master branch.

History of checks performed by CI is available at [RedHatInsights / insights-results-aggregator](https://travis-ci.org/RedHatInsights/insights-results-aggregator).
