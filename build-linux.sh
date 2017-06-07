#!/bin/bash
export GOPATH=`pwd`
export GOARCH=amd64
export GOOS=linux
cd bin

go build -o server -ldflags "-w -s" ../src/server.go
go build -o robot -ldflags "-w -s" ../src/robot.go


upx -9 server
upx -9 robot

read -p "Press any key to continue." var


