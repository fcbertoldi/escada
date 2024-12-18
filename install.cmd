@echo off
md "%ProgramFiles%\Escada"
set GOBIN="%ProgramFiles%\Escada"
go install "%~dp0bin\escada"
