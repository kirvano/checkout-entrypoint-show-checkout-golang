@echo off
echo Building Lambda function for AWS...

REM Create bin directory if it doesn't exist
if not exist bin mkdir bin

REM Set environment variables for Linux build
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0

echo Building for Linux (AWS Lambda)...

REM Build the Lambda function
go build -ldflags="-s -w" -o bin/bootstrap cmd/lambda/main.go

if %ERRORLEVEL% EQU 0 (
    echo.
    echo Creating deployment package...
    
    REM Create zip file for Lambda deployment
    cd bin
    if exist checkout-lambda.zip del checkout-lambda.zip
    
    REM Use PowerShell to create zip (available on all Windows 10+ systems)
    powershell -Command "Compress-Archive -Path bootstrap -DestinationPath checkout-lambda.zip -Force"
    
    if exist checkout-lambda.zip (
        echo.
        echo ✅ Lambda package created successfully!
        echo.
        echo Package location: bin\checkout-lambda.zip
        echo Package size:
        dir checkout-lambda.zip | findstr checkout-lambda.zip
        echo.
        echo To deploy to AWS Lambda:
        echo 1. Upload bin\checkout-lambda.zip to your Lambda function
        echo 2. Set the handler to: bootstrap
        echo 3. Set runtime to: provided.al2 or provided.al2023
        echo.
    ) else (
        echo ❌ Failed to create zip package!
        exit /b 1
    )
    
    cd ..
) else (
    echo.
    echo ❌ Build failed!
    echo Check the error messages above.
    exit /b 1
)

REM Reset environment variables
set GOOS=
set GOARCH=
set CGO_ENABLED=