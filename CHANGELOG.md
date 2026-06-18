# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.5.0] - 2026-06-12

### Added
- Timestamp tracking in status command and monitor

## [1.4.0] - 2026-06-07

### Changed
- Refactored project to use a standard project template
- Fixed config naming

## [1.3.0] - 2026-06-01

### Added
- Proxy support for Telegram client

### Removed
- Terraform deployment actions

## [1.2.2] - 2025-10-11

### Changed
- New errors model with improved error handling

## [1.2.1] - 2025-10-11

### Changed
- Improved task timing in monitor
- Improved graceful shutdown
- Use real service IDs in logs
- Minor fixes to error naming

### Maintenance
- Updated Go workflow
- Fixed lint errors

## [1.2.0] - 2025-10-08

### Added
- Active mode with `/status` command
- Services status list for monitor
- `escape` function for message templates
- Command type support in Telegram client
- Telegram updates query for development

### Changed
- Refactored messages to use additional `templates` package
- Refactored bot internal architecture
- Telegram client provides chat ID instead of user
- Upgraded `go-infra-fx` dependency
- Improved Telegram graceful shutdown
- Track last error and probe timestamp in monitor
- Removed unused storage code
- Renamed bot constant

### Fixed
- Race condition in messages

### Documentation
- Updated changelog
- Added messages templates and status command to README
- Updated example config

## [1.1.1] - 2025-10-02

### Changed
- Updated Go version
- Migrated to `samber/lo` instead of internal implementation

## [1.1.0] - 2023-12-26

### Added
- Storage interface with read-only storages
- Docker Compose setup
- Uber FX application framework
- Docker Swarm deployment with Terraform
- Notification message templates

### Removed
- Old code after migration to FX framework

### Fixed
- Graceful stop handling

### Changed
- Updated README

## [1.0.1] - 2023-02-09

### Added
- More items in roadmap

### Fixed
- Fake notifications in pending state
- Service name escaping
- Use pure Markdown for table of contents

### Changed
- Updated changelog and README

## [1.0.0] - 2023-01-20

### Added
- HTTP GET probe
- TCP probe
- Messages to Telegram channel
- Docker build and publish
- Simple message formatting
- Support of configurable thresholds
- Wait group for channel write protection

### Changed
- Refactored probes into dedicated `probes` package
- Renamed project to `service-monitor-tgbot`
- Renamed "ping" to "probe"
- Updated README
- Updated config comments

### Fixed
- Docker build parameters

[Unreleased]: https://github.com/capcom6/tgbot-service-monitor/compare/v1.5.0...HEAD
[1.5.0]: https://github.com/capcom6/tgbot-service-monitor/compare/v1.4.0...v1.5.0
[1.4.0]: https://github.com/capcom6/tgbot-service-monitor/compare/v1.3.0...v1.4.0
[1.3.0]: https://github.com/capcom6/tgbot-service-monitor/compare/v1.2.2...v1.3.0
[1.2.2]: https://github.com/capcom6/tgbot-service-monitor/compare/v1.2.1...v1.2.2
[1.2.1]: https://github.com/capcom6/tgbot-service-monitor/compare/v1.2.0...v1.2.1
[1.2.0]: https://github.com/capcom6/tgbot-service-monitor/compare/v1.1.1...v1.2.0
[1.1.1]: https://github.com/capcom6/tgbot-service-monitor/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/capcom6/tgbot-service-monitor/compare/v1.0.1...v1.1.0
[1.0.1]: https://github.com/capcom6/tgbot-service-monitor/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/capcom6/tgbot-service-monitor/releases/tag/v1.0.0
