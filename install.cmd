@echo off
set destination_dir="%ProgramFiles%\Escada"
if not exist %destination_dir% mkdir %destination_dir%
copy /Y "%~dp0bin\escada.exe" %destination_dir%
