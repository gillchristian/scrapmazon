# scrapmazon

A simple scraper of Amazon Prime Movies (AMP).

```
GET /movie/amazon/:amazon_id
```
Where `amazon_id` is the AMP's movie id (the `B00K19SD8Q` in `http://www.amazon.de/gp/product/B00K19SD8Q`).

Example:

```
$ curl http://localhost:8080/movie/amazon/B00K19SD8Q
{
  "title": "Um Jeden Preis",
    "release_year": 2013,
    "actors": [
      "Dennis Quaid",
    "Zac Efron",
    "Kim Dickens"
    ],
    "poster": "https://images-eu.ssl-images-amazon.com/images/I/51VELYHd4TL._SX200_QL80_.jpg",
    "similar_ids": [
      "B00L9KET84",
      "B00JM0JXYI",
      "B00HDZMP94",
    ...
    ]
}

```


## Third party packages

- [httprouter](https://github.com/julienschmidt/httprouter)

I picked it because of it's simplicity to set route handlers and to handle route params. 
They only downside I see on this package is the handlers not being compatible with `net/http` ones.

- [scrape](https://github.com/yhat/scrape)

A really simple and small HTML scraping interface on top of [Go's HTML parsing library](https://godoc.org/golang.org/x/net/html).

## Nice to have / Improvements

- Tests.
- Use only net/http, ditch httprouter.
