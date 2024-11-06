linux:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -buildid=" -trimpath

windows:
	CGO_ENABLED=0 GOOS=windows go build -ldflags="-s -w -buildid=" -trimpath

all: linux windows
