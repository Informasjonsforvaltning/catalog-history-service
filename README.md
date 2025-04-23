# Catalog History Service

This application provides an API to keep track of changes (JSON Patch operations) to datasets, dataservices, concepts
and services.

For a broader understanding of the system’s context, refer to
the [architecture documentation](https://github.com/Informasjonsforvaltning/architecture-documentation) wiki. For more
specific context on this application, see the **Registration** subsystem section.

## Getting Started

These instructions will give you a copy of the project up and running on your local machine for development and testing
purposes.

### Prerequisites

Ensure you have the following installed:

- Go (version 1.17 or higher)
- Docker

### Running locally

Clone the repository.

```sh
git clone https://github.com/Informasjonsforvaltning/catalog-history-service.git
cd catalog-history-service
```

Run `go get` to install the required dependencies.

```shell
go get
```

Start MongoDB and the application (either through your IDE, or via CLI):

```sh
docker compose up -d
go run main.go
```

### API Documentation (OpenAPI)

The API documentation is available at ```openapi.yaml```.

### Running tests

```shell
go test ./test
```

To generate a test coverage report, use the following command:

```shell
go test -v -race -coverpkg=./... -coverprofile=coverage.txt -covermode=atomic ./test
```
