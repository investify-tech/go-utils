# go-utils

Collection for utils to provide some recurrent convenient functionalities.

## Versioning

The project currently contains all sub-modules. As soon as this triggers issues in versioning (ex.: you want to have
different versions for `log` and `binexec`), it can be split into single projects. But for now it's OK that way. 

## Tests

Enter the subdir and execute: `go test`.

Execute all tests at once: `go test ./...`.

Some interesting and useful options:
 - `-v`: Detailed output
 - `-cover`: Include coverage
 - `-timeout=30s`: Set a timeout of 30s
 - `-count=1`: No caching for tests
