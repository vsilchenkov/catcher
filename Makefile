.PHONY: build-win

build-win:
	go test -v ./...
# 	catcher
	cd app/cmd/catcher && goversioninfo versioninfo.json
	go build -o catcher.exe ./app/cmd/catcher
	rm -f app/cmd/catcher/resource.syso
# 	sentry_exporter
	cd app/cmd/sentry_exporter && goversioninfo versioninfo.json
	go build -o sentry_exporter.exe ./app/cmd/sentry_exporter
	rm -f app/cmd/sentry_exporter/resource.syso
	