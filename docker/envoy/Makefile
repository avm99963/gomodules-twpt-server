MUTABLE_VERSION ?= latest
VERSION ?= $(shell git rev-parse --short HEAD)

IMAGE_INTERNAL := docker-registry.corp.avm99963.com/twpt-server-envoy

docker-internal: Dockerfile
	docker build --force-rm -f Dockerfile --tag $(IMAGE_INTERNAL):$(VERSION) .
	docker tag $(IMAGE_INTERNAL):$(VERSION) $(IMAGE_INTERNAL):$(MUTABLE_VERSION)

push-internal: docker-internal
	docker push $(IMAGE_INTERNAL):$(VERSION)
	docker push $(IMAGE_INTERNAL):$(MUTABLE_VERSION)
