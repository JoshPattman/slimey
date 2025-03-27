all: clean-bin mac-apple-arm linux-x86

mac-apple-arm: clean-bin
	GOOS=darwin GOARCH=amd64 go build -o bin/slimey-mac-arm

linux-x86: clean-bin
	GOOS=linux GOARCH=amd64 go build -o bin/slimey-linux

clean-bin:
	rm -rfd bin/
	mkdir bin/