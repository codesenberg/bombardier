@echo off
set GOOS=windows
set GOARCH=amd64
go build -o ../../bombardier-ui.exe
if %ERRORLEVEL% neq 0 (
    echo Build failed!
    exit /b %ERRORLEVEL%
)
echo Build completed successfully!
echo You can now run ..\..\bombardier-ui.exe