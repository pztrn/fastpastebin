# Flagger

[![GoDoc](https://godoc.org/go.dev.pztrn.name/flagger?status.svg)](https://godoc.org/go.dev.pztrn.name/flagger) [![Drone (self-hosted)](https://img.shields.io/drone/build/libraries/flagger?server=https%3A%2F%2Fci.dev.pztrn.name)](https://ci.dev.pztrn.name/libraries/flagger/) [![Discord](https://img.shields.io/discord/632359730089689128)](https://discord.gg/qHN6KsD) ![Keybase XLM](https://img.shields.io/keybase/xlm/pztrn) [![Go Report Card](https://goreportcard.com/badge/go.dev.pztrn.name/flagger)](https://goreportcard.com/report/go.dev.pztrn.name/flagger)

Flagger is an arbitrary CLI flags parser, like argparse in Python.
Flagger is able to parse boolean, integer and string flags.

## Installation

```bash
go get -u -v go.dev.pztrn.name/flagger
```

## Usage

Flagger requires logging interface to be passed on initialization.
See ``loggerinterface.go`` for required logging functions.
It is able to run with standart log package, in that case
initialize flagger like:

```go
flgr = flagger.New("My Super Program", flagger.LoggerInterface(log.New(os.Stdout, "testing logger: ", log.Lshortfile)))
flgr.Initialize()
```

Adding a flag is easy, just fill ``Flag`` structure and pass to ``AddFlag()`` call:

```go
flag_bool := Flag{
    Name: "boolflag",
    Description: "Boolean flag",
    Type: "bool",
    DefaultValue: true,
}
err := flgr.AddFlag(&flag_bool)
if err != nil {
    ...
}
```

After adding all neccessary flags you should issue ``Parse()`` call to get
them parsed:

```go
flgr.Parse()
```

After parsed they can be obtained everywhere you want, like:

```go
val, err := flgr.GetBoolValue("boolflag")
if err != nil {
    ...
}
```

For more examples take a look at ``flagger_test.go`` file or [at GoDoc](https://godoc.org/go.dev.pztrn.name/flagger).
