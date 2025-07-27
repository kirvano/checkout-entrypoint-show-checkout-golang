# PowerShell script to test local checkout server
Write-Host "Testing local checkout server..." -ForegroundColor Green

# Check if server is running
Write-Host "Checking if server is running on port 8080..."
$connection = Test-NetConnection -ComputerName localhost -Port 8080 -WarningAction SilentlyContinue
if (-not $connection.TcpTestSucceeded) {
    Write-Host ""
    Write-Host "❌ Server is not running on port 8080!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please start the server first:" -ForegroundColor Yellow
    Write-Host "   bin\checkout-local.exe" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Or run: build-local.bat then start the server" -ForegroundColor Yellow
    exit 1
}

Write-Host "✅ Server is running!" -ForegroundColor Green
Write-Host ""

# Test health endpoint
Write-Host "Testing health endpoint..."
try {
    $healthResponse = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get
    Write-Host "Health Response:" -ForegroundColor Cyan
    $healthResponse | ConvertTo-Json
    Write-Host "✅ Health check passed!" -ForegroundColor Green
} catch {
    Write-Host "❌ Health check failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "Testing checkout endpoint..."

# Test checkout endpoint with sample data
$testUrl = "http://localhost:8080/checkout/123e4567-e89b-12d3-a456-426614174000?browser=Chrome&os=Windows&country=BR&utm_source=google&aff=affiliate123"

try {
    Write-Host "Making request to:" -ForegroundColor Cyan
    Write-Host $testUrl -ForegroundColor Yellow
    Write-Host ""
    
    $checkoutResponse = Invoke-RestMethod -Uri $testUrl -Method Get
    Write-Host "✅ Checkout test completed!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Response Preview:" -ForegroundColor Cyan
    Write-Host "Billing Type: $($checkoutResponse.billing_type)" -ForegroundColor White
    Write-Host "Is Free: $($checkoutResponse.is_free)" -ForegroundColor White
    Write-Host "Product Name: $($checkoutResponse.product.name)" -ForegroundColor White
    Write-Host "Product Price: $($checkoutResponse.product.price)" -ForegroundColor White
    Write-Host ""
    
    # Save full response to file
    $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
    $outputFile = "test-response-$timestamp.json"
    $checkoutResponse | ConvertTo-Json -Depth 10 | Out-File -FilePath $outputFile -Encoding UTF8
    Write-Host "Full response saved to: $outputFile" -ForegroundColor Green
    
} catch {
    Write-Host "❌ Checkout test failed: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        Write-Host "Status Code: $($_.Exception.Response.StatusCode)" -ForegroundColor Red
        Write-Host "Status Description: $($_.Exception.Response.StatusDescription)" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "You can also test manually with these URLs:" -ForegroundColor Yellow
Write-Host ""
Write-Host "Basic test:" -ForegroundColor Cyan
Write-Host "http://localhost:8080/checkout/123e4567-e89b-12d3-a456-426614174000" -ForegroundColor White
Write-Host ""
Write-Host "With parameters:" -ForegroundColor Cyan  
Write-Host "http://localhost:8080/checkout/123e4567-e89b-12d3-a456-426614174000?browser=Chrome&os=Windows&country=BR&utm_source=google&aff=affiliate123" -ForegroundColor White
Write-Host ""

Read-Host "Press Enter to continue..."