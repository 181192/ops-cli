## ops completion bash

Generates bash completion scripts

### Synopsis

To load completion run

source <(ops completion bash)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
source <(ops completion bash)

If you are stuck on Bash 3 (macOS) use

source /dev/stdin <<<"$(ops completion bash)"



```
ops completion bash [flags]
```

### Options

```
  -h, --help   help for bash
```

### Options inherited from parent commands

```
      --log-level string   Log level (debug, info, warn, error, fatal, panic) (default "info")
```

### SEE ALSO

* [ops completion](ops_completion.md)	 - Generates shell completion scripts

