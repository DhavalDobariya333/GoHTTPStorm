#!/bin/bash

cd GoHTTPStorm

go mod init gohttpstorm.go

go get github.com/valyala/fasthttp

go run gohttpstorm.go
