![CI](https://github.com/mwazovzky/assistant/actions/workflows/test.yml/badge.svg)

# mwazovzky/assistant

Package mwazovzky/assistant implements simple open ai api client.

## Install

```
go get github.com/mwazovzky/assistant
```

## Basic Usage Example

For basic usage examples, please refer to `example_test.go`.

## Test

```
go test
go test -v ./...
go test -v assistant_test.go
go test -v -run TestCreateThread
go test example_test
```

## Test Coverage

```
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
open coverage.html
```

## Tag version

```
git tag -a v0.1.1 -m "Version 0.1.1: Revised the assistant.Ask() method signature to eliminate the return of usage statistics."
git push origin v0.1.1
```
