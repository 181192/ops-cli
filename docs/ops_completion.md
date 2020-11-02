## ops completion

Generate completion script

### Synopsis

To load completions:

Bash:

$ source <(ops completion bash)

# To load completions for each session, execute once:
Linux:
  $ ops completion bash > /etc/bash_completion.d/ops
MacOS:
  $ ops completion bash > /usr/local/etc/bash_completion.d/ops

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ ops completion zsh > "${fpath[1]}/_ops"

# You will need to start a new shell for this setup to take effect.

Fish:

$ ops completion fish | source

# To load completions for each session, execute once:
$ ops completion fish > ~/.config/fish/completions/ops.fish

PowerShell:

# Create a powershell profile (run as admin terminal)
if (!(Test-Path -Path $PROFILE)) {
	New-Item -ItemType File -Path $PROFILE -Force
}

# Open the profile.ps1
notepad $PROFILE

# And add the following
Invoke-Expression -Command $(ops completion powershell | Out-String)



```
ops completion [bash|zsh|fish|powershell]
```

### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
      --log-level string   Log level (debug, info, warn, error, fatal, panic) (default "info")
```

### SEE ALSO

* [ops](ops.md)	 - ops-cli is a wrapper for devops tools

