server: server.go Makefile
	CGO_ENABLED=0 go build -a -x -installsuffix cgo server.go

