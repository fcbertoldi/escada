@echo off
set destination_dir="%ProgramFiles%\Escada"
md "%destination_dir%"
copy /Y "%~dp0bin\escada.exe" "%destination_dir%"
