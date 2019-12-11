PREFIX=/usr/local

build: bindata format
	go build

format:
	gofmt -s -w .

bindata:
	go-bindata-assetfs -nometadata html assets

install:
	mkdir -p $(PREFIX)/bin
	cp -f servus $(PREFIX)/bin

install-service:
	cp -f systemd/servus.service /etc/systemd/system/
