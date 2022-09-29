#!/bin/bash
echo "Running tracker"
go env -w GO111MODULE=off
go build 
./tracker 
