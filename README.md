## tips

go test -gcflags=-l -v -run TestExtConv .
go test -gcflags=-l -coverprofile=coverage.out .
go tool cover -html=coverage.out -o coverage.htm