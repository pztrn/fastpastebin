[![GoDoc](https://godoc.org/github.com/pztrn/flagger?status.svg)](https://godoc.org/gitlab.com/pztrn/flagger)

# Flagger

Flagger is an arbitrary CLI flags parser, like argparse in Python.
Flagger is able to parse boolean, integer and string flags.

# Installation

```
go get -u -v lab.pztrn.name/golibs/flagger
```

# Usage

Flagger requires logging interface to be passed on initialization.
See ``loggerinterface.go`` for required logging functions.
It is able to run with standart log package, in that case
initialize flagger like:

```
flgr = flagger.New("My Super Program", flagger.LoggerInterface(log.New(os.Stdout, "testing logger: ", log.Lshortfile)))
flgr.Initialize()
```

Adding a flag is easy, just fill ``Flag`` structure and pass to ``AddFlag()`` call:

```
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

```
flgr.Parse()
```

After parsed they can be obtained everywhere you want, like:

```
val, err := flgr.GetBoolValue("boolflag")
if err != nil {
    ...
}
```

For more examples take a look at ``flagger_test.go`` file or [at GoDoc](https://godoc.org/gitlab.com/pztrn/flagger).

# Get help

If you want to report a bug - feel free to report it via Gitlab's issues system. Note that everything that isn't a bug report or feature request will be closed without any futher comments.

If you want to request some help (without warranties), propose a feature request or discuss flagger in any way - please use our mailing lists at flagger@googlegroups.com. To be able to send messages there you should subscribe by sending email to ``flagger+subscribe@googlegroups.com``, subject and mail body can be random.