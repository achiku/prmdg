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


## <a name="resource-task">タスク</a>

Stability: `prototype`

タスクを表します。

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **completedAt** | *date-time* | 直近の終了時間日時 | `"2016-02-01T12:13:14Z"` |
| **createdAt** | *date-time* | 作成日時 | `"2016-02-01T12:13:14Z"` |
| **id** | *uuid* | タスクID | `"ec0a1edc-062e-11e7-8b1e-040ccee2aa06"` |
| **name** | *string* | タスク名 | `"コーヒーを買う"` |
| **startedAt** | *date-time* | 直近の開始日時 | `"2016-02-01T12:13:14Z"` |
| **status** | *string* | タスク状態<br/> **one of:**`"done"` or `"doing"` or `"todo"` | `"done"` |
| **totalDuration** | *integer* | 累積タスク対応時間(秒) | `120` |
| **[user:id](#resource-user)** | *uuid* | ユーザーID | `"ec0a1edc-062e-11e7-8b1e-040ccee2aa06"` |
| **[user:name](#resource-user)** | *string* | ユーザー名 | `"8maki"` |

### <a name="link-GET-task-/tasks/{(%23%2Fdefinitions%2Ftask%2Fdefinitions%2Fidentity)}">タスク 詳細</a>

タスクの詳細を取得します。(正常ステータスコード: `200`)

```
GET /tasks/{task_id}
```


#### Curl Example

```bash
$ curl -n https://tasky.io/v1/tasks/$TASK_ID
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "id": "ec0a1edc-062e-11e7-8b1e-040ccee2aa06",
  "name": "コーヒーを買う",
  "user": {
    "id": "ec0a1edc-062e-11e7-8b1e-040ccee2aa06",
    "name": "8maki"
  },
  "status": "done",
  "totalDuration": 120,
  "startedAt": "2016-02-01T12:13:14Z",
  "createdAt": "2016-02-01T12:13:14Z",
  "completedAt": "2016-02-01T12:13:14Z"
}
```

### <a name="link-POST-task-/tasks">タスク 作成</a>

タスクを作成します。(正常ステータスコード: `201`)

```
POST /tasks
```

#### Required Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **name** | *string* | タスク名 | `"コーヒーを買う"` |
| **userId** | *uuid* | ユーザーID | `"ec0a1edc-062e-11e7-8b1e-040ccee2aa06"` |



#### Curl Example

```bash
$ curl -n -X POST https://tasky.io/v1/tasks \
  -d '{
  "name": "コーヒーを買う",
  "userId": "ec0a1edc-062e-11e7-8b1e-040ccee2aa06"
}' \
  -H "Content-Type: application/json"
```


#### Response Example

```
HTTP/1.1 201 Created
```

```json
{
  "id": "ec0a1edc-062e-11e7-8b1e-040ccee2aa06",
  "name": "コーヒーを買う",
  "user": {
    "id": "ec0a1edc-062e-11e7-8b1e-040ccee2aa06",
    "name": "8maki"
  },
  "status": "done",
  "totalDuration": 120,
  "startedAt": "2016-02-01T12:13:14Z",
  "createdAt": "2016-02-01T12:13:14Z",
  "completedAt": "2016-02-01T12:13:14Z"
}
```

### <a name="link-GET-task-/tasks">タスク 一覧</a>

タスクの一覧を取得します。(正常ステータスコード: `200`)

```
GET /tasks
```

#### Optional Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **limit** | *integer* | 取得する要素数(デフォルト20) | `20` |
| **offset** | *integer* | 取得する要素のオフセット(0から開始) | `20` |


#### Curl Example

```bash
$ curl -n https://tasky.io/v1/tasks
 -G \
  -d limit=20 \
  -d offset=20
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
[
  {
    "id": "ec0a1edc-062e-11e7-8b1e-040ccee2aa06",
    "name": "コーヒーを買う",
    "user": {
      "id": "ec0a1edc-062e-11e7-8b1e-040ccee2aa06",
      "name": "8maki"
    },
    "status": "done",
    "totalDuration": 120,
    "startedAt": "2016-02-01T12:13:14Z",
    "createdAt": "2016-02-01T12:13:14Z",
    "completedAt": "2016-02-01T12:13:14Z"
  }
]
```


## <a name="resource-user">ユーザー</a>

Stability: `prototype`

ユーザーを表します。

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **id** | *uuid* | ユーザーID | `"ec0a1edc-062e-11e7-8b1e-040ccee2aa06"` |
| **name** | *string* | ユーザー名 | `"8maki"` |

### <a name="link-GET-user-/me">ユーザー 詳細</a>

ログイン中のユーザー情報を取得します。(正常ステータスコード: `200`)

```
GET /me
```


#### Curl Example

```bash
$ curl -n https://tasky.io/v1/me
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "id": "ec0a1edc-062e-11e7-8b1e-040ccee2aa06",
  "name": "8maki"
}
```


