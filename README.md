Chat at #fastpastebin @ irc.oftc.net or in Matrix at #fastpastebin:pztrn.name

# Fast Pastebin

Easy-to-use-and-install pastebin software written in Go. No bells or
whistles, no websockets and even NO JAVASCRIPT!

# Current functionality.

* Create and view public and private pastes.
* Syntax highlighting.
* Pastes expiration.
* Passwords for pastes.
* Multiple storage backends. Currently: ``flatfiles`` and ``mysql``.

# Caveats.

* No links at lines numbers. See https://github.com/alecthomas/chroma/issues/132

# Installation and updating

Just issue:

```
go get -u -v github.com/pztrn/fastpastebin/cmd/fastpastebin
```

This command can be used to update Fast Paste Bin.

**WARNING:** installation by compiling Fast Paste Bin from sources **require**
at least 300 megabytes of free RAM! Eventually it'll run even on 64MB-powered
VM, it's only a compilation issue.

# Configuration.

Take a look at [example configuration file](examples/fastpastebin.yaml.dist)
which contains all supported options and their descriptions.

Configuration file position is irrelevant, there is no hardcoded paths where
Fast Paste Bin looking for it's configuration. Use ``-config`` CLI parameter
or ``FASTPASTEBIN_CONFIG`` environment variable to specify path.

# Developing

Developers should install https://github.com/UnnoTed/fileb0x/ which is used
as replacement to go-bindata for embedding assets into binary. After changing
assets they should be recompiled into Go code. At repository root execute
this command and you'll be fine:

```
fileb0x fileb0x.yml
```

Also if you're changed list of assets (by creating or deleting them) be sure
to fix files list in ``fileb0x.yml`` file!

The rest is default - use linters, formatters, etc. VSCode with Go plugin is 
recommended for developing as it will perform most of linting-formatting
actions automagically. Try to follow https://github.com/golang/go/wiki/CodeReviewComments
with few exceptions:

* Imports should be organized in 3 groups: stdlib, local, other. See
https://github.com/pztrn/fastpastebin/blob/master/pastes/api_http.go for
example.

* We're not forcing any limits on line length for code, only for comments,
they should be 72-76 chars long.

# ToDo

This is a ToDo list which isn't sorted by any parameter at all. Just a list
of tasks you can help with.

* User CP.
* Files uploading.
* Passwords for files.
* Pastes forking and revisioning (like git or github gists).
* Possibility to copy-paste-edit WISYWIG content.
* CLI client for pastes and files uploading.