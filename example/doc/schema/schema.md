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

## The table of contents

- <a href="#resource-error">Error</a>
- <a href="#resource-task">Task</a>
  - <a href="#link-GET-task-/tasks/{(%23%2Fdefinitions%2Ftask%2Fdefinitions%2Fidentity)}">GET /tasks/{task_id}</a>
  - <a href="#link-POST-task-/tasks">POST /tasks</a>
  - <a href="#link-GET-task-/tasks">GET /tasks</a>
- <a href="#resource-user">User</a>
  - <a href="#link-GET-user-/me">GET /me</a>

## <a name="resource-error">Error</a>

Stability: `prototype`

This resource represents API error

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **code** | *string* | error code<br/> **one of:**`"invalid_params"` or `"invalid_request"` or `"unauthorized"` or `"unsupported_client_version"` | `"invalid_params"` |
| **detail** | *string* | error detail | `"invalid param"` |
| **errorFields/message** | *string* | error message for invalid param field | `"invalid status"` |
| **errorFields/name** | *string* | param field name | `"status"` |


## <a name="resource-task">Task</a>

Stability: `prototype`

This resource represents task

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **completedAt** | *date-time* | time completed a task | `"2016-02-01T12:13:14Z"` |
| **createdAt** | *date-time* | time created a task | `"2016-02-01T12:13:14Z"` |
| **id** | *uuid* | task id | `"ec0a1edc-062e-11e7-8b1e-040ccee2aa06"` |
| **spent** | *integer* | time spent doing task in minutes | `12` |
| **startedAt** | *date-time* | time started a task | `"2016-02-01T12:13:14Z"` |
| **status** | *string* | task status<br/> **one of:**`"done"` or `"doing"` or `"stopped"` | `"done"` |
| **tags** | *array* | tags | `["study"]` |
| **title** | *string* | task title | `"Buy coffee"` |
| **[user:id](#resource-user)** | *uuid* | user id | `"ec0a1edc-062e-11e7-8b1e-040ccee2aa06"` |
| **[user:name](#resource-user)** | *string* | user name | `"8maki"` |

### <a name="link-GET-task-/tasks/{(%23%2Fdefinitions%2Ftask%2Fdefinitions%2Fidentity)}">Task detail</a>

Get task detail

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
  "title": "Buy coffee",
  "user": {
    "id": "ec0a1edc-062e-11e7-8b1e-040ccee2aa06",
    "name": "8maki"
  },
  "status": "done",
  "spent": 12,
  "startedAt": "2016-02-01T12:13:14Z",
  "createdAt": "2016-02-01T12:13:14Z",
  "completedAt": "2016-02-01T12:13:14Z",
  "tags": [
    "study"
  ]
}
```

### <a name="link-POST-task-/tasks">Task create</a>

Create task

```
POST /tasks
```

#### Required Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **title** | *string* | task title | `"Buy coffee"` |


#### Optional Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **tags** | *array* | tags | `["study"]` |


#### Curl Example

```bash
$ curl -n -X POST https://tasky.io/v1/tasks \
  -d '{
  "title": "Buy coffee",
  "tags": [
    "study"
  ]
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
  "title": "Buy coffee",
  "user": {
    "id": "ec0a1edc-062e-11e7-8b1e-040ccee2aa06",
    "name": "8maki"
  },
  "status": "done",
  "spent": 12,
  "startedAt": "2016-02-01T12:13:14Z",
  "createdAt": "2016-02-01T12:13:14Z",
  "completedAt": "2016-02-01T12:13:14Z",
  "tags": [
    "study"
  ]
}
```

### <a name="link-GET-task-/tasks">Task list</a>

Get task list

```
GET /tasks
```

#### Optional Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **limit** | *integer* | limit | `20` |
| **offset** | *integer* | offset | `20` |


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
    "title": "Buy coffee",
    "user": {
      "id": "ec0a1edc-062e-11e7-8b1e-040ccee2aa06",
      "name": "8maki"
    },
    "status": "done",
    "spent": 12,
    "startedAt": "2016-02-01T12:13:14Z",
    "createdAt": "2016-02-01T12:13:14Z",
    "completedAt": "2016-02-01T12:13:14Z",
    "tags": [
      "study"
    ]
  }
]
```


## <a name="resource-user">User</a>

Stability: `prototype`

This resource represents user

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **id** | *uuid* | user id | `"ec0a1edc-062e-11e7-8b1e-040ccee2aa06"` |
| **name** | *string* | user name | `"8maki"` |

### <a name="link-GET-user-/me">User detail</a>

Get authenticated user detail

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


