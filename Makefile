PREFIX=/usr/local
VERSION=v$(shell cat assets/version.json | jq .version)

format:
	gofmt -s -w .

version:
	git tag $(VERSION)

build: format version
	go build -trimpath -buildmode=pie -mod=readonly -modcacherw

install: build
	mkdir -p $(PREFIX)/bin
	cp -f servus $(PREFIX)/bin

install-service:
	cp -f systemd/servus.service /etc/systemd/system/
