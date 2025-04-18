@echo off
if not exist "%~dp0bin" mkdir "%~dp0bin"
go build -o "%~dp0bin"
