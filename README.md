# Receipt Processor in Go

## Intro

Receipt Processor project that is written fully in Go. This is my very first introductory to the language, and enjoyed writing in it very much. I feel as if there were better ways to write some of the logic, but the main goal for me was to learn the langauge as simply and concise as I can. 

## API Specification

### Process Receipt

- Path: `/Receipts/process`
- Method: `POST`
- Payload: Receipt JSON
- Response: JSON containing an id for the receipt

Passing in JSON receipt into `POST` method:

```json
{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}
```

Response:

**Note: the `id` is always going to be randomly generated.**

```json
{
  "id": "e4b135a1-304f-4518-8562-9480dbc796ce"
}
```

### Get points

- Path: `/Receipts/{id}/points`
- Method: `GET`
- Response: A JSON object containing the receipt by the ID and returns an object specifying the points.

```json
{
  "points": 109
}
```
