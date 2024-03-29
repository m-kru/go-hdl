# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.6.0] 2022-11-25
### Added
- [VHDL, doc] Support alias in package.
- [VHDL, doc] Bold more standard types.
### Changed
- Rename project from `thdl` to `hdl`.
### Fixed
- [VHDL] Ignore directories having \.vhdl? extension in name.
- [VHDL, gen] Add missing `to_str()` call for record attributes.
- [VHDL] Fix reading .hdl.yml config file when no extra arguments are provided.

## [0.5.0] 2022-06-22
### Added
- [VHDL, doc] Support package variables.
- [VHDL, gen] Support unknown types in generated records.
- [VHDL, gen] Support generated enumerations in generated records.
- [vet] Add possibility to vet only single file.

## [0.4.0] 2022-05-16
### Added
- [VHDL, doc] Protected types are documented in html.
- [VHDL, doc] Protected types are included in the package summary.
- [VHDL, gen] Initial support for enumeration and record types.

## [0.3.0] 2022-05-16
### Added
- [VHDL, doc] Initial support for HTML generation.
### Changed
- [VHDL, doc] Change library summary sections to: Entities, Packages, Testbenches.

## [0.2.0] - 2022-04-13
### Added
- [VHDL, doc] Support protected type.
- [VHDL, doc] Actual names used in library summary, not lowercased.
- [VHDL, doc] Improved array, enumeration and record type declarations scanning.
### Changed
- [VHDL, doc] Package summary is printed instead of full package code.
### Fixed
- [VHDL, doc] Fixed keywords bolding in strings and comments.

## [0.1.0] - 2022-04-04
First users release.
