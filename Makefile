# Platform builds via Wails. Default `build` targets the host OS/arch.
#
# macOS note: Go 1.27 emits macOS 13 objects; Wails defaults CGO to 10.13.
# Darwin targets set MACOSX_DEPLOYMENT_TARGET / CGO_* so the linker matches.

WAILS ?= wails
MACOSX_DEPLOYMENT_TARGET ?= 13.0

.PHONY: build build-debug build-darwin build-darwin-arm64 build-darwin-amd64 \
	build-linux build-windows build-all dev help

help:
	@echo "Targets:"
	@echo "  make build              Native production build (host OS/arch)"
	@echo "  make build-debug        Native debug build"
	@echo "  make build-darwin       macOS Apple Silicon (darwin/arm64)"
	@echo "  make build-darwin-amd64 macOS Intel (darwin/amd64)"
	@echo "  make build-linux        Linux amd64"
	@echo "  make build-windows      Windows amd64"
	@echo "  make build-all          darwin/arm64 + linux/amd64 + windows/amd64"
	@echo "  make dev                wails dev (live reload)"

build:
	$(WAILS) build

build-debug:
	$(WAILS) build -debug

build-darwin: build-darwin-arm64

build-darwin-arm64:
	MACOSX_DEPLOYMENT_TARGET=$(MACOSX_DEPLOYMENT_TARGET) \
	CGO_CFLAGS="-mmacosx-version-min=$(MACOSX_DEPLOYMENT_TARGET)" \
	CGO_LDFLAGS="-mmacosx-version-min=$(MACOSX_DEPLOYMENT_TARGET)" \
	$(WAILS) build -platform darwin/arm64

build-darwin-amd64:
	MACOSX_DEPLOYMENT_TARGET=$(MACOSX_DEPLOYMENT_TARGET) \
	CGO_CFLAGS="-mmacosx-version-min=$(MACOSX_DEPLOYMENT_TARGET)" \
	CGO_LDFLAGS="-mmacosx-version-min=$(MACOSX_DEPLOYMENT_TARGET)" \
	$(WAILS) build -platform darwin/amd64

build-linux:
	$(WAILS) build -platform linux/amd64

build-windows:
	$(WAILS) build -platform windows/amd64

build-all: build-darwin-arm64 build-linux build-windows

dev:
	$(WAILS) dev
