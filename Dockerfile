FROM scratch

COPY ./build/ops_cli_linux_amd64 /ops-cli

ENTRYPOINT ["/ops-cli"]
CMD ["version"]
