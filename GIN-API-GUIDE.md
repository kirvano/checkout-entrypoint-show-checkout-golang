# Gin-Based Checkout API Guide

This document explains how to use the new Gin-based API server for the checkout service.

## Overview

The Gin-based API provides a modern, high-performance HTTP server with the following features:

- **Fast routing** with Gin framework
- **Middleware support** for CORS, logging, recovery, and request tracking
- **Graceful shutdown** with proper signal handling
- **Environment-based configuration**
- **Structured logging** with request IDs
- **Backward compatibility** with existing endpoints

## Quick Start

### 1. Install Dependencies

First, add Gin to your project:

```bash
go get github.com/gin-gonic/gin
```

### 2. Build the Server

#### On Linux/macOS:
```bash
chmod +x build-gin.sh
./build-gin.sh
```

#### On Windows:
```bash
build-gin.bat
```

### 3. Run the Server

#### Default (port 8080):
```bash
./checkout-gin
```

#### Custom port:
```bash
PORT=3000 ./checkout-gin
```

#### With environment variables:
```bash
GIN_MODE=release PORT=8080 ./checkout-gin
```

## API Endpoints

### Health Check
```
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "checkout-api"
}
```

### Show Checkout (Main Endpoint)

#### New API format:
```
GET /api/v1/checkout/{offerUuid}
```

#### Legacy format (backward compatibility):
```
GET /checkout/{offerUuid}
```

**Parameters:**
- `offerUuid` (path) - The UUID of the checkout offer

**Query Parameters:**
- `isMobile` (boolean) - Whether the request is from a mobile device
- `ip` (string) - Client IP address
- `browser` (string) - Browser name
- `browserVersion` (string) - Browser version
- `os` (string) - Operating system
- `osVersion` (string) - OS version
- `country` (string) - Country code
- `state` (string) - State/region
- `city` (string) - City name
- `lat` (string) - Latitude
- `lon` (string) - Longitude
- `src` (string) - Traffic source
- `utm_source` (string) - UTM source
- `utm_medium` (string) - UTM medium
- `utm_campaign` (string) - UTM campaign
- `utm_term` (string) - UTM term
- `utm_content` (string) - UTM content
- `aff` (string) - Affiliate ID
- `fbclid` (string) - Facebook click ID
- `gclid` (string) - Google click ID
- `ttclid` (string) - TikTok click ID
- `clickId` (string) - General click ID
- `originalUrl` (string) - Original URL

**Headers:**
- `User-Agent` - Automatically extracted
- `Cookie` - Session cookies

**Example Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/checkout/123e4567-e89b-12d3-a456-426614174000?isMobile=false&country=BR&utm_source=google" \
  -H "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
```

**Response Example:**
```json
{
  "billing_type": "ONE_TIME",
  "is_free": false,
  "back_redirect_url": "https://example.com/back",
  "config": {
    "checkout_uuid": "checkout-123",
    "checkout_date": "2024-01-15T10:30:00Z",
    "has_discount": true,
    "logo_enabled": true,
    "logo": "https://cdn.example.com/logo.png",
    "color_primary": "#007bff",
    "color_secondary": "#6c757d",
    "credit_card_enabled": true,
    "pix_enabled": true,
    "bank_slip_enabled": true
  },
  "product": {
    "uuid": "product-456",
    "name": "Premium Course",
    "price": 99.90,
    "photo": "https://cdn.example.com/product.jpg",
    "format": "digital"
  },
  "order_bumps": [],
  "reviews": [],
  "pixels": [],
  "plans": []
}
```

## Environment Configuration

### Environment Variables

- `PORT` - Server port (default: 8080)
- `GIN_MODE` - Gin mode: `debug`, `release`, or `test`
- `ENVIRONMENT` - Alternative to GIN_MODE: `production`, `development`, `test`

### AWS Configuration

The server requires AWS credentials for DynamoDB access:

- `AWS_REGION` - AWS region (default: us-east-1)
- `AWS_ACCESS_KEY_ID` - AWS access key
- `AWS_SECRET_ACCESS_KEY` - AWS secret key

Or use AWS credential profiles/IAM roles.

## Features

### Middleware

1. **CORS Middleware**: Handles cross-origin requests
2. **Logger Middleware**: Structured request logging
3. **Recovery Middleware**: Panic recovery with proper error responses
4. **Request ID Middleware**: Adds unique request IDs for tracing

### Error Handling

All errors return a consistent JSON format:

```json
{
  "error": true,
  "message": "Error description",
  "status": 400
}
```

### Request Tracing

Each request gets a unique `X-Request-ID` header for tracking:

```
X-Request-ID: 20240115103000-ABC123
```

### Graceful Shutdown

The server handles `SIGINT` and `SIGTERM` signals gracefully:

1. Stops accepting new connections
2. Finishes processing existing requests (up to 30 seconds)
3. Closes all connections
4. Exits cleanly

## Comparison with Original Server

| Feature | Original (net/http) | New (Gin) |
|---------|-------------------|-----------|
| Framework | Standard library | Gin |
| Routing | Manual | Automatic |
| Middleware | Manual | Built-in |
| JSON handling | Manual encoding | Automatic |
| CORS | Manual headers | Middleware |
| Logging | Basic | Structured |
| Recovery | Basic | Advanced |
| Performance | Good | Better |
| Development | More code | Less code |

## Testing

### Manual Testing

```bash
# Start the server
./checkout-gin

# Test health endpoint
curl http://localhost:8080/health

# Test checkout endpoint
curl "http://localhost:8080/api/v1/checkout/123e4567-e89b-12d3-a456-426614174000?isMobile=false"
```

### Load Testing

You can use tools like `ab`, `wrk`, or `k6` for load testing:

```bash
# Using ab (Apache Bench)
ab -n 1000 -c 10 http://localhost:8080/health

# Using wrk
wrk -t12 -c400 -d30s http://localhost:8080/health
```

## Deployment

### Docker

Create a `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o checkout-gin ./cmd/gin/

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/checkout-gin .
EXPOSE 8080
CMD ["./checkout-gin"]
```

Build and run:

```bash
docker build -t checkout-gin .
docker run -p 8080:8080 checkout-gin
```

### Production

For production deployment:

1. Set `GIN_MODE=release`
2. Configure proper AWS credentials
3. Use a reverse proxy (nginx, Apache)
4. Set up monitoring and logging
5. Configure proper timeouts and limits

## Migration from Original Server

The new Gin server is fully backward compatible. You can:

1. Replace the original server gradually
2. Run both servers simultaneously on different ports
3. Use the same business logic and database

No changes needed to:
- Database schemas
- Business logic
- Use cases
- Repository implementations
- AWS configuration

## Troubleshooting

### Common Issues

1. **Port already in use**: Change the PORT environment variable
2. **AWS credentials**: Ensure proper AWS configuration
3. **DynamoDB access**: Check IAM permissions
4. **Build errors**: Run `go mod tidy` to fix dependencies

### Debug Mode

Run in debug mode for detailed logging:

```bash
GIN_MODE=debug ./checkout-gin
```

### Health Check

Always verify the server is running properly:

```bash
curl http://localhost:8080/health
```

## Future Enhancements

The Gin architecture enables easy addition of:

- Rate limiting middleware
- Authentication/authorization
- Request validation middleware
- Response caching
- Metrics collection
- API versioning
- WebSocket support
- gRPC endpoints

## Support

For issues or questions:

1. Check the server logs
2. Verify AWS configuration
3. Test with the health endpoint
4. Check environment variables
5. Review this documentation 