@echo off
echo Testing local checkout server...

REM Check if server is running
echo Checking if server is running on port 8080...
netstat -an | findstr ":8080" >nul
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo ❌ Server is not running on port 8080!
    echo.
    echo Please start the server first:
    echo    bin\checkout-local.exe
    echo.
    echo Or run: build-local.bat then start the server
    exit /b 1
)

echo ✅ Server is running!
echo.

REM Test health endpoint
echo Testing health endpoint...
curl -s "http://localhost:8080/health" 2>nul
if %ERRORLEVEL% EQU 0 (
    echo.
    echo ✅ Health check passed!
) else (
    echo ❌ Health check failed! Make sure curl is installed or use PowerShell version.
)

echo.
echo Testing checkout endpoint...

REM Test checkout endpoint with sample data
curl -s "http://localhost:8080/checkout/123e4567-e89b-12d3-a456-426614174000?browser=Chrome&os=Windows&country=BR&utm_source=google" 2>nul

if %ERRORLEVEL% EQU 0 (
    echo.
    echo.
    echo ✅ Checkout test completed!
    echo.
    echo You can also test manually with these URLs:
    echo.
    echo Basic test:
    echo http://localhost:8080/checkout/123e4567-e89b-12d3-a456-426614174000
    echo.
    echo With parameters:
    echo http://localhost:8080/checkout/123e4567-e89b-12d3-a456-426614174000?browser=Chrome^&os=Windows^&country=BR^&utm_source=google^&aff=affiliate123
    echo.
) else (
    echo.
    echo ❌ Checkout test failed! 
    echo Try the PowerShell version: test-local.ps1
)

pause