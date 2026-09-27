# Go URL Shortener

A small URL shortener built with Go's standard library. It generates a six-character short code for a URL and redirects visitors to the original address.

## Requirements

- Go 1.26.4 or later

## Run the server

From the project directory, run:

```bash
go run .
```

The server listens at <http://localhost:8080>.

## Run with Docker Compose

Make sure Docker Desktop is running, then start the app from the project directory:

```bash
docker compose up --build
```

Open <http://localhost:8080> in your browser. Press `Ctrl+C` in the terminal to stop the app. To remove the Compose container, run:

```bash
docker compose down
```

## Deploy to Render

This project can run as a Render web service. In Render, create a new **Web Service** and connect this GitHub repository. Use these settings:

| Setting | Value |
| --- | --- |
| Runtime | Go |
| Build command | `go build -o app .` |
| Start command | `./app` |
| Health check path | `/health` |

Render supplies the listening port through the `PORT` environment variable. The server reads this value automatically and binds to all network interfaces.

The app stores short URLs and rate-limit counts in memory. They are lost whenever the service restarts or redeploys. Keep the service to one instance so all requests share the same in-memory data.

## API

### Health check

```http
GET /health
```

Returns `OK` when the server is running.

### Shorten a URL

```http
POST /shorten
Content-Type: application/json
```

Request body:

```json
{
  "url": "https://example.com"
}
```

Example using `curl`:

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}'
```

The response includes a `short_code`, which can be used with the redirect endpoint.

### Follow a short URL

```http
GET /short/{short_code}
```

For example, `GET http://localhost:8080/short/abc123` redirects to the original URL associated with `abc123`.

### Rate limit

The `POST /shorten` endpoint allows up to five requests per IP address in each one-minute window. Further requests during that window receive `429 Too Many Requests` with this response:

```json
{"error":"Rate limit exceeded. Try again later."}
```

## Storage

Short-code mappings and rate-limit counts are kept in memory. They are cleared when the server stops.
