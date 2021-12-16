:: MedHash Tools
:: Copyright (c) 2021 GHIFARI160
:: MIT License

@echo off

SETLOCAL
if exist ..\out (
    echo Cleaning build directory
    rmdir /s /q ..\out
)

if exist ..\dist (
    echo Cleaning release directory
    rmdir /s /q ..\dist
)
ENDLOCAL
