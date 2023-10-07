# Fiber API Authentication

A user authentication api written with fiber and sqlx. A good base to build off of for simple projects that don't require an ORM.

## Features

- Includes JWT Middleware
- Database Schema
- Automatically Seed User Roles
- Create Users, Login, Use Middleware On Protected Routes To Automatically Add User ID and Role To Request Header.

### To Use:

#### Add database settings to config

```
go run main.go
```

##### It's that easy!
