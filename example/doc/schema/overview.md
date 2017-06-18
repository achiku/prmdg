## ヘッダ

APIを利用するアプリのバージョンやID、端末を一意に識別するIDをヘッダに含める必要があります。

- `Tasky-UUID`: デバイスを識別する一意なID
- `Tasky-App-Version`: 利用しているアプリのバージョン
- `Tasky-App-ID`: 利用しているアプリのID

## 認証

認証が必要なリクエストの場合は通常のヘッダ項目に加えて、Authorizationをヘッダに含める必要があります。

```
Authorization Bearer abcdefghijklmnopqrstuvwxyzabcdefghijklmn
Tasky-UUID FCDBD8EF-62FC-4ECB-B2F5-92C9E79AC7F9
Tasky-App-Version 1.0.0
```
