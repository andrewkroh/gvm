# Change Log
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

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
