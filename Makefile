PREFIX=/usr/local

format:
	gofmt -s -w .

build: format
	go build -trimpath -buildmode=pie -mod=readonly -modcacherw

install: build
	mkdir -p $(PREFIX)/bin
	cp -f servus $(PREFIX)/bin

install-service:
	cp -f systemd/servus.service /etc/systemd/system/
