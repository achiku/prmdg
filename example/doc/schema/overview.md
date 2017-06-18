## tasky.io API reference

This is psuedo Todo management service API (tasky.io) reference.


## Request Headers

It is necessary to include the below headers to request API.

- `Tasky-UUID`: device uuid
- `Tasky-App-Version`: app version
- `Tasky-App-ID`: app id

## Authorization

Whenever authorization is required, add the `Authorization` header with above headers.

```
Authorization Bearer abcdefghijklmnopqrstuvwxyzabcdefghijklmn
Tasky-UUID FCDBD8EF-62FC-4ECB-B2F5-92C9E79AC7F9
Tasky-App-Version 1.0.0
```
