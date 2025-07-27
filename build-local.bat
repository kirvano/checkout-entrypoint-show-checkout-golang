@echo off
echo Building local development server...

REM Build the local server
go build -o bin/checkout-local.exe cmd/local/main.go

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ✅ Build successful!
    echo.
    echo To start the server, run:
    echo    bin\checkout-local.exe
    echo.
    echo Server will be available at: http://localhost:8080
    echo Test endpoint: http://localhost:8080/checkout/123e4567-e89b-12d3-a456-426614174000
    echo.
) else (
    echo.
    echo ❌ Build failed!
    echo Check the error messages above.
    exit /b 1
)