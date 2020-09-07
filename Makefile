ifeq ($(DEBUG), true)
	COMPILEFLAGS := -gcflags "all=-N -l"
else
	COMPILEFLAGS :=
endif

BUILDFLAGS := $(COMPILEFLAGS)

.PHONY: build
build:
	go build -o build/multer -i cmd/multer/main.go $(BUILDFLAGS) 

.PHONY: test
test:
	go test -race --tags="testing" ./...

.PHONY: golangci-lint
golangci-lint:
	golangci-lint run ./...

.PHONY: run
run:
	build/multer -c deploy/multer.example.json