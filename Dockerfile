FROM scratch

COPY ops /

ENTRYPOINT ["/ops"]
