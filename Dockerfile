FROM golang:1.23-bookworm AS base

WORKDIR /golangProject


COPY go.mod go.sum ./


RUN go mod download


RUN go install github.com/swaggo/swag/cmd/swag@latest
ENV PATH="/go/bin:${PATH}"

COPY . .

RUN ls -al /golangProject

RUN mkdir -p docs/swagger

RUN swag init \
    --parseDependency \
    --parseInternal \
    --parseDepth 5 \
    -g internal/app/start/start.go \
    --output internal/app/swagger

WORKDIR /golangProject/cmd/app

RUN ls -al .

RUN go build -o /main

EXPOSE 8080

CMD ["/main"]
