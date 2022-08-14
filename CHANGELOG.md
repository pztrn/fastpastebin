# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.4.1] - 2022-08-14

### Changed

* Update docker images - alpine to 3.16.1, golang to 1.19, golangci-lint to 1.48.0.
* Update chroma to v2.2.0.
* Update bulma to 0.9.4.
* Update bulma-tooltip to 1.2 (was 3.0.0, but repo was switched to CreativeBulma).
* Update github.com/dchest/captcha to v1.0.0.

### Fixed

* Put valid repository's link in footer.

### Removed

* Removed `flagger` dependency.

## [0.4.0] - 2021-01-09

### Added

* PostgreSQL support.
* Docker containerization for every commit and tag.
* Pastes cleanup procedure.

### Changed

* Updated bulma to v 0.7.5.
* Moved from `io/ioutil` to `os` package for reading/writing files and directories.

### Fixed

* Dirty hack to get database connection reestablish (for sure).

## 0.3.0

Release changelogs lost :(.

## [0.2.0] - 2018-05-27

### Added

* Possibility to create different database backends. Currently `mysql` and `flatfiles` are available.

### Changed

* De-hardcoded pagination configuration, it is now configurable via configuration file.

## [0.1.1] - 2018-05-26

### Added

* Footer copyrights.

### Changed

* Refactored templates: now they're included in each other if neccessary.

### Fixed

* Fixed nasty bugs with private pastes that causing fastpastebin to crash.
* Logger level from configuration now properly set.

## [0.1.0] - 2018-05-19

First normal release. Fast Paste Bin is able to handle public, private
and passworded pastes.

[Unreleased]: https://code.pztrn.name/apps/fastpastebin/compare/v0.4.1...HEAD
[0.4.1]: https://code.pztrn.name/apps/fastpastebin/compare/0.4.0...v0.4.1
[0.4.0]: https://code.pztrn.name/apps/fastpastebin/compare/v0.2.0...0.4.0
[0.2.0]: https://code.pztrn.name/apps/fastpastebin/compare/v0.1.1...v0.2.0
[0.1.1]: https://code.pztrn.name/apps/fastpastebin/compare/v0.1.0...v0.1.1
[0.1.0]: https://code.pztrn.name/apps/fastpastebin/src/tag/v0.1.0
