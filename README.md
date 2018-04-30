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