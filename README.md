# Go Hiring Challenge

This repository contains a Go application for managing products and their prices, including functionalities for CRUD operations and seeding the database with initial data.

## Project Structure

1. **cmd/**: Contains the main application and seed command entry points.

   - `server/main.go`: The main application entry point, serves the REST API.
   - `seed/main.go`: Command to seed the database with initial product data.

2. **app/**: Contains the application logic.
3. **sql/**: Contains a very simple database migration scripts setup.
4. **models/**: Contains the data models and repositories used in the application.
5. `.env`: Environment variables file for configuration.

## Setup Code Repository

1. Create a github/bitbucket/gitlab repository and push all this code as-is.
2. Create a new branch, and provide a pull-request against the main branch with your changes. Instructions to follow.

## Application Setup

- Ensure you have Go installed on your machine.
- Ensure you have Docker installed on your machine.
- Important makefile targets:
  - `make tidy`: will install all dependencies.
  - `make docker-up`: will start the required infrastructure services via docker containers.
  - `make seed`: ⚠️ Will destroy and re-create the database tables.
  - `make test`: Will run the tests.
  - `make run`: Will start the application.
  - `make docker-down`: Will stop the docker containers.

## Mockery (Mock Generation)

This project uses [mockery](https://github.com/vektra/mockery) to generate mock implementations for repository interfaces used in unit tests.

### Installation

```bash
go install github.com/vektra/mockery/v2@latest
```

> **Note:** Make sure the mockery binary is compiled with the same (or newer) Go version that the project requires. If you see errors like `package requires newer Go version`, clone the [mockery repo](https://github.com/vektra/mockery) and build from source:
>
> ```bash
> git clone --depth 1 https://github.com/vektra/mockery.git
> cd mockery
> go build -o $(go env GOPATH)/bin/mockery ./
> ```

### Usage

The configuration lives in [`.mockery.yaml`](.mockery.yaml). To regenerate all mocks, run from the project root:

```bash
mockery
```

This will generate/update mock files under each interface's `mocks/` subdirectory:

- `app/catalog/mocks/mock_product_repository.go`
- `app/categories/mocks/mock_category_repository.go`

### Adding a new mock

1. Add the interface to `.mockery.yaml` under `packages`:
   ```yaml
   packages:
     github.com/mytheresa/go-hiring-challenge/app/your_package:
       interfaces:
         YourInterface: {}
   ```
2. Run `mockery` to generate the mock.
3. Use it in tests:
   ```go
   repo := mocks.NewMockYourInterface(t)
   repo.EXPECT().YourMethod(args).Return(result, nil)
   ```

Follow up for the assignemnt here: [ASSIGNMENT.md](ASSIGNMENT.md)
