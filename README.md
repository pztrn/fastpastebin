# Fast Pastebin

[![Build Status](https://ci.code.pztrn.name/api/badges/apps/fastpastebin/status.svg)](https://ci.code.pztrn.name/apps/fastpastebin)

Easy-to-use-and-install pastebin software written in Go. No bells or whistles, no websockets and even NO JAVASCRIPT!

**Please, use [my gitea](https://code.pztrn.name/apps/fastpastebin) for bug reporting. All other places are mirrors!**

## Current functionality

* Create and view public and private pastes.
* Syntax highlighting.
* Pastes expiration.
* Passwords for pastes.
* Multiple storage backends. Currently: ``flatfiles``, ``mysql`` and ``postgresql``.

## Caveats

* No links at lines numbers. See [this Chroma bug](https://github.com/alecthomas/chroma/issues/132)

## Installation and updating

Just issue:

```bash
CGO_ENABLED=0 go install go.dev.pztrn.name/fastpastebin/cmd/fastpastebin@latest
```

This command can be used to update Fast Paste Bin.

## Configuration

Take a look at [example configuration file](examples/fastpastebin.yaml.dist) which contains all supported options and their descriptions.

Configuration file position is irrelevant, there is no hardcoded paths where Fast Paste Bin looking for it's configuration. Use ``-config`` CLI parameter or ``FASTPASTEBIN_CONFIG`` environment variable to specify path.

## Developing

Use linters, formatters, etc. VSCode with Go plugin is recommended for developing as it will perform most of linting-formatting
actions automagically.

Also, Sublime Text with LSP-gopls will also work just fine.

Try to follow [Go's code review comments](https://github.com/golang/go/wiki/CodeReviewComments) with few exceptions:

* We're not forcing any limits on line length for code, only for comments, they should be 72-76 chars long.

## ToDo

This is a ToDo list which isn't sorted by any parameter at all. Just a list of tasks you can help with.

* User CP.
* Files uploading.
* Passwords for files.
* Pastes forking and revisioning (like git or github gists).
* Possibility to copy-paste-edit WYSIWYG content.
* CLI client for pastes and files uploading.
