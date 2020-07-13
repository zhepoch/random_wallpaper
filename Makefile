SHELL := /bin/bash
BASEDIR = $(shell pwd)

# build with version infos
versionDir = "main"
gitTag = $(shell if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
buildDate = $(shell TZ=Asia/Shanghai date +%FT%T%z)
gitCommit = $(shell git log --pretty=format:'%H' -n 1)
gitTreeState = $(shell if git status|grep -q 'clean'; then echo clean; else echo dirty; fi)

ldflags="-w -X ${versionDir}.gitTag=${gitTag} -X ${versionDir}.buildDate=${buildDate} -X ${versionDir}.gitCommit=${gitCommit} -X ${versionDir}.gitTreeState=${gitTreeState}"


all: gotool
	@go build -v -o bin/random_wallpaper -ldflags ${ldflags} -gcflags=-trimpath=${BASEDIR} -asmflags=-trimpath=${BASEDIR}.
clean:
	rm -f bin/random_wallpaper
	find . -name "[._]*.s[a-w][a-z]" | xargs rm -f {}
gotool:
	gofmt -w .
	go vet . | grep -v vendor;true
help:
	@echo "make - compile the source code"
	@echo "make clean - remove binary file and vim swp files"
	@echo "make gotool - run go tool 'fmt' and 'vet'"
	@echo "make install - to install into local machine"
install:
	cp bin/random_wallpaper /usr/local/bin/random_wallpaper
ci-test:
	./bin/random_wallpaper --version

.PHONY: clean gotool install help ci-test
