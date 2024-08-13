############################ develop #############################################
start:
	go run .
	#arun sh -c 'swag init -g pkg/main.go && go run ./pkg'
	#sh -c 'swag init --parseDependency --parseInternal --parseDepth 1 && go run .'

init: 
	#go install github.com/swaggo/swag/cmd/swag@latest -insecure
	#go install github.com/swaggo/swag/cmd/swag@v1.8.1
	go install github.com/swaggo/swag/cmd/swag@latest
generate:
	#//go:generate swag init --parseDependency --parseInternal -g description.go -o ./docs
	go generate
doc: 
	# swag init -g cmd/main.go
	# swag init --parseDependency --parseInternal --parseDepth 1 #&& open http://m:4500/swagger/index.html
	swag init  #&& open http://m:4500/swagger/index.html

############################# k8s ###############################################
docker:
	docker build -t ginapp .
	docker run -p 4501:4500 --name gin1 --rm -it  ginapp ./main -p 4500

deploy-step:
	zsh k8s/deploy.sh

######################### bench/perf #################################
benchcpu:
	go-wrk  -d=50 -c=50  http://localhost:4500/cpu/5

profcpu:
	go tool pprof http://127.0.0.1:4500/debug/pprof/profile

######################### test & pkg ##########################################
msg?=
.ONESHELL:
gitcheck:
	if [[ "$(msg)" = "" ]] ; then echo "Usage: make pkg msg='commit msg'";exit 20; fi

.ONESHELL:
pkg: gitcheck test
	{ hash newversion.py 2>/dev/null && newversion.py version;} ;  { echo version `cat version`; }
	git commit -am "$(msg)"
	#jfrog "rt" "go-publish" "go-pl" $$(cat version) "--url=$$GOPROXY_API" --user=$$GOPROXY_USER --apikey=$$GOPROXY_PASS
	v=`cat version` && git tag "$$v" && git push origin "$$v" && git push origin HEAD

LDFLAGS=-ldflags="-s -w -X ginapp/conf.BuildDate=$(shell date -Iseconds) -X ginapp/conf.BuildBranch=$(shell git rev-parse --abbrev-ref HEAD)"

install:
	# go install .
	go install $(LDFLAGS) 
	#go build -o ~/bin/ginapp .

build:
	#go build $(LDFLAGS) -o ginapp main.go
	go build $(LDFLAGS) -o ginapp

.ONESHELL:
build-linux: test
	# gh auth -h ; https://github.com/settings/applications
	{ hash newversion.py 2>/dev/null && newversion.py version;} ;  { echo version `cat version`; }
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/ginapp-linux-amd64 
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/ginapp-darwin-arm64
	v=`cat version` && \
	gh release create $$v --notes $$v dist/ginapp-linux-amd64 dist/ginapp-darwin-arm64


.PHONY: test
test: 
	go test -race -coverprofile cover.out -coverpkg "./..." -failfast ./...
cover: test
	go tool cover -html=cover.out
race: 
	go test -race -failfast ./...
fmt:
	gofmt -s -w .
