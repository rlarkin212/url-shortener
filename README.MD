# Url Shortener
Shortens url into 9 digit code

## Endpoints

### POST
https://go-urlshortener.herokuapp.com/

Request
```
{
    "url" : "https://github.com/rlarkin212/"
}
```
Response
```
{
  "code": "gjGVmoLZR",
  "url": "https://github.com/rlarkin212",
  "created_at": 1577739808
}
```

### GET
https://go-urlshortener.herokuapp.com/{code}
