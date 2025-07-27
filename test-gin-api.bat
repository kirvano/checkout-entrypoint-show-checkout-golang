@echo off
setlocal enabledelayedexpansion

REM Test script for Gin-based checkout API (Windows)

set BASE_URL=http://localhost:8080
set TEST_OFFER_UUID=123e4567-e89b-12d3-a456-426614174000

echo ===================================
echo Testing Gin-based Checkout API
echo ===================================
echo.

REM Check if curl is available
curl --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ curl is not available. Please install curl first.
    echo You can download it from: https://curl.se/download.html
    echo Or use Windows 10/11 built-in curl.
    pause
    exit /b 1
)

REM Function to check if server is running
echo 🔍 Checking if server is running...
curl -s %BASE_URL%/health >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Server is not running on %BASE_URL%
    echo Please start the server first:
    echo   checkout-gin.exe
    echo.
    pause
    exit /b 1
)
echo ✅ Server is running!
echo.

REM Test 1: Health Check
echo 🧪 Testing: Health Check
echo 📍 Endpoint: %BASE_URL%/health
echo ⏳ Making request...
curl -s %BASE_URL%/health
echo.
echo ✅ Health check completed
echo.
echo -----------------------------------
echo.

REM Test 2: Basic checkout request (new API format)
echo 🧪 Testing: Basic Checkout (New API)
echo 📍 Endpoint: %BASE_URL%/api/v1/checkout/%TEST_OFFER_UUID%
echo ⏳ Making request...
curl -s "%BASE_URL%/api/v1/checkout/%TEST_OFFER_UUID%" ^
    -H "User-Agent: Test-Agent/1.0" ^
    -H "Cookie: aff.product123=affiliate-uuid; _fbp=fb.1.123456789"
echo.
echo ✅ New API format tested
echo.
echo -----------------------------------
echo.

REM Test 3: Basic checkout request (legacy format)
echo 🧪 Testing: Basic Checkout (Legacy API)
echo 📍 Endpoint: %BASE_URL%/checkout/%TEST_OFFER_UUID%
echo ⏳ Making request...
curl -s "%BASE_URL%/checkout/%TEST_OFFER_UUID%" ^
    -H "User-Agent: Test-Agent/1.0"
echo.
echo ✅ Legacy API format tested
echo.
echo -----------------------------------
echo.

REM Test 4: Checkout with query parameters
set "query_params=isMobile=false&country=BR&state=SP&city=SaoPaulo&utm_source=google&utm_medium=cpc&utm_campaign=test&aff=affiliate123&fbclid=fb123&gclid=gc123"
echo 🧪 Testing: Checkout with Parameters
echo 📍 Endpoint: %BASE_URL%/api/v1/checkout/%TEST_OFFER_UUID%?%query_params%
echo ⏳ Making request...
curl -s "%BASE_URL%/api/v1/checkout/%TEST_OFFER_UUID%?%query_params%" ^
    -H "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
echo.
echo ✅ Parameters tested
echo.
echo -----------------------------------
echo.

REM Test 5: Invalid UUID
echo 🧪 Testing: Invalid UUID (Should fail)
echo 📍 Endpoint: %BASE_URL%/api/v1/checkout/invalid-uuid
echo ⏳ Making request...
curl -s "%BASE_URL%/api/v1/checkout/invalid-uuid"
echo.
echo ✅ Invalid UUID error handling tested
echo.
echo -----------------------------------
echo.

REM Test 6: Non-existent endpoint
echo 🧪 Testing: Non-existent Endpoint (Should fail)
echo 📍 Endpoint: %BASE_URL%/non-existent
echo ⏳ Making request...
curl -s "%BASE_URL%/non-existent"
echo.
echo ✅ 404 error handling tested
echo.
echo -----------------------------------
echo.

echo 🏁 Testing completed!
echo.
echo 📊 Summary:
echo - Health check endpoint tested
echo - Main checkout endpoint tested (both formats)
echo - Parameter parsing tested
echo - Error handling tested
echo.
echo 💡 Tips:
echo - Check server logs for detailed request information
echo - Use different offer UUIDs to test various scenarios
echo - Test with real data when ready
echo.
echo 🔧 Customization:
echo - Edit BASE_URL to test different servers
echo - Edit TEST_OFFER_UUID to test specific offers
echo - Add more test cases as needed
echo.
pause 