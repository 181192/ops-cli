FROM scratch

COPY dist/ops_cli_linux_amd64/ops_cli_linux_amd64 /ops-cli

ENTRYPOINT ["/ops-cli"]
