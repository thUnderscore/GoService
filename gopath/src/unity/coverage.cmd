rem go test -cover unity/example/shared
go test -coverprofile=coverage.out unity/example/shared
go tool cover -html=coverage.out