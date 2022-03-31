@echo off
go build -ldflags "-s -w" -trimpath

echo Finished