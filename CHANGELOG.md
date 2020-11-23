# Change Log
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

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
