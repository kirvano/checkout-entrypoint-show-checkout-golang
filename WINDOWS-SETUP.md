# Windows Setup Guide

This guide provides Windows-specific instructions for building, running, and testing the Go checkout application.

## 🚀 Quick Start for Windows

### 1. Build Local Development Server
```cmd
build-local.bat
```
This creates `bin\checkout-local.exe` for local testing.

### 2. Start the Server
```cmd
bin\checkout-local.exe
```
Server starts on `http://localhost:8080`

### 3. Test the Server
In another terminal:
```cmd
REM Basic test with curl
test-local.bat

REM Or use PowerShell for better output
powershell -ExecutionPolicy Bypass -File test-local.ps1
```

### 4. Build for AWS Lambda
```cmd
build-lambda.bat
```
This creates `bin\checkout-lambda.zip` ready for AWS deployment.

## 📋 Available Scripts

| Script | Purpose | Output |
|--------|---------|--------|
| `build-local.bat` | Build local HTTP server | `bin\checkout-local.exe` |
| `build-lambda.bat` | Build Lambda deployment package | `bin\checkout-lambda.zip` |
| `test-local.bat` | Test local server (curl) | Console output |
| `test-local.ps1` | Test local server (PowerShell) | Detailed output + JSON file |

## 🧪 Testing Examples

### Manual Testing with Browser
- Health check: http://localhost:8080/health
- Basic checkout: http://localhost:8080/checkout/123e4567-e89b-12d3-a456-426614174000
- With parameters: http://localhost:8080/checkout/123e4567-e89b-12d3-a456-426614174000?browser=Chrome&os=Windows&country=BR

### PowerShell Testing
```powershell
# Test health endpoint
Invoke-RestMethod -Uri "http://localhost:8080/health"

# Test checkout with parameters
$response = Invoke-RestMethod -Uri "http://localhost:8080/checkout/123e4567-e89b-12d3-a456-426614174000?browser=Chrome&os=Windows&country=BR&aff=test123"
$response | ConvertTo-Json -Depth 10
```

### cURL Testing
```cmd
REM Health check
curl "http://localhost:8080/health"

REM Checkout test
curl "http://localhost:8080/checkout/123e4567-e89b-12d3-a456-426614174000?browser=Chrome&os=Windows&country=BR"
```

## 🔧 Development Workflow

1. **Edit Code** - Make changes to Go files
2. **Build** - Run `build-local.bat`
3. **Test** - Run `test-local.ps1` or test manually
4. **Deploy** - Run `build-lambda.bat` and upload to AWS

## 📁 Project Structure

```
checkout-go/
├── bin/                          # Build outputs
│   ├── checkout-local.exe        # Local server (after build-local.bat)
│   └── checkout-lambda.zip       # Lambda package (after build-lambda.bat)
├── cmd/
│   ├── lambda/main.go            # AWS Lambda entry point
│   └── local/main.go             # Local HTTP server
├── build-local.bat               # Build local server
├── build-lambda.bat              # Build Lambda package
├── test-local.bat                # Test script (curl)
├── test-local.ps1                # Test script (PowerShell)
└── ...rest of the project
```

## 🐛 Troubleshooting

### Build Issues
```cmd
REM Check Go installation
go version

REM Install dependencies
go mod tidy

REM Clean and rebuild
del bin\*.exe
build-local.bat
```

### Server Issues
```cmd
REM Check if port 8080 is in use
netstat -an | findstr ":8080"

REM Kill process using port 8080
taskkill /f /im checkout-local.exe
```

### PowerShell Execution Policy
If you get execution policy errors:
```powershell
# Run this once as Administrator
Set-ExecutionPolicy RemoteSigned

# Or run scripts with bypass
powershell -ExecutionPolicy Bypass -File test-local.ps1
```

## 📊 Expected Performance

- **Startup Time**: <1 second
- **Response Time**: <50ms for healthy requests
- **Memory Usage**: ~20-30MB for local server
- **Build Time**: ~5-10 seconds

## 🚀 AWS Deployment

After running `build-lambda.bat`:

1. Upload `bin\checkout-lambda.zip` to AWS Lambda
2. Set handler to: `bootstrap`
3. Set runtime to: `provided.al2` or `provided.al2023`
4. Configure environment variables as needed

## 💡 Tips

- Use **PowerShell** scripts for better output and error handling
- The **local server** mimics AWS Lambda behavior for testing
- **Test files** are saved with timestamps for comparison
- **Build scripts** include error checking and helpful output

---

Happy coding! 🎉