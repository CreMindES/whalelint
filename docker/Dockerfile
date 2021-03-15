FROM golang:1.15 as goBuilder

COPY . /go/src/github.com/CreMindES/whalelint/

WORKDIR /go/src/github.com/CreMindES/whalelint/
RUN CGO_ENABLED=0 go build -o /app/whalelint

FROM hairyhenderson/upx:3.96 as upxBuilder

WORKDIR /app
COPY --from=goBuilder /app/whalelint /app/whalelint
RUN upx whalelint

FROM scratch

ARG gitCommitHash

LABEL maintainer="Tamas Gabor Barna"
LABEL org.opencontainers.image.source="https://github.com/CreMindES/whalelint"
LABEL org.opencontainers.image.description="Dockerfile linter."
LABEL org.opencontainers.image.revision=${gitCommitHash}

COPY --from=upxBuilder /app/whalelint /usr/local/bin/

CMD ["--help"]
ENTRYPOINT ["/usr/local/bin/whalelint"]
