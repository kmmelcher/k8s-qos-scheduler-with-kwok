all: build
	sudo docker push giovannifs/qos-driven-scheduler:v0.2.5.12
	sudo docker pull giovannifs/qos-driven-scheduler:v0.2.5.12

build:
	sudo docker build . -t giovannifs/qos-driven-scheduler:v0.2.5.12

compile: format
	CGO_ENABLED=1 go build -race -ldflags '-w' -o bin/kube-scheduler cmd/main.go

format: import
	go fmt cmd/main.go

import: go.mod
	go mod tidy && go mod vendor
