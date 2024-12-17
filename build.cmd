@echo off
if not exist "%~dp0bin" mkdir "%~dp0bin"
cd /d "%~dp0bin"
go build "%~dp0"