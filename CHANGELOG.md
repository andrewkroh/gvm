# Change Log
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

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

[Unreleased]: https://github.com/andrewkroh/gvm/compare/v0.3.0...HEAD
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