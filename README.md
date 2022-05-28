[![Go Report Card](https://goreportcard.com/badge/github.com/TrevorEdris/api-template)](https://goreportcard.com/report/github.com/TrevorEdris/api-template)
![CodeQL](https://github.com/TrevorEdris/api-template/workflows/CodeQL/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GoT](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://go.dev)

# API Template

This project provides an easy to modify template for an API using the labstack/echo framework.
The Repository code pattern is used for data access, making the API highly flexible for whatever
storage solution is necessary.

**Note:** A significant portion of the code has been slightly modified based on [mikestefanello/pagoda](https://github.com/mikestefanello/pagoda). 

## Local Development

For local development, use the `make (up|restart|down|logs)` commands provided by the Makefile.

```md
❯ make help
. . .
up            Run the API locally and print logs to stdout
down          Stop all containers
restart       Restart all containers
logs          Print logs in stdout
. . .
```

The above commands use `deployments/local/docker-compose.dev.yaml` to run the API. The binary for the API will be rebuilt automatically
when a change to one of the source `.go` files is detected (configurable in `.air.toml`).

```bash
app               |    ____    __
app               |   / __/___/ /  ___
app               |  / _// __/ _ \/ _ \
app               | /___/\__/_//_/\___/ v4.6.3
app               | High performance, minimalist Go web framework
app               | https://echo.labstack.com
app               | ____________________________________O/_______
app               |                                     O\
app               | ⇨ http server started on [::]:8000
```

See [https://github.com/cosmtrek/air](https://github.com/cosmtrek/air) for more details.

### Configuration

To configure parameters for the local instance of the API, copy the `sample.env` file into `.env`. The API uses [`joho/godotenv`](https://github.com/joho/godotenv) to read environment variables from this file and apply them to the container at runtime.
Once the environment variables are set, the API will then parse the environment variables using [`joeshaw/envdecode`](https://github.com/joeshaw/envdecode).

**Warning:** The `sample.env` file has values for `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`. These are
sensitive values that should _never_ be comitted to github. The provided [`blackbox`](https://github.com/StackExchange/blackbox) tool can be used to
encrypt sensitive files using a GPG key, such as `secrets/<region>/some-secret-file.gpg`. Only users listed
in `.blackbox/blackbox-admins.txt` whose GPG keys have been used to encrypt the listed `.gpg` files will
be able to decrypt those files.

## SLO

| Endpoint | Requests/s | p99  |
|---|---|---|
| `GET /` | 20000 | 5ms |
| `GET /item/:id` | 20000 | 5ms |
| `POST /item` | 20000 | 5ms |
| `PUT /item/:id` | 20000 | 5ms |
| `DELETE /item/:id` | 20000 | 5ms |

## Data Model

### Local

Local storage consists of a `sync.Map`, where the key is a string and the value is an `item.Model`, defined in `./app/model/item/model.go`.

### DynamoDB

DynamoDB storage consists of a single table, `items`, with the following definition:

```json
{
    "TableName": "items",
    "AttributeDefinitions": [
        {
            "Attributename": "id",
            "AttributeType": "S"
        }
    ],
    "KeySchema": [
        {
            "KeyType": "HASH",
            "AttributeName": "id"
        }
    ],
}
```

## Authentication

TODO: Impelement authentication

## Endpoints

TODO: Create auto-generated OpenAPI definition

## Deployment Procedure

Deployment of this API to kubernetes is made simpler with the following Makefile targets.

```md
❯ make help
. . .
build         Build and tag the docker container for the API
test          Run unit tests
finalize      Build, test, and tag the docker container with the finalized tag (typically, the full docker registery will be tagged here)
publish_only  Push the tagged docker image to the docker registry
publish       Finalize and publish the docker container
deploy_only   Fill out the .yaml.tmpl files and apply them to the specified namespace
deploy        Build, test, finalize, publish, and then deploy the docker container to kube
```

**Notes:**

* The default namespace is `tedris`.
* The default region is `dev`.

### Deploying to dev

1. `git checkout vX.Y.Z`
2. `make deploy NAMESPACE=<some_namespace_here>`

### Deploying to a non-dev environment

1. `git checkout vX.Y.Z`
2. `make deploy REGION=<some_non-dev_region>`

### Deploying _without_ running any tests

1. `git checkout vX.Y.Z`
2. `make deploy_only REGION=<some_region> NAMESPACE=<some_namespace>`

## Common Maintenance

_What common, repeated actions are necessary to ensure this API continues to run?_

_Do any API keys need rotated frequently? Does any data need to be deleted at some interval? etc._

## List of 3rd Party Libraries

The following is a list of all 3rd party libraries in use by this API

* _TODO_
