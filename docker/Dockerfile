FROM golang:1.17 as goBuilder

COPY . /go/src/github.com/CreMindES/whalelint/

WORKDIR /go/src/github.com/CreMindES/whalelint/
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /app/whalelint

FROM hairyhenderson/upx:3.96 as upxBuilder

WORKDIR /app
COPY --from=goBuilder /app/whalelint /app/whalelint
RUN upx -9 whalelint

FROM scratch

ARG gitCommitHash

LABEL maintainer="Tamas Gabor Barna"
LABEL source="https://github.com/CreMindES/whalelint"
LABEL description="Dockerfile linter."
LABEL revision=${gitCommitHash}

COPY --from=upxBuilder /app/whalelint /usr/local/bin/

CMD ["--help"]
ENTRYPOINT ["/usr/local/bin/whalelint"]
