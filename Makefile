BINARY = ticket-broker
GOARCH = amd64

VERSION?=$(shell git describe HEAD | tr -d 'v')
BUILD=$(shell git rev-parse HEAD)
PACKAGE=${BINARY}_${VERSION}_${GOARCH}.deb
SERVER?=ticket-broker.proctoru.com

LDFLAGS = -ldflags "-s -w -X main.VERSION=${VERSION} -X main.BUILD=${BUILD}"

.PHONY: linux darwin fmt clean prep package 

all: clean prep linux darwin package

linux:
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o build/linux-${GOARCH}/${BINARY} main.go

darwin:
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o build/darwin-${GOARCH}/${BINARY} main.go

package:
	fpm \
		-f \
		-s dir \
		-t deb \
		-n ${BINARY} \
		-v ${VERSION} \
		-p build \
		--after-install resources/after-install.sh \
		--after-upgrade resources/after-upgrade.sh \
		--before-remove resources/before-remove.sh \
		build/linux-${GOARCH}/${BINARY}=/usr/local/bin/${BINARY} \
		resources/ticket-broker.service=/lib/systemd/system/ticket-broker.service

fmt:
	go fmt

prep:
	mkdir -p build/darwin-${GOARCH}
	mkdir -p build/linux-${GOARCH}

clean:
	- rm -r build

