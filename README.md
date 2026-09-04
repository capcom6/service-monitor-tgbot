<a name="readme-top"></a>

<!-- PROJECT SHIELDS -->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]

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
  - [Storage Backends](#storage-backends)
- [Configuration](#configuration)
  - [Environment Variables](#environment-variables)
  - [YAML Config Structure](#yaml-config-structure)
  - [Service Definition Fields](#service-definition-fields)
- [Examples](#examples)
  - [HTTP Service Monitoring Example](#http-service-monitoring-example)
  - [TCP Service Monitoring Example](#tcp-service-monitoring-example)
- [Development](#development)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)
- [Acknowledgments](#acknowledgments)

<!-- ABOUT THE PROJECT -->
## About The Project

Monitoring the availability of network services is an important task for any project. At the same time, it is not always necessary to deploy universal solutions like Prometheus -- a fairly simple solution suffices for many cases. It is for such cases that this bot was created.

The bot monitors the availability of HTTP(S) and TCP services and sends Telegram notifications when their status changes. It also supports on-demand status queries via commands and periodic heartbeat messages.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Built With

* [![Golang][Golang]][Golang-url]
* [![Uber FX][fx-shield]][fx-url]
* [![Fiber][fiber-shield]][fiber-url]
* [![Telegram Bot API][telegram-shield]][telegram-url]
* [![Redis][redis-shield]][redis-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
## Getting Started

Follow the instructions below to run the bot.

### Prerequisites

- Docker, or
- Go 1.25+ and Make

### Installation

**Docker (single container):**

1. Register a new bot and get a token: https://core.telegram.org/bots/features#creating-a-new-bot
2. Create a [channel](https://telegram.org/tour/channels) or [group](https://telegram.org/tour/groups) for notifications
3. Add the bot to the channel/group as an administrator with message-sending permissions
4. Copy the configuration file [config.example.yml](configs/config.example.yml) to your working directory as `config.yml`
5. Edit `config.yml`:
    - Set the bot token
    - Set the channel/group ID (find it via `https://api.telegram.org/bot<token>/getUpdates?allowed_updates=[]` after adding the bot)
    - List the services to monitor (or configure a [storage backend](#storage-backends))
6. Run:
    ```bash
    docker run -d -v "$(pwd)/config.yml:/app/config.yml:ro" --name tgbot ghcr.io/capcom6/service-monitor-tgbot:latest
    ```

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

The bot uses a flexible template system for customizing notification messages. Templates use Go `text/template` syntax with a custom `escape` function for Telegram MarkdownV2 escaping.

Available templates:

| Template        | Variables                                                                                                                                                                                                                                 | Description                           |
| --------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------- |
| `online`        | `Name` - name of the service<br>`ChangedAt` - when the service entered this state<br>`Duration` - how long the service has been in this state                                                                                             | message when a service goes "online"  |
| `offline`       | `Name` - name of the service<br>`Error` - error message<br>`ChangedAt` - when the service entered this state<br>`Duration` - how long the service has been in this state                                                                  | message when a service goes "offline" |
| `services_list` | `.` - list of all services:<br>`Name` - name of the service<br>`State` - state of the service<br>`Error` - error message<br>`ChangedAt` - when the service entered this state<br>`Duration` - how long the service has been in this state | message with a list of all services   |
| `heartbeat`     | `TotalServices` - total number of services<br>`OnlineServices` - number of online services<br>`OfflineServices` - number of offline services<br>`CheckedAt` - time of the heartbeat check                                                 | periodic health summary               |

Templates can be overridden via the `TELEGRAM__MESSAGES` environment variable (JSON) or the `telegram.messages` YAML key.

### Commands

- `/status` - Get the current status of all monitored services

### Heartbeat

The heartbeat feature sends periodic "all clear" summaries to help distinguish between "everything is fine" and "the bot has crashed."

When enabled, the bot sends a message at a configurable interval showing how many services are online/offline.

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
  chatId: -1001234567890  # optional
```

Default heartbeat message: `💓 Heartbeat: 5/5 services online`

### Storage Backends

The bot supports pluggable storage backends for loading the list of monitored services. The backend is selected via DSN in the `storage.dsn` config field (or `STORAGE__DSN` environment variable).

| Scheme  | Backend             | Use Case                                      |
| ------- | ------------------- | --------------------------------------------- |
| `file`  | YAML file (default) | Simple setups, static configs                 |
| `redis` | Redis key           | Dynamic service lists, distributed management |

**File DSN format:** `file:///path/to/services.yml`

**Redis DSN format:** `redis://[[user]:pass@]host:port[/db][?key=<key>&channel=<channel>]`

Default key: `service-monitor:services`, default channel: `service-monitor:reload`.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONFIGURATION -->
## Configuration

Configuration is loaded in order: defaults < YAML file < environment variables.

Set `CONFIG_PATH` to point to your YAML config file. Environment variables override YAML values using double underscores for nesting (e.g., `TELEGRAM__TOKEN`).

### Environment Variables

| Variable                     | Type     | Default                   | Required | Description                                                          |
| ---------------------------- | -------- | ------------------------- | -------- | -------------------------------------------------------------------- |
| `CONFIG_PATH`                | string   |                           | **Yes**  | Path to main YAML config file                                        |
| `DEBUG`                      | bool     | `false`                   | No       | Enable dev logging (human-readable, debug level)                     |
| `HTTP__ADDRESS`              | string   | `127.0.0.1:3000`          | No       | Health/metrics HTTP server bind address                              |
| `HTTP__PROXY_HEADER`         | string   | `X-Forwarded-For`         | No       | Trusted proxy header name                                            |
| `HTTP__PROXIES`              | string   |                           | No       | Comma-separated trusted proxy IPs/CIDRs                              |
| `HTTP__OPENAPI__ENABLED`     | bool     | `true`                    | No       | Enable Swagger/OpenAPI docs                                          |
| `HTTP__OPENAPI__PUBLIC_HOST` | string   |                           | No       | Public hostname for OpenAPI spec                                     |
| `HTTP__OPENAPI__PUBLIC_PATH` | string   |                           | No       | Path prefix for OpenAPI endpoint                                     |
| `TELEGRAM__TOKEN`            | string   |                           | **Yes**  | Telegram bot API token                                               |
| `TELEGRAM__CHATID`           | int64    | `0`                       | **Yes**  | Target chat/group/channel ID for notifications                       |
| `TELEGRAM__PROXY_URL`        | string   |                           | No       | SOCKS5 proxy URL for Telegram API                                    |
| `TELEGRAM__TIMEOUT`          | duration | `1m`                      | No       | Telegram API client timeout                                          |
| `TELEGRAM__MESSAGES`         | JSON map | (built-in)                | No       | Custom message templates (online, offline, services_list, heartbeat) |
| `HEARTBEAT__ENABLED`         | bool     | `false`                   | No       | Enable periodic heartbeat messages                                   |
| `HEARTBEAT__INTERVAL`        | duration | `6h`                      | No       | Heartbeat send interval                                              |
| `HEARTBEAT__CHATID`          | int64    | (uses `TELEGRAM__CHATID`) | No       | Target chat for heartbeat messages                                   |
| `STORAGE__DSN`               | string   | `file://$CONFIG_PATH`     | No       | Storage backend DSN (`file://` or `redis://`)                        |

### YAML Config Structure

```yaml
telegram:
  token: <bot token>
  chatId: -1234567890123
  proxyUrl: socks5://<user>:<pass>@127.0.0.1:1080  # optional
  messages:
    online: "✅ {{.Name | escape}} is *online*"
    offline: "❌ {{.Name | escape}} is *offline*: {{.Error | escape}}"
    services_list: |
      {{ range . }}...{{ end }}
    heartbeat: "💓 Heartbeat: {{.OnlineServices}}/{{.TotalServices}} services online"

heartbeat:
  enabled: false
  interval: 6h
  chatId: -1234567890123  # optional, defaults to telegram.chatId

storage:
  dsn: file://./configs/services.yml  # or redis://...

services:   # inline services (used when storage is file:// with this same file)
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

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- EXAMPLES -->
## Examples

### HTTP Service Monitoring Example

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

### TCP Service Monitoring Example

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

<!-- DEVELOPMENT -->
## Development

### Prerequisites

- Go 1.25+
- [golangci-lint](https://golangci-lint.run/) v2
- [Air](https://github.com/air-verse/air) (for hot-reload, optional)
- [GoReleaser](https://goreleaser.com/) (for release builds, optional)

### Make Targets

| Target              | Description                                                 |
| ------------------- | ----------------------------------------------------------- |
| `make all`          | Format + lint + test coverage                               |
| `make build`        | Build binary to `bin/`                                      |
| `make test`         | Run tests with race detection and coverage                  |
| `make coverage`     | Generate coverage report (`coverage.out` + `coverage.html`) |
| `make lint`         | Run golangci-lint                                           |
| `make fmt`          | Format code                                                 |
| `make air`          | Development server with hot-reload                          |
| `make release`      | GoReleaser snapshot build                                   |
| `make docker-build` | Build Docker image                                          |
| `make docker-up`    | Start via Docker Compose                                    |
| `make docker-down`  | Stop Docker Compose                                         |
| `make swagger`      | Regenerate Swagger docs                                     |
| `make clean`        | Clean build artifacts                                       |
| `make help`         | Show available targets                                      |

### Quick Start (Development)

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

Project Link: [https://github.com/capcom6/service-monitor-tgbot](https://github.com/capcom6/service-monitor-tgbot)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

* [Best README Template](https://github.com/othneildrew/Best-README-Template)
* [Uber FX](https://github.com/uber-go/fx)
* [Fiber](https://gofiber.io/)
* [Telegram Bot API](https://core.telegram.org/bots/api)
* [go-redis](https://github.com/redis/go-redis)
* [koanf](https://github.com/knadh/koanf)

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
[fx-shield]: https://img.shields.io/badge/Uber%20FX-000000?style=for-the-badge&logo=uber&logoColor=white
[fx-url]: https://github.com/uber-go/fx
[fiber-shield]: https://img.shields.io/badge/Fiber-000000?style=for-the-badge
[fiber-url]: https://gofiber.io/
[telegram-shield]: https://img.shields.io/badge/Telegram%20Bot%20API-0088cc?style=for-the-badge&logo=telegram&logoColor=white
[telegram-url]: https://core.telegram.org/bots/api
[redis-shield]: https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white
[redis-url]: https://github.com/redis/go-redis
