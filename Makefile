BIN=freeipa-dhcp
OUTPUT_DIR=build
PROJECT=freeipa-dhcp
PKG_NAME=$(PROJECT)
PKG_VERSION=0.1.0

build: $(BIN)

$(BIN):
	mkdir -p build
	go build -o $(OUTPUT_DIR)/$(BIN)

$(BIN)-linux:
	mkdir -p build
	GOOS=linux go build -o $(OUTPUT_DIR)/$(BIN)

deb: $(BIN)-linux
	docker run --rm -v $(PWD):/build/$(PROJECT) alanfranz/fpm-within-docker:debian-stretch \
	fpm -s dir -t deb -n $(PKG_NAME) -v $(PKG_VERSION) -p /build/$(PROJECT)/build/ \
	--config-files /etc/default/$(BIN) \
	/build/$(PROJECT)/build/$(BIN)=/usr/bin/ \
	/build/$(PROJECT)/conf/$(BIN)=/etc/default/ \
	/build/$(PROJECT)/conf/$(BIN).service=/lib/systemd/system/

clean:
	rm build/*
