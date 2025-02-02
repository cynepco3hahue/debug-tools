RUNTIME ?= podman
REPOOWNER ?= openshift-kni
IMAGENAME ?= debug-tools
IMAGETAG ?= latest

all: dist

.PHONY: build
build: dist

.PHONY: ci-job
ci-job: test-e2e-knit

outdir:
	mkdir -p _output || :

.PHONY: deps-update
deps-update:
	go mod tidy && go mod vendor

.PHONY: deps-clean
deps-clean:
	rm -rf vendor

.PHONY: dist
dist: binaries

.PHONY: binaries
binaries: outdir deps-update
	# go flags are set in here
	./hack/build-binaries.sh

.PHONY: clean
clean:
	rm -rf _output

.PHONY: image
image: binaries
	@echo "building image"
	$(RUNTIME) build -f Dockerfile -t quay.io/$(REPOOWNER)/$(IMAGENAME):$(IMAGETAG) .

.PHONY: push
push: image
	@echo "pushing image"
	$(RUNTIME) push quay.io/$(REPOOWNER)/$(IMAGENAME):$(IMAGETAG)

.PHONY: test-unit
test-unit: test-unit-pkg

.PHONY: test-unit-pkg
test-unit-pkg:
	go test ./pkg/...

.PHONY: test-e2e-knit
test-e2e-knit: binaries
	ginkgo test/e2e
