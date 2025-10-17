# Change Log
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

### Changed

### Fixed

### Added

## [0.6.0]

### Changed

- Migrated from Google Cloud Storage to the official Go downloads API at go.dev/dl to resolve HTTP 403 errors when downloading Go binaries. [#117](https://github.com/andrewkroh/gvm/pull/117)

## [0.5.3]

### Changed

- This release is built with Go 1.24.6. [#109](https://github.com/andrewkroh/gvm/pull/109)

### Fixed

- Fixed path separator format for Powershell on Linux [#81](https://github.com/andrewkroh/gvm/issues/81)

### Added

## [0.5.2]

### Changed

- Add 10s wait between retries when there are download errors. [#65](https://github.com/andrewkroh/gvm/pull/65)
 
## [0.5.1]

### Changed

- This release is built with Go 1.20.7. [#63](https://github.com/andrewkroh/gvm/pull/63)

### Fixed

- Fixed installation of Go 1.21.0. Versioning conventions changed as per golang/go#57631. [#61](https://github.com/andrewkroh/gvm/issues/61)
- Fixed issues unpacking tar files for Go 1.21.0 releases. Relates to golang/go#61862. [#61](https://github.com/andrewkroh/gvm/issues/61)

## [0.5.0]

### Changed

- Updated releases to use Go 1.18. [#54](https://github.com/andrewkroh/gvm/pull/54)
- Report Go module version from `gvm --version` if installed via `go install`. [#57](https://github.com/andrewkroh/gvm/pull/57)

### Fixed

- Fix `--arch` flag and associated `GVM_ARCH` env variable. [#53](https://github.com/andrewkroh/gvm/pull/53)

## [0.4.1]

### Fixed

- Fixed builds under Go 1.18. [#46](https://github.com/andrewkroh/gvm/issues/46) [#47](https://github.com/andrewkroh/gvm/pull/47)
- Fixed several issues identified by golangci-lint. [#48](https://github.com/andrewkroh/gvm/pull/48)

### Added

- Added macos universal build and switched to goreleaser for release builds. [#50](https://github.com/andrewkroh/gvm/pull/50)

## [0.4.0]

### Added

- Add `--http-timeout` flag to control the timeout for HTTP requests. [#43](https://github.com/andrewkroh/gvm/issues/43) [#45](https://github.com/andrewkroh/gvm/pull/45)

## [0.3.2]

### Added

- Add an artifact for Apple M1 (darwin/arm64).

## [0.3.1]

### Changed

- Use Go 1.17 to build project. [#40](https://github.com/andrewkroh/gvm/pull/40)

### Fixed

- Fix staleness issues with the `gvm available` output. [#39](https://github.com/andrewkroh/gvm/issues/39) [#41](https://github.com/andrewkroh/gvm/pull/41)

## [0.3.0]

### Added

- Added new `gvm use` flag `--no-install` (or `-n`) to disable installing
  and updating from source. For example, `gvm use tip -n` will use the tip
  assuming you have it already installed, but will not trigger an update.
  [#35](https://github.com/andrewkroh/gvm/pull/35)

## [0.2.4]

### Fixed

- Fix errors with renames failing across disks by falling back to a copy/delete. [#31](https://github.com/andrewkroh/gvm/pull/31)
- Fall back to a source code based install only when binary package URL returns
  with HTTP 404. [#30](https://github.com/andrewkroh/gvm/issues/30) [#32](https://github.com/andrewkroh/gvm/pull/32)
  
## [0.2.3]

### Changed 

- Use Go 1.15 to build project. [#28](https://github.com/andrewkroh/gvm/pull/28)
- Log a warning when binary package install fails. [#29](https://github.com/andrewkroh/gvm/pull/29)

## [0.2.2]

### Fixed

- Fix output of `--format match` to remove a stray `$` character. #25

## [0.2.1]

### Changed

- Update repo to use `go mod` for dependency management. #17

## [0.2.0]

## Added

- Add logic to retry failed downloads. #15

## [0.1.0]

### Added

- Add the ability to build Go from source (e.g. gvm tip).
- Add `gvm init` command that clones the Go git repository or forces
  a git fetch.
- Add `gvm install` command.
- Add `gvm list` command.
- Add `gvm available` command.
- Add `gvm purge` command.
- Add `gvm remove` command.

## [0.0.5]

### Fixed

- Download armv6l releases when GOARCH is ARM. be730193cac29bd64f751a7104a32883703741b1

## [0.0.4]

### Fixed

- Fixed documentation for batch usage (cmd.exe) on README. #5

### Changed

- Binary releases are built with Go 1.9.2. #6

### Added

- Added ARM releases for Linux and FreeBSD. #6

## [0.0.3]

### Added

- Added tests to check the environment variables returned by gvm.

## [0.0.2]

### Changed

- Changed code to extract the Go distribution directly into HOME instead
  of extracting to a temp dir and then moving to HOME. #1

## [0.0.1]

Initial release.

[Unreleased]: https://github.com/andrewkroh/gvm/compare/v0.6.0...HEAD
[0.6.0]: https://github.com/andrewkroh/gvm/releases/tag/v0.6.0
[0.5.3]: https://github.com/andrewkroh/gvm/releases/tag/v0.5.3
[0.5.2]: https://github.com/andrewkroh/gvm/releases/tag/v0.5.2
[0.5.1]: https://github.com/andrewkroh/gvm/releases/tag/v0.5.1
[0.5.0]: https://github.com/andrewkroh/gvm/releases/tag/v0.5.0
[0.4.1]: https://github.com/andrewkroh/gvm/releases/tag/v0.4.1
[0.4.0]: https://github.com/andrewkroh/gvm/releases/tag/v0.4.0
[0.3.2]: https://github.com/andrewkroh/gvm/releases/tag/v0.3.2
[0.3.1]: https://github.com/andrewkroh/gvm/releases/tag/v0.3.1
[0.3.0]: https://github.com/andrewkroh/gvm/releases/tag/v0.3.0
[0.2.4]: https://github.com/andrewkroh/gvm/releases/tag/v0.2.4
[0.2.3]: https://github.com/andrewkroh/gvm/releases/tag/v0.2.3
[0.2.2]: https://github.com/andrewkroh/gvm/releases/tag/v0.2.2
[0.2.1]: https://github.com/andrewkroh/gvm/releases/tag/v0.2.1
[0.2.0]: https://github.com/andrewkroh/gvm/releases/tag/v0.2.0
[0.1.0]: https://github.com/andrewkroh/gvm/releases/tag/v0.1.0
[0.0.5]: https://github.com/andrewkroh/gvm/releases/tag/v0.0.5
[0.0.4]: https://github.com/andrewkroh/gvm/releases/tag/v0.0.4
[0.0.3]: https://github.com/andrewkroh/gvm/releases/tag/v0.0.3
[0.0.2]: https://github.com/andrewkroh/gvm/releases/tag/v0.0.2
[0.0.1]: https://github.com/andrewkroh/gvm/releases/tag/v0.0.1
