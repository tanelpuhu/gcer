build:
	@mkdir -p build
	GOOS=linux GOARCH=amd64 go build -o build/gcer_linux_amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -o build/gcer_darwin_amd64 main.go
