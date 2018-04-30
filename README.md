[![Join the chat at https://gitter.im/fastpastebin/Lobby](https://badges.gitter.im/fastpastebin/Lobby.svg)](https://gitter.im/fastpastebin/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

# Fast Pastebin

Easy-to-use-and-install pastebin software written in Go.

# Installation and updating

Just issue:

```
go get -u -v github.com/pztrn/fastpastebin
```

This command can be used to update Fast Paste Bin.

# Configuration.

Take a look at [example configuration file](examples/fastpastebin.yaml.dist)
which contains all supported options and their descriptions.

Configuration file position is irrelevant, there is no hardcoded paths where
Fast Paste Bin looking for it's configuration. Use ``-config`` CLI parameter
or ``FASTPASTEBIN_CONFIG`` environment variable to specify path.

# ToDo

This is a ToDo list which isn't sorted by any parameter at all. Just a list
of tasks you can help with.

* Pastes expiration. It saves time to database but isn't blocking access.
* User CP.
* Files uploading.
* Passwords for pastes and files.
* Code highlight.
* Pastes forking and revisioning (like git or github gists).
* Possibility to copy-paste-edit WISYWIG content.
* CLI client for pastes and files uploading.