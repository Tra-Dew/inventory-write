# Tradew Inventory Write

<div align="center"><img src="docs/logo.png" alt="logo" width="300"/></div>

## Overview
The `inventory-write` microservice is a write only service, from the the `inventory` API. It's responsible for managing queries of all items

### Why separate the same API into 2 microservices?
We could expect the read side of the `inventory` api to be a lot bigger. So by separating then into 2 microservices we can scale them individually. For more informatin read [CQRS](https://en.wikipedia.org/wiki/Command%E2%80%93query_separation) and [Read/Write Ratio](https://support.liveoptics.com/hc/en-us/articles/229590547-Live-Optics-Basics-Read-Write-Ratio)

## Usage

### HTTP
To start the http server on port `9001` run the command:
```
go run main.go api
```

### GRPC
To start the grpc server on port `9005` run the command:
```
go run main.go grpc
```

### Workers
To start the worker `dispatch-item-updated-worker` run the command:
```
go run main.go dispatch-item-updated-worker
```

## Docker

You can also run using docker, go in the root of the workspace and run:
```
docker build . -t inventory-write
docker run inventory-write -p 9001:9001
```

## Migrations
To create a SQL migration, first download the [migrate](https://github.com/golang-migrate/migrate) tool and then run the first migration

```
migrate create -ext sql -dir ./migrations -seq create_database
```

Run migrations:
```
migrate -database <CONNECTION_STRING> -path ./migrations up
```

Roolback migrations:
```
migrate -database <CONNECTION_STRING> -path ./migrations down
```

## Protobuf
Install [protoc](https://grpc.io/docs/protoc-installation/) and then run the command to generate the pb files
```
protoc --go_out=. --go-grpc_out=. pkg/inventory/proto/service.proto
```


## Architecture overview

![image of the architecture](./docs/architecture.png)