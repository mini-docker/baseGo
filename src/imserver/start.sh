#!/bin/bash -x 
cd cmd/comet
go run main.go
cd ../../cmd/job
go run main.go
cd ../../cmd/logic
go run main.go


