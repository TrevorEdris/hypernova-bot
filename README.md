[![Go Report Card](https://goreportcard.com/badge/github.com/TrevorEdris/hypernova-bot)](https://goreportcard.com/report/github.com/TrevorEdris/hypernova-bot)
![CodeQL](https://github.com/TrevorEdris/hypernova-bot/workflows/CodeQL/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GoT](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://go.dev)

# HypernovaBot

This project is a Discord bot, enabling Discord+Minecraft server owners to provide their members
with a consistent economy system, rewarding them for being active in both the Discord server
as well as the Minecraft server.

## Local Development

For local development, use the `make (dev-up|dev-restart|dev-down|dev-logs)` commands provided by the Makefile.

```md
❯ make help
. . .
dev-up            Run the API locally and print logs to stdout
dev-down          Stop all containers
dev-restart       Restart all containers
dev-logs          Print logs in stdout
. . .
```

The above commands use `deployments/local/docker-compose.dev.yaml` to run the API and Bot. The binaries for each will be rebuilt automatically
when a change to one of the source `.go` files is detected (configurable in `.air.(api|bot).toml`).

```bash
bot               | {"level":"info","ts":1653848623.2347543,"caller":"controller/controller.go:46","msg":"Successfully opened discord session","session_id":"39d4923af2295bedc044979df6815077","bot_username":"pingpong-bot#1315"}
api               | running...
api               | {"time":"2022-05-29T18:23:43.6939776Z","level":"INFO","prefix":"echo","file":"container.go","line":"88","message":"Configured for local storage"}
api               | {"time":"2022-05-29T18:23:43.6945775Z","level":"INFO","prefix":"echo","file":"main.go","line":"48","message":"Starting HTTP server"}
api               |
api               |    ____    __
api               |   / __/___/ /  ___
api               |  / _// __/ _ \/ _ \
api               | /___/\__/_//_/\___/ v4.6.3
api               | High performance, minimalist Go web framework
api               | https://echo.labstack.com
api               | ____________________________________O/_______
api               |                                     O\
api               | ⇨ http server started on [::]:8000
bot               | {"level":"info","ts":1653848675.20739,"caller":"controller/controller.go:74","msg":"Handling new message event","author":"MuchUsername#5604"}
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

Local storage consists of a `sync.Map`, where the key is a string and the value is an `domain.Item`, defined in `./app/domain/model.go`.

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

TODO: Describe deployments

## Common Maintenance

_What common, repeated actions are necessary to ensure this API continues to run?_

_Do any API keys need rotated frequently? Does any data need to be deleted at some interval? etc._

## List of 3rd Party Libraries

The following is a list of all 3rd party libraries in use by this API

* _TODO_
