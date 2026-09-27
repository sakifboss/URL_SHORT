# URL Shortener

এটি Go দিয়ে তৈরি একটি সহজ URL shortener HTTP server। এটি একটি URL-এর জন্য ৬ অক্ষরের short code তৈরি করে এবং সেই code-এ গেলে মূল URL-এ redirect করে।

## প্রয়োজনীয়তা

- Go 1.26.4 বা পরবর্তী সংস্করণ

## চালানোর নিয়ম

প্রজেক্ট ফোল্ডারে গিয়ে চালান:

```bash
go run .
```

সার্ভার চালু হলে এটি `http://localhost:8080`-এ অনুরোধ নেবে।

## API

### Health check

```http
GET /health
```

সফল হলে `OK` ফেরত দেয়।

### URL ছোট করা

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

উদাহরণ:

```bash
curl -X POST http://localhost:8080/shorten -H "Content-Type: application/json" -d '{"url":"https://example.com"}'
```

Response-এ `short_code` পাওয়া যাবে।

### Short URL ব্যবহার

```http
GET /short/{short_code}
```

এটি মূল URL-এ redirect করে। যেমন, short code `abc123` হলে `http://localhost:8080/short/abc123` খুলুন।

## সীমাবদ্ধতা

Short code এবং URL-গুলোর mapping শুধু চলমান server-এর memory-তে রাখা হয়। Server বন্ধ হলে তৈরি করা short URL-গুলো আর থাকবে না।