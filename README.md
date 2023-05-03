# github.com/NewStreetTechnologies/go-backend-boilerplate

## Run in local
```go
go run main.go
```
## Update repository mock command

```go
cd ./repository
mockery --all --output ../test/mocks
cd ../integration
mockery --all --output ../test/mocks
```

## Unit Test command
```go
// Run Test for controllers unit test
cd ./test/controllers
go test -v

// Get Test Coverage
go test ./test/... -coverpkg ./...

// Get Test Coverage of controllers
go test ./test/... -coverpkg ./controllers/