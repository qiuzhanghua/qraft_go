default: pre build post

all: pre build build-linux build-darwin build-windows post

pre:
	autotag write

build:
	go build -o qraft

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o qraft_`autotag current`_linux_arm64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o qraft_`autotag current`_linux_amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o qraft_`autotag current`_linux_386

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o qraft_`autotag current`_darwin_amd64
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o qraft_`autotag current`_darwin_arm64

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o qraft_`autotag current`_windows_amd64.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o qraft_`autotag current`_windows_386.exe

post:
	git restore autotag.go

clean:
	rm qraft
	rm qraft*.*