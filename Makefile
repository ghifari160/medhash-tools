BUILD_DIR ?= $(shell pwd)/out/build
OBJ_DIR ?= $(BUILD_DIR)/cache

.PHONY: default_target
default_target: all

$(BUILD_DIR)/medhash-gen:
	mkdir -p $(@D)
	go build -o $@ src/medhash-gen/*.go

$(BUILD_DIR)/medhash-chk:
	mkdir -p $(@D)
	go build -o $@ src/medhash-chk/*.go

$(BUILD_DIR)/medhash-upgrade:
	mkdir -p $(@D)
	go build -o $@ src/medhash-upgrade/*.go

.PHONY: all
all: $(BUILD_DIR)/medhash-gen $(BUILD_DIR)/medhash-chk $(BUILD_DIR)/medhash-upgrade

.PHONY: linux
linux:
	$(MAKE) -C . clean
	export GOOS=linux && export GOARCH=386 && $(MAKE) -C . all

.PHONY: macos
macos: darwin

.PHONY: darwin
darwin:
	$(MAKE) -C . clean
	export GOOS=darwin && export GOARCH=amd64 && $(MAKE) -C . all

.PHONY: windows
windows:
	$(MAKE) -C . clean
	export GOOS=windows && export GOARCH=386 && $(MAKE) -C . all
	mv $(BUILD_DIR)/medhash-gen $(BUILD_DIR)/medhash-gen.exe
	mv $(BUILD_DIR)/medhash-chk $(BUILD_DIR)/medhash-chk.exe
	mv $(BUILD_DIR)/medhash-upgrade $(BUILD_DIR)/medhash-upgrade.exe

.PHONY: clean
clean:
	-rm -rf $(BUILD_DIR) $(OBJ_DIR)
