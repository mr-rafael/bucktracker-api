# Finance Calculator

**A REST API with useful Savings and Loans calculations.**


## Technologies Used

- Go
- PostgreSQL

## Getting Started

### Prerequisites

Before running the project, make sure you have the following installed:

- [Go](https://go.dev/doc/install) 
- [PostgreSQL](https://www.postgresql.org/download/)
- [Goose](https://pkg.go.dev/github.com/pressly/goose/v3#section-readme) (Database migration tool)

---

### Installation

Clone the repository:

```bash
git clone git@github.com:Mr-Rafael/finance-calculator.git
cd finance-calculator
```

---

### Configure the Database

Create a PostgreSQL database:

```sql
CREATE DATABASE your_database_name;
```

Set the database connection environment variables as needed by the project.

Example:

```bash
ALLOWED_ORIGIN=http://localhost:5173
POSTGRES_CONNECTION_STRING=postgres://<username>:<password>@localhost:5432/finance_calculator?sslmode=disable
ACCESS_SECRET=DEVENVIRONMENTSECRET
REFRESH_SECRET=DEVENVREFRESHSECRET
ENV=develop
```

- **ALLOWED_ORIGIN**: Used for CORS. Necessary when you're running both the server and a client on the same computer.
- **POSTGRES_CONNECTION_STRING**: The user, password and address of the PostgreSQL database you're running.
- **ACCESS_SECRET**: The secret that will be used to sign Access Tokens.
- **REFRESH_SECRET**: The secret that will be used to sign Refresh Tokens.
- **ENV**: If set to "production", the Refresh Token cookie will be set to Secure, and only be sent via HTTPS.

---

### Installing Goose

If Goose is not installed:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

---

### Run Database Migrations

Apply the migrations using Goose:

```bash
goose postgres "host=localhost port=5432 user=postgres password=postgres dbname=your_database_name sslmode=disable" up
```

---

### Install Dependencies

```bash
go mod tidy
```

---

### Run the Project

```bash
go run ./cmd/server/main.go
```

---

### Notes

- Make sure PostgreSQL is running before starting the application.
- Ensure the database credentials match your local setup.

## API Endpoints

| route               | description                                          
|----------------------|-----------------------------------------------------
| <kbd>GET /api/healthz</kbd>     | Check if server is running.
| <kbd>POST /app/users/create</kbd>     | Create a new user.
| <kbd>POST /app/login</kbd>     | Login user.
| <kbd>POST /app/refresh</kbd>     | Refresh the access token.
| <kbd>POST /app/savings/calculate</kbd>     | Generate a Savings Plan without saving.
| <kbd>POST /app/loans/calculate</kbd>     | Calculate a Loan Payment Plan without saving.
| <kbd>POST /app/savings/save</kbd>     | Calculate and save a Savings Plan.
| <kbd>POST /app/loans/save</kbd>     | Calculate and save a Loan Payment Plan.
| <kbd>GET /app/savings/list</kbd>     | List a User's saved Savings Plans. 
| <kbd>GET /app/loans/list</kbd>     | List a User's saved Loan Payment Plans.
| <kbd>GET /app/savings/{id}</kbd>     | Get a previously saved Savings Plan.
| <kbd>GET /app/loans/{id}</kbd>     | Get a previously saved Loan Payment Plan.
| <kbd>PATCH /app/savings/{id}</kbd>     | Update and recalculate a Savings Plan.
| <kbd>PATCH /app/loans/{id}</kbd>     | Update and recalculate a Loan Payment Plan.
| <kbd>DELETE /app/savings/{id}</kbd>     | Delete a Savings Plan.
| <kbd>DELETE /app/loans/{id}</kbd>     | Delete a Loan Payment Plan.

## `GET /api/healthz`

Short description of what the endpoint does.

---

### Request

### URL

```http
METHOD /endpoint
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | No | Bearer token if authentication is required |
| Content-Type | Yes | Usually `application/json` |

### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| id | string | Yes | Resource identifier |

### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| limit | integer | No | Number of items to return |

### Request Body

```json
{
  "field1": "value",
  "field2": 123
}
```

### Request Fields

| Field | Type | Required | Description |
|---|---|---|---|
| field1 | string | Yes | Example string field |
| field2 | integer | No | Example numeric field |

---


### Success Response

**Status Code:** `200 OK`

```json
{
  "id": "123",
  "field1": "value",
  "created_at": "2026-05-28T12:00:00Z"
}
```

### Response Fields

| Field | Type | Description |
|---|---|---|
| id | string | Resource ID |
| field1 | string | Example field |
| created_at | string | ISO 8601 timestamp |

---

### Error Responses

### `400 Bad Request`

```json
{
  "error": "invalid request body"
}
```

### `404 Not Found`

```json
{
  "error": "resource not found"
}
```

### `500 Internal Server Error`

```json
{
  "error": "internal server error"
}
```

---


| route               | description                                          
|----------------------|-----------------------------------------------------
| <kbd>GET /api/healthz</kbd>     | Check if server is running.
| <kbd>POST /app/users/create</kbd>     | Create a new user.
| <kbd>POST /app/login</kbd>     | Login user.
| <kbd>POST /app/refresh</kbd>     | Refresh the access token.
| <kbd>POST /app/savings/calculate</kbd>     | Generate a Savings Plan without saving.
| <kbd>POST /app/loans/calculate</kbd>     | Calculate a Loan Payment Plan without saving.
| <kbd>POST /app/savings/save</kbd>     | Calculate and save a Savings Plan.
| <kbd>POST /app/loans/save</kbd>     | Calculate and save a Loan Payment Plan.
| <kbd>GET /app/savings/list</kbd>     | List a User's saved Savings Plans. 
| <kbd>GET /app/loans/list</kbd>     | List a User's saved Loan Payment Plans.
| <kbd>GET /app/savings/{id}</kbd>     | Get a previously saved Savings Plan.
| <kbd>GET /app/loans/{id}</kbd>     | Get a previously saved Loan Payment Plan.
| <kbd>PATCH /app/savings/{id}</kbd>     | Update and recalculate a Savings Plan.
| <kbd>PATCH /app/loans/{id}</kbd>     | Update and recalculate a Loan Payment Plan.
| <kbd>DELETE /app/savings/{id}</kbd>     | Delete a Savings Plan.
| <kbd>DELETE /app/loans/{id}</kbd>     | Delete a Loan Payment Plan.

## `GET /api/healthz`

Check if the server is running.

### URL

```http
GET /api/healthz
```
---


### Success Response

**Status Code:** `200 OK`

```text
OK
```
---

## Collaborators</h2>

<table>
  <tr>
    <td align="center">
      <a href="#">
        <img src="https://avatars.githubusercontent.com/u/35672719?s=48&v=4" width="100px;" alt="Rafael Mazariegos picture"/><br>
        <sub>
          <b>Rafael Mazariegos</b>
        </sub>
      </a>
    </td>
  </tr>
</table>