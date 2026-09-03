<a name="readme-top"></a>

<!-- PROJECT SHIELDS -->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![Apache-2.0 License][license-shield]][license-url]

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/capcom6/service-monitor-tgbot">
    <img src="assets/logo.png" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">Service Monitor Telegram Bot</h3>

  <p align="center">
    Telegram bot for monitoring the availability of network services.
    <br />
    <a href="https://github.com/capcom6/service-monitor-tgbot/issues">Report Bug</a>
    ·
    <a href="https://github.com/capcom6/service-monitor-tgbot/issues">Request Feature</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
- [About The Project](#about-the-project)
  - [Built With](#built-with)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
  - [Messages Template System](#messages-template-system)
  - [Commands](#commands)
  - [Heartbeat](#heartbeat)
  - [API Documentation](#api-documentation)
- [Examples](#examples)
  - [HTTP service monitoring example](#http-service-monitoring-example)
  - [TCP service monitoring example](#tcp-service-monitoring-example)
- [Configuration](#configuration)
  - [Environment Variables](#environment-variables)
  - [YAML Config Structure](#yaml-config-structure)
  - [Service Definition Fields](#service-definition-fields)
  - [Storage Backends](#storage-backends)
- [Deployment](#deployment)
  - [Docker](#docker)
  - [GoReleaser](#goreleaser)
- [Development](#development)
  - [Prerequisites](#prerequisites-1)
  - [Make Targets](#make-targets)
  - [Quick Start](#quick-start)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)
- [Acknowledgments](#acknowledgments)

<!-- ABOUT THE PROJECT -->
## About The Project

Monitoring the availability of network services is an important task for any project. At the same time, it is not always necessary to deploy universal solutions like Prometheus — a simpler solution often suffices. It is for such cases that this bot was created.

The bot monitors the availability of HTTP(S) and TCP services and notifies a Telegram channel or group about changes in their status. It also supports on-demand status queries via commands and periodic heartbeat summaries to confirm the bot is alive.

**Key features:**

- HTTP(S) and TCP service probing with configurable thresholds
- Flexible notification message templates (Go template syntax)
- Pluggable storage backends: YAML file or Redis
- `/status` command for querying current service states
- Periodic heartbeat messages
- Telegram SOCKS5 proxy support
- OpenAPI/Swagger documentation
- Prometheus metrics endpoint

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Built With

* [![Go][Golang]][Golang-url]
* [![Fiber][Fiber]][Fiber-url]
* [![Redis][Redis]][Redis-url]
* [![Telegram Bot API][Telegram]][Telegram-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
## Getting Started

Follow the instructions below to run the bot.

### Prerequisites

Choose one of the following:

- **Docker** (recommended) — no other dependencies required
- **Go 1.26+** — for building from source

### Installation

1. Register a new bot and get a token: https://core.telegram.org/bots/features#creating-a-new-bot
2. Create a [channel](https://telegram.org/tour/channels) or [group](https://telegram.org/tour/groups) for notifications
3. Add the bot as an administrator with permission to send messages
4. Copy [config.example.yml](configs/config.example.yml) to your working directory as `config.yml`
5. Edit the configuration file:
    - Set the bot token
    - Set the channel/group ID (find it via `https://api.telegram.org/bot<token>/getUpdates?allowed_updates=[]` after adding the bot — look for `my_chat_member.chat.id`)
    - List the services to monitor (or configure a [storage backend](#storage-backends))
6. Run the bot:

```bash
docker run -d \
  -v "$(pwd)/config.yml:/app/config.yml:ro" \
  --name tgbot \
  ghcr.io/capcom6/service-monitor-tgbot:latest
```

Alternatively, copy [`.env.example`](.env.example) to `.env` and set environment variables there.

**From source:**

```bash
make deps
make build
./bin/service-monitor-tgbot
```

> **Note:** By default, the service list is read from the same YAML file (file storage backend). To use Redis or another backend, see the [Storage Backends](#storage-backends) section and set the `storage.dsn` field or `STORAGE__DSN` environment variable.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USAGE -->
## Usage

### Messages Template System

The bot uses a flexible template system for customizing notification messages. Templates support Go `text/template` syntax and a custom `escape` function for Telegram MarkdownV2 escaping. Templates can be customized in the configuration file under the `telegram.messages` key or via the `TELEGRAM__MESSAGES` environment variable (JSON object).

| Template        | Variables                                                                                                                                                                                                                                      | Description                           |
| --------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------- |
| `online`        | `Name` - name of the service<br>`ChangedAt` - when the service entered this state<br>`Duration` - how long the service has been in this state                                                                                                  | message when a service goes "online"  |
| `offline`       | `Name` - name of the service<br>`Error` - error message<br>`ChangedAt` - when the service entered this state<br>`Duration` - how long the service has been in this state                                                                       | message when a service goes "offline" |
| `services_list` | list of all services, each with:<br>`Name` - name of the service<br>`State` - state of the service<br>`Error` - error message<br>`ChangedAt` - when the service entered this state<br>`Duration` - how long the service has been in this state | message with a list of all services   |
| `heartbeat`     | `TotalServices` - total number of services<br>`OnlineServices` - number of online services<br>`OfflineServices` - number of offline services<br>`CheckedAt` - time of the heartbeat check                                                      | periodic heartbeat status summary     |

### Commands

- `/status` — Get the current status of all monitored services

### Heartbeat

The heartbeat feature sends periodic status summaries to help distinguish between "everything is fine" and "the bot has crashed."

When enabled, the bot sends a summary message at a configurable interval regardless of service states. The summary shows how many monitored services are online and offline. By default, an offline count is appended when at least one service is offline.

Configure via environment variables or YAML:

```bash
HEARTBEAT__ENABLED=true
HEARTBEAT__INTERVAL=1h
HEARTBEAT__CHATID=-1001234567890  # optional, defaults to TELEGRAM__CHATID
```

Or in YAML:

```yaml
heartbeat:
  enabled: true
  interval: 1h
  chatId: -1001234567890  # optional, defaults to telegram.chatId
```

Default heartbeat message: `💓 Heartbeat: 5/5 services online` (with the `(1 offline)` suffix when any service is offline). The message text can be customized via the `heartbeat` template.

### API Documentation

When OpenAPI is enabled (default), Swagger documentation is available at:

```
http://localhost:3000/api/v1/docs
```

Disable it with `HTTP__OPENAPI__ENABLED=false` or the `http.openapi.enabled: false` config field.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- EXAMPLES -->
## Examples

### HTTP service monitoring example

```yaml
services:
  - name: Google
    initialDelaySeconds: 5
    periodSeconds: 10
    timeoutSeconds: 1
    successThreshold: 1
    failureThreshold: 3
    httpGet:
      scheme: https
      host: google.com
      path: /
      port: 443
      httpHeaders:
        - name: X-Header
          value: value
```

![HTTP Alert][http-alert]

### TCP service monitoring example

```yaml
services:
  - name: MySQL
    initialDelaySeconds: 5
    periodSeconds: 10
    timeoutSeconds: 1
    successThreshold: 1
    failureThreshold: 3
    tcpSocket:
      host: localhost
      port: 3306
```

![TCP Alert][tcp-alert]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONFIGURATION -->
## Configuration

Configuration is loaded in order: defaults < YAML file < `.env` file < environment variables. Set `CONFIG_PATH` to point to your YAML config file.

### Environment Variables

| Variable                     | Description                                                | Default                   | Example                             |
| ---------------------------- | ---------------------------------------------------------- | ------------------------- | ----------------------------------- |
| `CONFIG_PATH`                | Path to main YAML configuration file                       | —                         | `./configs/config.yml`              |
| `DEBUG`                      | Enable development logging (human-readable console output) | `false`                   | `1`                                 |
| `HTTP__ADDRESS`              | HTTP server bind address (health/metrics)                  | `127.0.0.1:3000`          | `0.0.0.0:3000`                      |
| `HTTP__PROXY_HEADER`         | Trusted proxy header for real client IP                    | `X-Forwarded-For`         | `X-Real-IP`                         |
| `HTTP__PROXIES`              | Trusted proxy IPs/CIDRs (comma-separated)                  | (empty)                   | `10.0.0.0/8,192.168.1.0/24`         |
| `HTTP__OPENAPI__ENABLED`     | Enable Swagger/OpenAPI docs endpoint                       | `true`                    | `false`                             |
| `HTTP__OPENAPI__PUBLIC_HOST` | Public hostname in OpenAPI spec                            | (empty)                   | `api.example.com`                   |
| `HTTP__OPENAPI__PUBLIC_PATH` | Path prefix for OpenAPI docs                               | (empty)                   | `/docs`                             |
| `TELEGRAM__TOKEN`            | Telegram Bot API token (required)                          | —                         | `123456:ABC-DEF...`                 |
| `TELEGRAM__CHATID`           | Target chat/group/channel ID (required)                    | —                         | `-1001234567890`                    |
| `TELEGRAM__PROXY_URL`        | SOCKS5 proxy for Telegram API                              | (empty)                   | `socks5://user:pass@127.0.0.1:1080` |
| `TELEGRAM__TIMEOUT`          | Telegram API client timeout                                | `1m`                      | `30s`                               |
| `TELEGRAM__MESSAGES`         | Custom message templates (JSON object)                     | (built-in defaults)       | `{"online":"✅ {{.Name}}"}`          |
| `HEARTBEAT__ENABLED`         | Enable periodic heartbeat messages                         | `false`                   | `true`                              |
| `HEARTBEAT__INTERVAL`        | Heartbeat send interval                                    | `6h`                      | `1h`                                |
| `HEARTBEAT__CHATID`          | Target chat for heartbeat messages                         | (uses `TELEGRAM__CHATID`) | `-1001234567890`                    |
| `STORAGE__DSN`               | Storage backend DSN                                        | `file://$CONFIG_PATH`     | `redis://localhost:6379/0`          |

> **Note:** Environment variable names use `__` (double underscore) as the separator for nested config keys (e.g., `TELEGRAM__TOKEN` maps to `telegram.token`). A `.env` file in the working directory is loaded automatically.

### YAML Config Structure

```yaml
telegram:
  token: 123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11
  chatId: -1234567890123
  proxyUrl: socks5://user:pass@127.0.0.1:1080  # optional
  messages: {}                                 # optional, overrides built-in templates

heartbeat:
  enabled: true   # optional, enables periodic heartbeat messages
  interval: 1h    # optional, default: 6h
  chatId: 0       # optional, 0 = use telegram.chatId

storage:
  dsn: file://./configs/services.yml  # optional, default: file://$CONFIG_PATH

services:  # optional inline list of services, see Examples and Service Definition Fields
```

### Service Definition Fields

| Field                 | Type   | Default       | Description                                                  |
| --------------------- | ------ | ------------- | ------------------------------------------------------------ |
| `name`                | string | (required)    | Human-readable service name                                  |
| `initialDelaySeconds` | int16  | `0`           | Delay before first check; negative = random 0..periodSeconds |
| `periodSeconds`       | uint16 | `10`          | Seconds between probes                                       |
| `timeoutSeconds`      | uint16 | `1`           | Probe timeout                                                |
| `successThreshold`    | uint8  | `1`           | Consecutive successes to mark "online"                       |
| `failureThreshold`    | uint8  | `3`           | Consecutive failures to mark "offline"                       |
| `httpGet`             | object |               | HTTP probe (mutually exclusive with `tcpSocket`)             |
| `httpGet.scheme`      | string | `http`        | `http` or `https`                                            |
| `httpGet.host`        | string | (required)    | Hostname                                                     |
| `httpGet.path`        | string | `/`           | URL path                                                     |
| `httpGet.port`        | uint16 | auto (80/443) | Port (defaults based on scheme)                              |
| `httpGet.httpHeaders` | list   |               | Custom HTTP headers                                          |
| `tcpSocket`           | object |               | TCP probe (mutually exclusive with `httpGet`)                |
| `tcpSocket.host`      | string | (required)    | Hostname                                                     |
| `tcpSocket.port`      | uint16 | (required)    | Port                                                         |

### Storage Backends

The bot supports pluggable storage backends for loading the list of monitored services. The backend is selected via DSN in the `storage.dsn` config field or `STORAGE__DSN` environment variable.

| Scheme  | Backend             | Use Case                                      |
| ------- | ------------------- | --------------------------------------------- |
| `file`  | YAML file (default) | Simple setups, static configs                 |
| `redis` | Redis key           | Dynamic service lists, distributed management |

**File DSN format:** `file:///path/to/services.yml`

**Redis DSN format:** `redis://[[user]:pass@]host:port[/db][?key=<key>&channel=<channel>]`

Default key: `service-monitor:services`, default channel: `service-monitor:reload`.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- DEPLOYMENT -->
## Deployment

### Docker

The recommended way to run the bot:

```bash
docker run -d \
  -v "$(pwd)/config.yml:/app/config.yml:ro" \
  --name tgbot \
  ghcr.io/capcom6/service-monitor-tgbot:latest
```

### GoReleaser

The project uses GoReleaser for building and publishing releases. Binaries are built for Linux, macOS, and Windows. Docker images are published for `linux/amd64` and `linux/arm64`.

```bash
make release
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- DEVELOPMENT -->
## Development

### Prerequisites

- Go 1.26+
- [golangci-lint](https://golangci-lint.run/) v2
- [Air](https://github.com/air-verse/air) (for hot-reload, optional)
- [GoReleaser](https://goreleaser.com/) (for release builds, optional)

### Make Targets

| Target              | Description                                                               |
| ------------------- | ------------------------------------------------------------------------- |
| `make all`          | Format + lint + test coverage                                             |
| `make build`        | Build binary to `bin/`                                                    |
| `make test`         | Run tests with race detection and coverage                                |
| `make coverage`     | Generate coverage report (`coverage.out` + `coverage.html`)               |
| `make lint`         | Run golangci-lint                                                         |
| `make fmt`          | Format code                                                               |
| `make air`          | Development server with hot-reload                                        |
| `make release`      | GoReleaser snapshot build                                                 |
| `make docker-build` | Build Docker image                                                        |
| `make docker-up`    | Start services via Docker Compose (requires a local `docker-compose.yml`) |
| `make docker-down`  | Stop Docker Compose services                                              |
| `make swagger`      | Regenerate Swagger docs                                                   |
| `make clean`        | Clean build artifacts                                                     |
| `make help`         | Show available targets                                                    |

### Quick Start

```bash
make deps        # install Go dependencies
make air         # start with hot-reload (requires Air)
```

Or build and run directly:

```bash
make build
CONFIG_PATH=configs/config.yml DEBUG=1 ./bin/service-monitor-tgbot
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- ROADMAP -->
## Roadmap

- [x] Add Changelog
- [x] Add the ability to change the text of messages
- [x] Pluggable storage backends (file, Redis)
- [x] Proxy support for Telegram client
- [x] Request current state of services (`/status` command)
- [x] Periodic heartbeat messages
- [ ] Send notifications to multiple channels/groups
- [ ] Display event time in notifications
- [ ] Online/offline time count
- [ ] Active bot mode
     - [ ] SLA report
     - [ ] The event log
- [ ] Separation of bot and monitoring service
- [ ] Dynamic list of services
- [ ] Service discovery
- [ ] Resource groups

See the [open issues](https://github.com/capcom6/service-monitor-tgbot/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- LICENSE -->
## License

Distributed under the Apache-2.0 license. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTACT -->
## Contact

**API Support:** i@capcom.me

Project Link: [https://github.com/capcom6/service-monitor-tgbot](https://github.com/capcom6/service-monitor-tgbot)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

* [Go Fiber](https://gofiber.io/) — Express-inspired web framework
* [Uber FX](https://uber-go.github.io/fx/) — Dependency injection framework
* [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) — Telegram Bot API client
* [go-redis](https://github.com/redis/go-redis) — Redis client for Go
* [Best README Template](https://github.com/othneildrew/Best-README-Template) — README structure inspiration
* [Keep a Changelog](https://keepachangelog.com/) — Changelog format

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
[contributors-shield]: https://img.shields.io/github/contributors/capcom6/service-monitor-tgbot.svg?style=for-the-badge
[contributors-url]: https://github.com/capcom6/service-monitor-tgbot/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/capcom6/service-monitor-tgbot.svg?style=for-the-badge
[forks-url]: https://github.com/capcom6/service-monitor-tgbot/network/members
[stars-shield]: https://img.shields.io/github/stars/capcom6/service-monitor-tgbot.svg?style=for-the-badge
[stars-url]: https://github.com/capcom6/service-monitor-tgbot/stargazers
[issues-shield]: https://img.shields.io/github/issues/capcom6/service-monitor-tgbot.svg?style=for-the-badge
[issues-url]: https://github.com/capcom6/service-monitor-tgbot/issues
[license-shield]: https://img.shields.io/github/license/capcom6/service-monitor-tgbot.svg?style=for-the-badge
[license-url]: https://github.com/capcom6/service-monitor-tgbot/blob/master/LICENSE
[http-alert]: assets/http-alert.png
[tcp-alert]: assets/tcp-alert.png
[Golang]: https://img.shields.io/badge/Golang-000000?style=for-the-badge&logo=go&logoColor=white
[Golang-url]: https://go.dev/
[Fiber]: https://img.shields.io/badge/Fiber-000000?style=for-the-badge&logo=go&logoColor=white
[Fiber-url]: https://gofiber.io/
[Redis]: https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white
[Redis-url]: https://github.com/redis/go-redis
[Telegram]: https://img.shields.io/badge/Telegram-26A5E4?style=for-the-badge&logo=telegram&logoColor=white
[Telegram-url]: https://core.telegram.org/bots/api
