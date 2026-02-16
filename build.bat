@echo off
echo Building notify.exe with optimizations...

:: Set Go flags for optimization
set GOOS=windows
set GOARCH=amd64

:: Build with optimizations
:: -ldflags "-s -w" strips debug information and DWARF symbol table
:: -trimpath removes file system paths from binary
go build -ldflags="-s -w" -trimpath -o notify.exe main.go

if %ERRORLEVEL% EQU 0 (
    echo Build successful! notify.exe created.
    echo.
    echo Binary size:
    for %%I in (notify.exe) do echo %%~zI bytes
) else (
    echo Build failed!
    exit /b 1
)

pause