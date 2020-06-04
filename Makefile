PREFIX=/usr/local
BINDIR=$(PREFIX)/bin

fishpaste: main.go
	go build .

install: fishpaste
	install -o root -g wheel -m 0755 fishpaste $(BINDIR)
