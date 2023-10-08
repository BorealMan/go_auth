# Fiber API Authentication

A user authentication api written with fiber and sqlx. A good base to build off of for simple projects that don't require an ORM.

## Features

- Includes JWT Middleware
- Database Schema
- Automatically Seed User Roles
- Create Users, Login, Use Middleware On Protected Routes To Automatically Add User ID and Role To Request Header.

### To Use:

You must have a mysql or mariadb database to use this service.

Add your database settings to config

#### Running Service

```
go run main.go
```

#### Create an executable

```
go build main.go
```

### Optionally Run With Docker Compose

#### Build the go-user-api image

```
docker build --no-cache -t go-user-api .
```

#### Start service

```
docker compose up
```

#### Stop service

```
docker compose down
```
