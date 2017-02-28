@echo off 
setlocal EnableDelayedExpansion
set exts[windows]=.exe

if [%1] == [] (
    goto :error
) else (
    goto build
)

:error
echo Please, specify commit hash or tag.
goto :exit

:build
for /f %%i in ('git rev-parse %1') do set VERSION=%%i
for /f %%i in ('git tag --points-at %VERSION%') do set VERSION=%%i
set CGO_ENABLED=0
for %%O in (linux windows darwin freebsd) do (
    for %%A in (386 amd64) do (
        set GOOS=%%O
        set GOARCH=%%A
        go build -ldflags "-X main.version=%VERSION%" -o bombardier-%%O-%%A!exts[%%O]!
    )
)

:exit