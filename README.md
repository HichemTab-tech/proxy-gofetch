# Proxy GoFetch

A simple HTTP proxy server written in Go that fetches content from specified URLs and forwards it to the client. This tool is particularly useful for bypassing CORS restrictions when developing web applications.

## Features

- Proxies HTTP requests to any URL
- Preserves original content type
- Sets CORS headers to allow cross-origin requests
- Simple API with a single endpoint

## Installation

### Prerequisites

- Go 1.16 or higher

### Steps

1. Clone the repository:
   ```
   git clone https://github.com/HichemTab-tech/proxy-gofetch.git
   cd proxy-gofetch
   ```

2. Build the application:
   ```
   go build
   ```

## Usage

1. Start the server:
   ```
   ./proxy-gofetch
   ```
   The server will start on port 8080.

2. Make requests to the proxy:
   ```
   http://localhost:8080/fetch?url=https://example.com/image.jpg
   ```

### Example

To fetch an image through the proxy:
```
<img src="http://localhost:8080/fetch?url=https://example.com/image.jpg">
```

To fetch JSON data:
```javascript
fetch('http://localhost:8080/fetch?url=https://api.example.com/data.json')
  .then(response => response.json())
  .then(data => console.log(data));
```

## Error Handling

The proxy will return appropriate HTTP status codes:
- 400 Bad Request: If the URL parameter is missing or invalid
- 502 Bad Gateway: If the target URL cannot be fetched
- 500 Internal Server Error: If there are issues with processing the response

## License

This project is for personal use.