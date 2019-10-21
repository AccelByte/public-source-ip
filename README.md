[![Build Status](https://travis-ci.com/AccelByte/public-source-ip.svg?branch=master)](https://travis-ci.com/AccelByte/public-source-ip)

# Public Source IP

This package determines the public IP used to connect to the first ingress load balancer, by the X-Forwarded-For header.

## Usage
### Importing
```go
import publicsourceip "github.com/AccelByte/public-source-ip"
```

### Call PublicIP
```go
publicIP := publicsourceip.PublicIP(&http.Request{Header: request.Request.Header})
```
