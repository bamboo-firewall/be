# Bamboo Firewall API Server

API for gathering bamboo firewall information:

* Host End Point
* Global Network Set
* Global Network Policy (Policy)

## How to run? 

Using https://github.com/spf13/viper to read configs. Priority:

- env
- json config
- key/value store
- default

Container way: modified file `deployment/.env`

"Go run" way: modified file `./config.json`

## Requirement

* mongodb

## Public API

1. Ping

```bash
curl --location 'localhost:8080/api/ping'
```

2. Register

```bash
curl -L 'localhost:8080/api/signup' \
-H 'Content-Type: application/json' \
--data-raw '{
    "name": "buycoffee+3c1c0z2b",
    "email": "buycoffee+3c1c0z2b@example.com",
    "password": "immorally8578"
}'
```

3. Login

```bash
curl -L 'localhost:8080/api/login' \
-H 'Content-Type: application/json' \
--data-raw '{
    "email": "buycoffee+3c1c0z2b@example.com",
    "password": "immorally8578"
}'
```

## Protected API

1. Fetch HEP

```bash
curl -L -X POST 'localhost:8080/api/v1/hep/fetch' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiQ2FvIFh1YW4gQW5oIiwiaWQiOiI2NDk5NGUzMDc0YWRhZWFiZDY0MWNmMDIiLCJleHAiOjE2ODc4NTc4NzF9.dOzvXgUM12epaJDXZ4jbF0KjZddh2B1UHr_MrbXIubk'
```

2. Fetch GNS - Global Network Set

```bash
curl -L -X POST 'localhost:8080/api/v1/gns/fetch' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiQ2FvIFh1YW4gQW5oIiwiaWQiOiI2NDk5NGUzMDc0YWRhZWFiZDY0MWNmMDIiLCJleHAiOjE2ODc4NTc4NzF9.dOzvXgUM12epaJDXZ4jbF0KjZddh2B1UHr_MrbXIubk'
```

3. Fetch Policy

```bash
curl -L -X POST 'localhost:8080/api/v1/policy/fetch' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiQ2FvIFh1YW4gQW5oIiwiaWQiOiI2NDk5NGUzMDc0YWRhZWFiZDY0MWNmMDIiLCJleHAiOjE2ODc4NTc4NzF9.dOzvXgUM12epaJDXZ4jbF0KjZddh2B1UHr_MrbXIubk'
```

4. Fetch options

Option - object using to searching
Option.key - field name
Option.value - value of field

Filter: value to filters

```bash
curl -L 'localhost:9091/api/v1/options/fetch' \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiQWRtaW4iLCJpZCI6IjY0YjYwMmQ0ODk2YTJhODc3ZTI3NGE0YSIsImV4cCI6MTY4OTkzMTU5OH0.iAyD2nsRCmfkaGgyQBA1LnsG7ly1wM59HC2e4Lm5F-U' \
-d '{
    "type": "hostendpoints",
    "label": "zone",
    "filter": [
        {
            "key": "namespace",
            "value": "non-production"
        },
        {
            "key": "role",
            "value": "lb"
        }
    ]
}'
```

Valid input:

| Type                  | Label                              |
| --------------------- | ---------------------------------- |
| hostendpoints         | ip, namespace, project, role, zone |
| globalnetworksets     | name, zone                         |
| globalnetworkpolicies | name                               |

5. Search HEP - Host End Point

```bash
curl -L 'localhost:9091/api/v1/hep/search' \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiQWRtaW4iLCJpZCI6IjY0YWZkYWNiNDhiNTBiNDhiMTk2MTFiOCIsImV4cCI6MTY4OTU4OTMwNH0.BqPgjbe644GST2uNWDCTgOC4LsuFB6B-f8SOoguZVEc' \
-d '{
    "options": [
        {
            "key": "namespace",
            "value": "non-production"
        },
        {
            "key": "role",
            "value": "lb"
        }
    ]
}'
```

6. Search Policy

```bash
curl -L 'localhost:9091/api/v1/policy/search' \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiQWRtaW4iLCJpZCI6IjY0YWZkYWNiNDhiNTBiNDhiMTk2MTFiOCIsImV4cCI6MTY4OTU5NjY2OX0.LcfpvySJHtRCuG4VvY4clMjmOAOYV7XanHVKVrdWR1E' \
-d '{
    "options": [
        {
            "key": "name",
            "value": "gchat-gdrives-app"
        }
    ]
}'
```

7. Search GNS - Global Network Set

```bash
curl -L 'localhost:9091/api/v1/gns/search' \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiQWRtaW4iLCJpZCI6IjY0YWZkYWNiNDhiNTBiNDhiMTk2MTFiOCIsImV4cCI6MTY4OTU5NjY2OX0.LcfpvySJHtRCuG4VvY4clMjmOAOYV7XanHVKVrdWR1E' \
-d '{
    "options": [
        {
            "key": "zone",
            "value": "gray"
        },
        {
            "key": "name",
            "value": "vpn-devops"
        }
    ]
}'
```

8. Statistic API - Summary

```bash
curl -L -X POST 'localhost:9091/api/v1/statistic/summary' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiQWRtaW4iLCJpZCI6IjY0YWZkYWNiNDhiNTBiNDhiMTk2MTFiOCIsImV4cCI6MTY4OTYwOTU3OX0.tObMwnGLrzeVehA5EvpXEloyRO63NManFQ6fUkGDleY'
```

