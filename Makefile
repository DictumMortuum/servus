PREFIX=/usr/local

build: format
	go build

format:
	gofmt -s -w .

install:
	mkdir -p $(PREFIX)/bin
	cp -f servus $(PREFIX)/bin

install-service:
	cp -f systemd/servus.service /etc/systemd/system/
