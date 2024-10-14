# URL Shortener

A simple URL shortener service built with Go, Gin, MongoDB, and rate limiting. This service allows users to shorten URLs, track the number of clicks, set custom expiration times, and limit requests based on IP addresses.

## Features

- URL Shortening: Converts long URLs into shorter, more manageable URLs.
- Custom Expiration: Users can specify expiration times for their URLs in days, hours, and minutes.
- Rate Limiting: Limits the number of requests from a single IP to prevent abuse. Individual users are allowed upto maximum of 100 requests per hour.
- Analytics: Tracks how many times each shortened URL has been clicked.
  MongoDB Integration: Stores shortened URLs, expiration times, and hit counts in MongoDB.

### Usage Example

Copy paste this code in your terminal

```
curl -X POST http://localhost:8080/shorten \
-H "Content-Type: application/json" \
-d '{
  "lurl": "https://example.com",
  "days": 1,
  "hours": 2
}'
```

You should get a sample shortened url as `http://localhost:8000/abc123`
