.PHONY: build-win

build-win:
	go test -v ./...
	cd app/cmd/catcher && goversioninfo versioninfo.json
	go build -o catcher.exe ./app/cmd/catcher
	rm -f app/cmd/catcher/resource.syso