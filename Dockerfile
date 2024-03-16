ARG build_folder="/opt/build/source"

FROM debian:bookworm-slim as builder

ARG build_folder

RUN apt update && \
    apt upgrade && \
    apt install --assume-yes golang ca-certificates

COPY source "$build_folder"

RUN cd "$build_folder/app" && \
    go build -ldflags "-s -w" -o app app.go && \
    ls -alh


FROM debian:bookworm-slim

ARG build_folder

COPY --from=builder "$build_folder/app/app" /app
COPY --from=builder "$build_folder/app/templates" /templates

ENTRYPOINT ["/app"]


