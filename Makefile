BINARY=msuser
OUTPUT=bin

test:
	go test -covermode=count -coverprofile=coverage.out
	go tool cover -func=coverage.out -o=coverage.out

local:
	echo "Building for local binary"
	go build -o ${OUTPUT}/${BINARY} main.go

build: local
	echo "Building for every OS and Platform"
	GOOS=freebsd GOARCH=amd64 go build -o ${OUTPUT}/${BINARY}-freebsd-amd64 main.go
	GOOS=linux GOARCH=amd64 go build -o ${OUTPUT}/${BINARY}-linux-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -o ${OUTPUT}/${BINARY}-windows-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -o ${OUTPUT}/${BINARY}-darwin-amd64 main.go

run: build
	echo "Run the default binary..."
	./${OUTPUT}/${BINARY}

clean:
	go clean
	rm -rf ${OUTPUT}

all: clean build
