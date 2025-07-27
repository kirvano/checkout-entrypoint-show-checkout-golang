@echo off

REM Build script for Gin-based checkout server (Windows)

echo Building Gin-based checkout server...

REM Clean previous builds
del /q checkout-gin.exe 2>nul
del /q bin\checkout-gin*.exe 2>nul

REM Build for current platform (Windows)
echo Building for Windows...
go build -o checkout-gin.exe ./cmd/gin/

REM Build for Linux (useful for Docker/deployment)
echo Building for Linux amd64...
set GOOS=linux
set GOARCH=amd64
go build -o bin/checkout-gin-linux ./cmd/gin/

REM Build for macOS
echo Building for macOS amd64...
set GOOS=darwin
set GOARCH=amd64
go build -o bin/checkout-gin-darwin ./cmd/gin/

REM Reset environment variables
set GOOS=
set GOARCH=

echo Build completed!
echo Files generated:
echo   - checkout-gin.exe (Windows)
echo   - bin/checkout-gin-linux (Linux)
echo   - bin/checkout-gin-darwin (macOS)
echo.
echo To run the server:
echo   checkout-gin.exe
echo.
echo Or with custom port:
echo   set PORT=3000 ^&^& checkout-gin.exe 