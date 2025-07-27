#!/bin/bash

echo "Testing local checkout server..."
echo ""

# Test health endpoint
echo "1. Testing health endpoint:"
echo "   GET http://localhost:8080/health"
echo ""
curl -s "http://localhost:8080/health" | jq . 2>/dev/null || curl -s "http://localhost:8080/health"
echo ""
echo ""

# Test checkout endpoint
echo "2. Testing checkout endpoint:"
echo "   GET http://localhost:8080/checkout/123e4567-e89b-12d3-a456-426614174000"
echo ""
curl -s "http://localhost:8080/checkout/123e4567-e89b-12d3-a456-426614174000" | jq . 2>/dev/null || curl -s "http://localhost:8080/checkout/123e4567-e89b-12d3-a456-426614174000"
echo ""
echo ""

# Test with parameters
echo "3. Testing with parameters:"
echo "   GET http://localhost:8080/checkout/123e4567-e89b-12d3-a456-426614174000?browser=Chrome&os=Windows&country=BR"
echo ""
curl -s "http://localhost:8080/checkout/123e4567-e89b-12d3-a456-426614174000?browser=Chrome&os=Windows&country=BR&aff=test123" | jq . 2>/dev/null || curl -s "http://localhost:8080/checkout/123e4567-e89b-12d3-a456-426614174000?browser=Chrome&os=Windows&country=BR&aff=test123"
echo ""
echo ""

echo "✅ Tests completed!"
echo ""
echo "If you see JSON responses above, the server is working correctly."
echo "If you see connection errors, make sure the server is running with: ./bin/checkout-local.exe"