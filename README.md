[![Tests](https://github.com/m-kru/go-hdl/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/m-kru/go-hdl/actions?query=main)

# hdl

Hdl is a tool for easing the work with hardware description languages.

## Documentation

The documentation is included in the `hdl` binary.
Simply run `hdl help` or `hdl help {command}` to read about particular command.

## Commands

The commands are:
* `doc`  - show or generate documentation,
* `gen`  - generate code by processing sources,
* `help` - print more information about a specific command,
* `ver`  - print hdl version,
* `vet`  - check for likely mistakes.

## Installation

### go
```
go install github.com/m-kru/go-hdl/cmd/hdl@latest
```

Go installation installs to go configured path.

### Manual

```
git clone https://github.com/m-kru/go-hdl.git
make
make install
```

Manual installation installs to `/usr/bin`.

## Examples

### doc

<p align="center"><img src="assets/doc_hctsp.png?raw=true"/></p>

<p align="center"><img src="assets/doc_reset_synchronizer.png?raw=true"/></p>

<p align="center"><img src="assets/doc_t_command.png?raw=true"/></p>
