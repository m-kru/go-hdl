[![Tests](https://github.com/m-kru/go-thdl/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/m-kru/go-thdl/actions?query=main)

# THDL

Thdl is a tool for easing the work with hardware description languages.
It is (and will always be) based solely on text processing, with no semantic analysis.
Such an approach draws a clear line between what might be included and what will never be supported.
The 'THDL' acronym doesn't have any expansion.
The first 'T' letter might be interpreted as 'text', as the tool is based on text processing.
However, don't read THDL as "Text Hardware Description Language" and do not treat it as such.
Part of the first prototype was implemented when I was on a train.
The 'train' word also starts with 't', so I thought 'thdl' would be a good name.

## Documentation

The documentation is included in the `thdl` binary.
Simply run `thdl help` or `thdl help {command}` to read about particular command.

## Commands

The commands are:
* `doc` - show or generate documentation,
* `gen` - generate code by processing sources,
* `vet` - check for likely mistakes.

## Installation

### go
```
go install github.com/m-kru/go-thdl/cmd/thdl@latest
```

Go installation installs to go configured path.

### Manual

```
git clone https://github.com/m-kru/go-thdl.git
make
make install
```

Manual installation installs to `/usr/bin`.

## Examples

### doc

<p align="center"><img src="assets/doc_hctsp.png?raw=true"/></p>

<p align="center"><img src="assets/doc_reset_synchronizer.png?raw=true"/></p>

<p align="center"><img src="assets/doc_t_command.png?raw=true"/></p>
