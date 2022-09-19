PREFIX=/usr/local
VERSION=v$(shell cat assets/version.json | jq .version)

build: format
	go build -trimpath -buildmode=pie -mod=readonly -modcacherw -ldflags="-s -w"

format:
	gofmt -s -w .

test:
#	go test $(wildcard ./pkg/*) -v
	go test ./pkg/boardgames -v

version:
	git tag -f $(VERSION)

install: build
	mkdir -p $(PREFIX)/bin
	cp -f servus $(PREFIX)/bin

install-service:
	cp -f systemd/servus.service /etc/systemd/system/
