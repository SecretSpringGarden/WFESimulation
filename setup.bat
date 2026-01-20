@echo off
REM Setup script for Workforce AI Transition Simulator (Windows)

echo Setting up Workforce AI Transition Simulator...

REM Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: Go is not installed. Please install Go 1.21 or higher from https://golang.org/dl/
    exit /b 1
)

REM Check Go version
for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i
echo Found Go version: %GO_VERSION%

REM Download dependencies
echo Downloading dependencies...
go mod download

REM Run tests to verify setup
echo Running tests...
go test ./...

if %ERRORLEVEL% EQU 0 (
    echo Setup complete! All tests passed.
    echo.
    echo To build the simulator, run: go build -o simulator.exe ./cmd/simulator
    echo To run tests, run: go test ./...
) else (
    echo Setup complete, but some tests failed. Please review the output above.
)
