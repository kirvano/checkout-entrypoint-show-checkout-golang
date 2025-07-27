# 🚀 Quick Start Guide

Choose your preferred method based on your terminal:

## For Git Bash / WSL / Linux

### 1. Build the server
```bash
chmod +x build-local.sh
./build-local.sh
```

### 2. Run the server
```bash
./bin/checkout-local.exe
```

### 3. Test the server (in another terminal)
```bash
chmod +x test-local.sh
./test-local.sh
```

### 4. Build for Lambda
```bash
chmod +x build-lambda.sh
./build-lambda.sh
```

## For Windows PowerShell / Command Prompt

### 1. Build the server
```cmd
build-local.bat
```

### 2. Run the server
```cmd
bin\checkout-local.exe
```

### 3. Test the server (in another terminal)
```powershell
.\test-local.ps1
```

### 4. Build for Lambda
```cmd
build-lambda.bat
```

## 🔍 Troubleshooting

### Server exits immediately?

The improved server now shows detailed startup logs:
- ✅ If you see "Press Ctrl+C to stop the server" - server is running correctly
- ❌ If you see an error before that - there's a startup issue

### Can't execute scripts?

**Git Bash:**
```bash
chmod +x *.sh
```

**PowerShell:**
```powershell
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Test manually

Open browser and visit:
- Health: http://localhost:8080/health
- Basic test: http://localhost:8080/checkout/123e4567-e89b-12d3-a456-426614174000

## 📊 Expected Output

When server starts correctly, you should see:
```
Starting local checkout server on :8080
Initializing dependency injection container...
DI container initialized successfully
ShowCheckoutUseCase retrieved successfully
Server endpoints registered:
  - GET /health
  - GET /checkout/{uuid}

Server running at http://localhost:8080
Example: http://localhost:8080/checkout/123e4567-e89b-12d3-a456-426614174000

Press Ctrl+C to stop the server
```

If the server stops before showing "Press Ctrl+C", check the error message above that line.

---

**Next:** Try rebuilding with the improved server and let me know what you see!