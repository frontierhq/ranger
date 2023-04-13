build:
	go build cmd/ranger/ranger.go

install:
	go install cmd/ranger/ranger.go

test:
	go test -v ./...
