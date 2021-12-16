:: MedHash Tools
:: Copyright (c) 2021 GHIFARI160
:: MIT License

@echo off

SETLOCAL
set ver=0.3.0
set buildDir=..\out
set releaseDir=..\dist

set buildArgs=-trimpath
set ldArgs=-s -w -X github.com/ghifari160/medhash-tools/src/common.VERSION=%ver%

goto main

:: Build a given tool for a given platform and architecture
:: %~1: GOOS
:: %~2: GOARCH
:: %~3: Tool name
:: %~4: Suffix
:build_helper
set GOOS=%~1
set GOARCH=%~2
echo Building %~3 for %~1 / %~2 (Suffix: %~4)
go build %buildArgs% -ldflags="%ldArgs%" -o %buildDir%\%~3%~4 .\%~3
exit /B 0

:: Build all tools
:: %~1: GOOS
:: %~2: GOARCH
:: %~3: Suffix
:build_all
call :build_helper %~1 %~2 medhash-gen %~3
call :build_helper %~1 %~2 medhash-chk %~3
call :build_helper %~1 %~2 medhash-upgrade %~3
exit /B 0

:: Move binary to the release directory
:: %~1: Tool name
:: %~2: Platform
:: %~3: Architecture
:: %~4: Suffix
:move_helper
echo Moving %~1 to %releaseDir%\%ver%\%~1-%~2_%~3-%ver%%~4
move /y %buildDir%\%~1 %releaseDir%\%ver%\%~1-%~2_%~3-%ver%%~4
exit /B 0

:: Linux
:build_linux
call :build_all linux, 386, -386

call :move_helper medhash-gen-386, linux, 386
call :move_helper medhash-chk-386, linux, 386
call :move_helper medhash-upgrade-386, linux, 386
exit /B 0

:: macOS (build not supported)
:build_macos
echo Cannot build for macOS on non-macOS host machine
exit /B 0

:: Windows
:build_windows
call :build_all windows, 386, -386

call :move_helper medhash-gen-386, windows, x86, .exe
call :move_helper medhash-chk-386, windows, x86, .exe
call :move_helper medhash-upgrade-386, windows, x86, .exe
exit /B 0

:main
if not exist %buildDir% (
    echo Creating build directory
    mkdir %buildDir%
)

if not exist %releaseDir% (
    echo Creating release directory for v%ver%
    mkdir %releaseDir%\%ver%
)

echo %~1

if "%~1"=="linux" (
    call :build_linux
) else if "%~1"=="macos" (
    call :build_macos
) else if "%~1"=="darwin" (
    call :build_macos
) else if "%~1"=="windows" (
    call :build_windows
) else (
    call :build_linux
    call :build_windows
)

:end
echo Done

ENDLOCAL