Sample response

```json
{
    "summary": {
        "total_global_network_set": 122,
        "total_policy": 431,
        "total_host_endpoint": 1397,
        "total_user": 2
    }
}
```

9. Statistic API - Project Summary

```bash
curl -L -X POST 'localhost:9091/api/v1/statistic/project-summary' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiQWRtaW4iLCJpZCI6IjY0YWZkYWNiNDhiNTBiNDhiMTk2MTFiOCIsImV4cCI6MTY4OTYwOTU3OX0.tObMwnGLrzeVehA5EvpXEloyRO63NManFQ6fUkGDleY'
```

Sample response

```
{
    "project_summary": [
        {
            "project_name": "deceiver5553",
            "total": 5
        },
        {
            "project_name": "disgrace3860",
            "total": 2
        },
        {
            "project_name": "taunt7475",
            "total": 201
        }
    ]
}
```

10. Get Profile

```bash
curl -L -X POST 'localhost:9091/api/v1/profile' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiTmd1eWVuIFF1YW5nIFRoaW5oIiwiaWQiOiI2NGI3NjM4M2Q4Nzg3NDhjZmE1ZTVhZTciLCJleHAiOjE2ODk3NDczNjN9.3rMgxjQH0yoZfoZJvZi5yqJv4neoNS-0evmR6bZyInk'
```

Sample response

```json
{
    "user_id:": "64b60588eed1f4a039c0fc9f",
    "name": "buycoffee+763eub9b",
    "email": "buycoffee+763eub9b@example.com",
    "role": "devops"
}
```

11. Update profile

```bash
curl -L 'localhost:9091/api/v1/profile/update' \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiQWRtaW4iLCJpZCI6IjY0YjYwMmQ0ODk2YTJhODc3ZTI3NGE0YSIsImV4cCI6MTY4OTkzMTU5OH0.iAyD2nsRCmfkaGgyQBA1LnsG7ly1wM59HC2e4Lm5F-U' \
-d '{
    "name": "buycoffee+763eub9b"
}'
```

**// TODO: revoke token when password is updated**

for now,  frontend will erase token in local storage and make "fake" revoke token :(

Sample response

```json
{
    "message": "Update profile successfully",
    "data": {
        "name": "fragrant1852",
        "email": "admin@example.com",
        "role": "admin"
    }
}
```



## Admin API - User management

Required: Role = admin

1. Fetch User

```bash
curl -L -X POST 'localhost:9091/api/v1/admin/user/fetch' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiQWRtaW4iLCJpZCI6IjY0YjYwMmQ0ODk2YTJhODc3ZTI3NGE0YSIsImV4cCI6MTY4OTY1NzY2NX0.dGji61f160f_eJUSRd7tLJU_wEJ4KDppNx5NZgz0nfE'
```

Sample response:

```json
[
    {
        "ID": "64b602d4896a2a877e274a4a",
        "name": "Admin",
        "email": "admin@example.com",
        "role": "admin"
    },
    {
        "ID": "64b604733b869b6e81dbee8d",
        "name": "buycoffee+3c1c0z2b",
        "email": "buycoffee+3c1c0z2b@example.com",
        "role": "devops"
    }
]
```



2. Create User

```bash
curl -L 'localhost:9091/api/v1/admin/user/create' \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiQWRtaW4iLCJpZCI6IjY0YjYwMmQ0ODk2YTJhODc3ZTI3NGE0YSIsImV4cCI6MTY4OTY1NzY2NX0.dGji61f160f_eJUSRd7tLJU_wEJ4KDppNx5NZgz0nfE' \
--data-raw '{
    "name": "buycoffee+0hb01cy5",
    "email": "buycoffee+0hb01cy5@example.com",
    "password": "6Ga9yGqMmgLdSRrs",
    "role": "admin"
}'
```



3. Delete User

```bash
curl -L 'localhost:9091/api/v1/admin/user/delete' \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiVHJpbmggRGluaCBCaWVuIiwiaWQiOiI2NGI2MDU4OGVlZDFmNGEwMzljMGZjOWYiLCJleHAiOjE2OTAwMzc1OTh9.rrgdQNJ8CCJey5xEHhmi5Zv9s3WRpIkr2wuV9qhkEwc' \
-d '{
    "id": "64bbd1ed0b175d0395806bd9"
}'
```

**// TODO: revoke token when password is updated**

Response:

```json
{
    "message": "User deleted successfully!"
}
```

4. Update User

```bash
curl -L 'localhost:9091/api/v1/admin/user/update' \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiQWRtaW4iLCJpZCI6IjY0YjYwMmQ0ODk2YTJhODc3ZTI3NGE0YSIsImV4cCI6MTY4OTkzMTU5OH0.iAyD2nsRCmfkaGgyQBA1LnsG7ly1wM59HC2e4Lm5F-U' \
-d '{
    "id": "64b60588eed1f4a039c0fc9f",
    "name": "buycoffee+9wvp2i6v",
    "password": "dedicate4673",
    "role": "admin"
}'
```

Required: ID

Can't not update email

**// TODO: revoke token when password is updated**
