VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
LDFLAGS = -X main.version=${VERSION}

build:
	go build -ldflags '$(LDFLAGS)' cmd/wallet-exporter.go

install:
	go install -ldflags '$(LDFLAGS)' cmd/wallet-exporter.go

lint:
	golangci-lint run --fix ./...

test:
	go test -coverpkg=./... -coverprofile cover.out -v ./...

coverage:
	go tool cover -html=cover.out


APP_NAME := wallet-exporter
DOCKER_REPO := choihocheol/$(APP_NAME)
VERSION := $(shell git describe --tags --abbrev=0)

PLATFORMS := linux/amd64,linux/arm64

run:
	go run cmd/wallet-exporter.go --config config.toml

docker-build:
	docker build -t $(APP_NAME):latest .

docker-push:
	@echo "Building image for all platforms: $(PLATFORMS)"
	docker buildx inspect multiarch-builder >/dev/null 2>&1 || docker buildx create --use --name multiarch-builder --bootstrap
	docker buildx use multiarch-builder
	docker buildx build \
		--platform $(PLATFORMS) \
		--tag ghcr.io/$(DOCKER_REPO):latest \
		--tag ghcr.io/$(DOCKER_REPO):$(VERSION) \
		--file Dockerfile \
		--push .
