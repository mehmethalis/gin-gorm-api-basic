## Installation

Run the following command to install packages and dependencies

```sh
go mod tidy
```

Create ``.env`` file in root directory:

```sh
DB_HOST=localhost
DB_USER=user
DB_PASSWORD=password
DB_NAME=test
DB_PORT=5432
JWT_SECRET=secret
```

## Endpoints
<br>

#### Auth

| Route | HTTP Verb	 |Body	 | Description	 |
| --- | --- | --- | --- |
| /api/auth/login | `POST` | {'email': 'email@email.com', 'password':'*****'} | Validate credentials |
| /api/auht/register | `POST` | {'name':'foo', 'email': 'email@email.com', 'password':'*****'} | Create a new user |

#### User

| Route | HTTP Verb	 |Body	 | Header | Description	 | 
| --- | --- | --- | --- | --- |
| /api/user/profile | `PUT` | {'name': 'john' , email': 'email@email.com', 'password':'*****'} | { 'Authorization' : 'Token'}| Validate credentials |
| /api/user/profile | `GET` | {'name':'foo', 'email': 'email@email.com', 'password':'*****'} | | Create a new user |