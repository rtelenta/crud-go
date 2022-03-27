# users crud
## available methods

```
GET `/` - healthcheck
```

```
GET `/api/users` - list all users
```

```
POST `/api/users/create` - create user

payload
{
    "name": "test",
    "email": "name@test.com"
}
```

```
GET `/api/users/{{USER_ID}}` - get user
```

```
DELETE `/api/users/{{USER_ID}}` - delete user
```

```
PATCH `/api/users/{{USER_ID}}` - update user

payload
{
    "name": "test",
    "email": "name@test.com"
}
```




# Notes
idk why `go get github.com/codegangsta/gin` dont work use `go install -v github.com/codegangsta/gin` instead