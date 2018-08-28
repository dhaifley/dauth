NOVENDOR = $(shell go list ./... | grep -v vendor | grep -v node_modules)
NOVENDOR_LINTER = $(shell go list ./... | grep -v vendor | grep -v ptypes | grep -v node_modules)

all: build

fix:
	go fix $(NOVENDOR)
.PHONY: fix

vet:
	go vet $(NOVENDOR)
.PHONY: vet

lint:
	printf "%s\n" "$(NOVENDOR)" | xargs -I {} sh -c 'golint -set_exit_status {}'
.PHONY: lint

test:
	go test -v -cover $(NOVENDOR)
.PHONY: test

metalinter:
	gometalinter --config .gometalinter.json $(NOVENDOR_LINTER)
.PHONY: metalinter

clean:
	rm -rf ./bin
.PHONY: clean

build: clean fix vet lint test
	mkdir bin
	GOOS=linux GOARCH=386 go build -v -o ./bin/dauth main.go
	# Uncomment if Windows build is needed.
	# GOOS=windows GOARCH=amd64 go build -v -o ./bin/dauth.exe main.go
.PHONY: build

docker: build
	docker build -f Dockerfile -t gcr.io/rf-services/dauth:latest .
	docker push gcr.io/rf-services/dauth:latest
	# Uncomment if there is a script container to build.
	# docker build -f script/Dockerfile -t dhaifley/dauth_scripts:latest .
	# docker push dhaifley/dauth_scripts:latest
	docker system prune --volumes -f
.PHONY: docker

deploy: docker
	scp update.sh dauth:dauth/update.sh
	scp docker-compose.yml dauth:dauth/docker-compose.yml
.PHONY: deploy
