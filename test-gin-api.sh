#!/bin/bash

# Test script for Gin-based checkout API

BASE_URL="http://localhost:8080"
TEST_OFFER_UUID="123e4567-e89b-12d3-a456-426614174000"

echo "==================================="
echo "Testing Gin-based Checkout API"
echo "==================================="
echo ""

# Function to check if server is running
check_server() {
    if ! curl -s "$BASE_URL/health" > /dev/null; then
        echo "❌ Server is not running on $BASE_URL"
        echo "Please start the server first:"
        echo "  ./checkout-gin"
        echo ""
        exit 1
    fi
}

# Function to make API call and format output
api_call() {
    local endpoint="$1"
    local description="$2"
    
    echo "🧪 Testing: $description"
    echo "📍 Endpoint: $endpoint"
    echo "⏳ Making request..."
    
    response=$(curl -s -w "\n%{http_code}" "$endpoint" \
        -H "User-Agent: Test-Agent/1.0" \
        -H "Cookie: aff.product123=affiliate-uuid; _fbp=fb.1.123456789; _gcl_au=gcl.123456789")
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" = "200" ]; then
        echo "✅ Status: $http_code (Success)"
        echo "📄 Response:"
        echo "$body" | jq . 2>/dev/null || echo "$body"
    else
        echo "❌ Status: $http_code (Error)"
        echo "📄 Response:"
        echo "$body"
    fi
    
    echo ""
    echo "-----------------------------------"
    echo ""
}

# Check if jq is available for JSON formatting
if ! command -v jq &> /dev/null; then
    echo "⚠️  Note: 'jq' is not installed. JSON responses will not be formatted."
    echo "   Install jq for better output formatting:"
    echo "   - Ubuntu/Debian: sudo apt-get install jq"
    echo "   - macOS: brew install jq"
    echo "   - Windows: choco install jq"
    echo ""
fi

# Check if server is running
echo "🔍 Checking if server is running..."
check_server
echo "✅ Server is running!"
echo ""

# Test 1: Health Check
api_call "$BASE_URL/health" "Health Check"

# Test 2: Basic checkout request (new API format)
api_call "$BASE_URL/api/v1/checkout/$TEST_OFFER_UUID" "Basic Checkout (New API)"

# Test 3: Basic checkout request (legacy format)
api_call "$BASE_URL/checkout/$TEST_OFFER_UUID" "Basic Checkout (Legacy API)"

# Test 4: Checkout with query parameters
query_params="isMobile=false&country=BR&state=SP&city=SaoPaulo&utm_source=google&utm_medium=cpc&utm_campaign=test&aff=affiliate123&fbclid=fb123&gclid=gc123"
api_call "$BASE_URL/api/v1/checkout/$TEST_OFFER_UUID?$query_params" "Checkout with Parameters"

# Test 5: Invalid UUID
api_call "$BASE_URL/api/v1/checkout/invalid-uuid" "Invalid UUID (Should fail)"

# Test 6: Missing UUID
api_call "$BASE_URL/api/v1/checkout/" "Missing UUID (Should fail)"

# Test 7: Non-existent endpoint
api_call "$BASE_URL/non-existent" "Non-existent Endpoint (Should fail)"

echo "🏁 Testing completed!"
echo ""
echo "📊 Summary:"
echo "- Health check endpoint tested"
echo "- Main checkout endpoint tested (both formats)"
echo "- Parameter parsing tested"
echo "- Error handling tested"
echo ""
echo "💡 Tips:"
echo "- Check server logs for detailed request information"
echo "- Use different offer UUIDs to test various scenarios"
echo "- Test with real data when ready"
echo ""
echo "🔧 Customization:"
echo "- Edit BASE_URL to test different servers"
echo "- Edit TEST_OFFER_UUID to test specific offers"
echo "- Add more test cases as needed" 