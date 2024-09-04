FROM golang:1.23

WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download


RUN --mount=type=bind,target=. \
   go build --ldflags "-s -w" -o /bin/server ./cmd/server

CMD ["/bin/server"]

