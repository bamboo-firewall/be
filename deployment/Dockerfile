FROM golang:1.20 as builder

WORKDIR /go/src/github.com/bamboo-firewall/be
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/service /go/src/github.com/bamboo-firewall/be/cmd/main.go

FROM alpine:3.15
WORKDIR /app
ENV TZ=Asia/Ho_Chi_Minh
RUN mkdir config
COPY --from=builder /bin/service /app/service
COPY --from=builder /go/src/github.com/bamboo-firewall/be/config.json /app/config.json
COPY --from=builder /go/src/github.com/bamboo-firewall/be/config/rbac_model.conf /app/config/rbac_model.conf
# add package for handle timezone in alpine
RUN apk add tzdata \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone \
    && chmod u+x /app/service

CMD [ "/app/service" ]