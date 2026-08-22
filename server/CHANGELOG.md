# Changelog

## [1.3.0] - 2026-08-22

### Added

- Table-driven unit tests for core functions (alias/version parsing, repository URL splitting, config loading, dependency stats, filter counts, summary generation, notification mapping)
- Typed `FilterCount` stats and `NotificationData` DTO to decouple API responses from go-github types

### Changed

- Upgrade to Go 1.27
- Migrate GoFiber v2 to v3
- Upgrade go-github v29 to v90 (drops `golang.org/x/oauth2`, uses `github.WithAuthToken`)
- Upgrade zap to 1.28.0
- Drop vendored dependencies in favor of go.sum-only resolution

### Removed

- Dead code in `handlers.go` (duplicated routes.go)
- Baked-in `GITHUB_AUTH_TOKEN` from the Dockerfile; must be supplied at runtime

## [1.2.0] - 2021-05-27

### Changed

- Upgrade GoFiber to 2.10.0
- Update the README

## [1.1.0] - 2020-07-05

### Added

- Update README adding instructions how to deploy into the Cloud Run

### Changed

- Migrate server to [gofiber](https://github.com/gofiber/fiber)
- Upgrade dependencies
- Change port `3000` to `8080`

### Fixed

- Bug that included the name of the wrong projects

## [1.0.0] - 2020-06-28

### Added

- Add validation show when there is an unnamend package json
- Add filter to Angular 10

### Fixed

- Bug that included the name of the wrong projects