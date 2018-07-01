IMAGE_TAG=v1alpha1
HASHED_FILE?=./test.txt
QUAY_PASS?=biggestsecret

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o decrypt-incremental main.go
	docker build -t quay.io/tamarakaufler/decrypt-incremental:$(IMAGE_TAG) .
	docker login quay.io -u tamarakaufler -p $(QUAY_PASS)
	docker push quay.io/tamarakaufler/decrypt-incremental:$(IMAGE_TAG)

dev:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o decrypt-incremental main.go
	docker build -t quay.io/tamarakaufler/decrypt-incremental:$(IMAGE_TAG) .

run:
	docker run \
	--name=decrypt-incremental \
	--rm \
	quay.io/tamarakaufler/decrypt-incremental:$(IMAGE_TAG) \
	-f=$(HASHED_FILE)
