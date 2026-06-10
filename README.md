# Finance Calculator

**A REST API with useful Savings and Loans calculations.**

Currently, it can help you:

1. Calculate your future savings on a savings account, depending on how much you wish to deposit every month, what interest rate it has, the term, and other parameters.
2. Calculate how long it will take you to pay a loan, what your cost of credit will be, depending on the interest rate, your monthly payments, and other parameters. 

## Motivation

I often find myself making finance-related questions, like:

- If I deposit x amount each month on my savings account, how much will I have after y years at this rate?
- How much are my savings really growing, considering inflation?
- If I take an x-year loan instead of a y-year one, how much more money will I end up paying?
- How much earlier will my loan be paid if I make a paydown now?

I used to make these calculations on spreadsheets, but they become harder to manage with more complex calculations. Instead, I decided to build this RESTful API to implement these calculations more easily.

## Technologies Used

- Go
- PostgreSQL

# Quick Start

## Prerequisites

To run the project, make sure you have [Docker](https://www.docker.com/products/docker-desktop/) installed and running.

## Installation

Clone the repository:

```bash
git clone git@github.com:Mr-Rafael/finance-calculator.git
cd finance-calculator
```

## Running the Project

After cloning the project, enter the project directory and run:

```
docker compose up --build -d
```

### What the docker compose command does (if interested)

The previous command will automatically:

- Start a container running PostgreSQL.
- Create the volume `finance-calculator_postgres_data` and attach it to the container.
- Run the database migrations located in `/internal/db/migrations` using Goose Migrations.
- Start a container running the Finance Calculator REST API server.

## Usage example

### Calculate a Savings Plan

You can calculate a Savings Plan by using the endpoint:

`http://localhost:8080/app/savings/calculate`

With a JSON body in this format:

```json
{
	"startingCapital": 10000000,
	"yearlyInterestRate": "4.75",
	"monthlyContribution": 10000,
	"durationYears": 1,
	"startDate": "1970-01-01"
}
```

You should receive a response in this format:

```json
{
    "monthlyInterestRate": "0.4074123784",
    "totalEarnings": 502726,
    "totalDeposited": 10120000,
    "rateOfReturn": "4.97",
    "inflationAdjustedROR": "4.97",
    "plan": [array of monthly statuses here]
}
```

### Calculate a Loan Payment Plan

You can calculate a Savings Plan by using the endpoint:

`http://localhost:8080/app/loans/calculate`

With a JSON body in this format:

```json
{
	"startingPrincipal": 10000000,
	"yearlyInterestRate": "5",
	"monthlyPayment": 750000,
	"escrowPayment": 10000,
	"startDate": "1970-01-01"
}
```

You should receive a response in this format:

```json
{
    "durationMonths": 14,
    "totalExpenditure": 454085,
    "totalPaid": 10454085,
    "costOfCreditPercent": "4.54",
    "plan": [array of monthly statuses here]
}
```

### What else can you do?

You can currently:
- Create a user.
- Save your Loan Payment Plans and Savings Plans for future reference.
- Get, Update or Delete the Plans you previously saved.

Check the API Endpoints section for advanced usage information.

# API Endpoints

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
## `POST /app/users/create`

Create a new user.

---

### URL

```http
POST /app/users/create
```

### Headers

| Header | Required | Description |
|---|---|---|
| Content-Type | Yes | `application/json` |

### Request Body

```json
{
	"email":"user@mail.com",
    "password":"password",
    "username":"User Name"
}
```

### Request Fields


| Parameter | Type | Required | Description |
|---|---|---|---|
| email | string | Yes | Email of the user to create. Must be unique to the user. |
| password | string | Yes | Password for the user. |
| username | string | Yes | Display name for the user. |

---

### Response

### Success Response

**Status Code:** `201 Created`

```json
{
    "ID": "fa3ad421-c507-4d3c-bd77-e7918a67aaae",
    "Email": "user@mail.com",
    "Username": "User Name",
    "CreatedAt": "1970-01-01T21:39:22.159435"
}
```

### Response Fields

| Field | Type | Description |
|---|---|---|
| ID | string | UUID of the created user. |
| Email | string | User's email. |
| Username | string | User's Display Name. |
| CreatedAt | string | ISO 8601 timestamp |

---

### Error Responses

### `400 Bad Request`

```json
{
  "error": "invalid request body"
}
```
### `500 Internal Server Error`

```json
{
  "error": "internal server error"
}
```
## `POST app/login`

Login to receive an access and refresh token. Access tokens are necessary to access some endpoints.

---

### URL

```http
POST /app/login
```

### Headers

| Header | Required | Description |
|---|---|---|
| Content-Type | Yes | `application/json` |

### Request Body

```json
{
  "email":"user@mail.com",
  "password":"password"
}
```

### Request Fields

| Field | Type | Required | Description |
|---|---|---|---|
| email | string | Yes | User's email. |
| password | string | Yes | User's password. |

---

### Response

### Success Response

**Status Code:** `200 OK`

```json
{
    "id": "70fe1047-9c22-42e8-baac-64772c5c475a",
    "email": "user@mail.com",
    "username": "Test",
    "access_token": "<access token here>"
}
```

### Response Fields

| Field | Type | Description |
|---|---|---|
| id | string | User's UUID (version 4) |
| email | string | User's email |
| username | string | User's display name |
| access_token | string | JWT token to access authenticated endpoints |

---

### Error Responses

### `401 Unauthorized`

```json
{
  "error": "invalid request body"
}
```
---
## `POST /app/refresh`

Obtain a new access token, if the old one has expired.

---

### URL

```http
POST /app/refresh
```

### Headers

| Header | Required | Description |
|---|---|---|
| Cookie | Yes | The cookie should contain the refresh token obtained from the Login endpoint. |

### Request Body

No body. All required information is on the cookie.

---

### Response

### Success Response

**Status Code:** `200 OK`

```json
{
    "access_token": "[access token here]"
}
```

### Response Fields

| Field | Type | Description |
|---|---|---|
| access_token | string | A new, valid access token to access authenticated endpoints. |

---

### Error Responses

### `401 Unauthorized`

```json
{
    "error": "missing refresh token"
}
```
## `POST /app/savings/calculate`

Calculate a Savings Plan without saving it.

---

### URL

```http
POST /app/savings/calculate
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer token `Bearer [access token here]` |
| Content-Type | Yes | Usually `application/json` |

### Request Body

```json
{
	"startingCapital": 700000,
	"yearlyInterestRate": "4.75",
    "interestRateType": "APR",
	"monthlyContribution": 15000,
	"durationYears": 1,
    "taxRate": "5",
    "yearlyInflationRate": "6",
	"startDate": "1970-01-01"
}
```

### Request Fields

| Field | Type | Required | Description |
|---|---|---|---|
| startingCapital | integer | Yes | How much money was deposited at the start of the term, in cents. For example, startingCapital = 100 would mean $1. |
| yearlyInterestRate | string | Yes | The yearly interest rate for the savings plan. Send as a percent. For example, "6.25" would be a 6.25% interest rate. |
| startDate | integer | Yes | The start date of the savings plan. |
| durationYears | integer | Yes | The term you want to calculate in years. 1 means "calculate the savings plan for 1 year".  ISO 8601|
| interestRateType | string | No | Send "APR" or "APY", depending on the type of interest rate. If empty, it defaults to APY. |
| monthlyContribution | integer | No | The monthly deposits that will be made (if any). Defaults to 0 if not in the request. The amount is in cents (e.g. 15000 means $150) |
| taxRate | string | No | The tax rate paid on returns. Send as a percent. For example, "5" means a 5% tax rate. Defaults to 0% if not in the request. |
| yearlyInflationRate | string | No | The yearly inflation rate, used for rate of return calculations. Send as a percent. For example, "6" means a 6% yearly inflation rate. Defaults to 0% if not in the request..  |

---

### Response

### Success Response

**Status Code:** `200 OK`

```json
{
    "monthlyInterestRate": "0.4074123784",
    "totalEarnings": 500000,
    "totalDeposited": 10000000,
    "rateOfReturn": "5",
    "inflationAdjustedROR": "5",
    "plan": [
        {
            "date": "2026-03-01T00:00:00Z",
            "interest": 40741,
            "tax": 0,
            "contribution": 0,
            "increase": 40741,
            "capital": 10040741
        },
        {
            "date": "2026-03-01T00:00:00Z",
            "interest": 40741,
            "tax": 0,
            "contribution": 100,
            "increase": 40741,
            "capital": 10040741
        }
    ]
}
```

### Response Fields

| Field | Type | Description |
|---|---|---|
| monthlyInterestRate | string | The monthly interest rate that was used on the calculations. |
| totalEarnings | integer | The total amount in cents that you earned in interest. |
| totalDeposited | integer | The total amount deposited in cents. Includes the initial deposit and monthly deposits. |
| rateOfReturn | string | The total in the account at the end of the term, divided by the total deposits made. A measure of how much return was made. The value is a percent (e.g. rateOfReturn = "5" means a 5% rate of return). |
| inflationAdjustedROR | string | The rate of return divided by the total inflation over the term. The value is a percent (e.g. rateOfReturn = "5" means a 5% rate of return). |
| plan | array of plan statuses | An array of monthly statuses of the Savings Plan. |

**Plan Status Fields**
| Field | Type | Description |
|---|---|---|
| date | string | The date of this monthly status. ISO 8601. |
| interest | integer | The interest earned this month in cents. |
| tax | integer | The tax paid this month in cents. |
| contribution | integer | The deposit made this month in cents. |
| increase | integer | The increase in savings at the end of this month, in cents. Includes deposits and interest earnings minus taxes. |
| capital | integer | The total money in the account at the end of this month, in cents. |

---

### Error Responses

### `400 Bad Request`

```json
{
  "error": "error message depending on the invalid or missing field"
}
```

### `401 Unauthorized`

```json
{
    "error": "error message dependimg on authentication error"
}
```

### `500 Internal Server Error`

```json
{
  "error": "internal server error"
}
```
## `POST /app/loans/calculate`

Calculate a Loan Payment Plan without saving it.

---

### URL

```http
POST /app/loans/calculate
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer token `Bearer [access token here]` |
| Content-Type | Yes | Usually `application/json` |

### Request Body

```json
{
	"startingPrincipal": 10000000,
	"yearlyInterestRate": "5",
	"monthlyPayment": 1500000,
	"escrowPayment": 10000,
	"startDate": "1970-01-01"
}
```

### Request Fields

| Field | Type | Required | Description |
|---|---|---|---|
| startingPrincipal | integer | Yes | Starting principal of the loan in cents. |
| yearlyInterestRate | string | Yes | Yearly interest rate of the loan (APR). Send as a percent ("5" would mean 5% yearly interest rate) |
| monthlyPayment | integer | Yes | Monthly payments made to the loan, in cents. |
| escrowPayment | integer | Yes | Additional, fixed payments that you make every month. For example, insurance, that are part of the monthly payments. |
| startDate | string | Yes | Start date of the loan. |

---

### Response

### Success Response

**Status Code:** `200 OK`

```json
{
    "durationMonths": 7,
    "totalExpenditure": 234054,
    "totalPaid": 10234054,
    "costOfCreditPercent": "1.02",
    "plan": [
        {
            "date": "1970-02-01T00:00:00Z",
            "payment": 1500000,
            "interest": 41667,
            "otherPayments": 10000,
            "paydown": 1448333,
            "principal": 8551667
        },
        {
            "date": "1970-03-01T00:00:00Z",
            "payment": 1500000,
            "interest": 35632,
            "otherPayments": 10000,
            "paydown": 1454368,
            "principal": 7097299
        }
    ]
}
```



### Response Fields

| Field | Type | Description |
|---|---|---|
| durationMonths | integer | Calculated duration of the loan in months. |
| totalExpenditure | integer | Total money paid to non-principal payments at the end of the loan. This includes interest and escrow payments. |
| totalPaid | integer | Total money paid at the end of the loan, including interest, principal and escrow payments. |
| costOfCreditPercent | string | How much more was paid than what was loaned. Calculated as the total paid over the loan's starting principal. For example "5" means a 5% cost of credit. |
| plan | array of plan statuses | An array of monthly statuses of the Payment plan. |

**Plan Status Fields**

| Field | Type | Description |
|---|---|---|
| date | string | The date of this monthly status. ISO 8601. |
| payment | integer | The payment made this month. |
| interest | integer | The interest accrued this month. |
| otherPayments | string | The amount that went into other payments this month (escrow payments like insurance). |
| paydown | array of plan statuses | The amount that was paid to principal this month. |
| principal | array of plan statuses | The loan principal at the end of this month. |

---

### Error Responses

### `400 Bad Request`

```json
{
  "error": "error message depending on the invalid or missing field"
}
```

### `401 Unauthorized`

```json
{
    "error": "error message depending on authentication error"
}
```

### `500 Internal Server Error`

```json
{
  "error": "internal server error"
}
```
## `POST /app/savings/save`

Calculate and then save a Savings Plan to database.

---

### URL

```http
POST /app/savings/save
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer token `Bearer [access token here]` |
| Content-Type | Yes | Usually `application/json` |

### Request Body

```json
{
    "name": "Test",
	"startingCapital": 700000,
	"yearlyInterestRate": "4.75",
    "interestRateType": "APR",
	"monthlyContribution": 15000,
	"durationYears": 1,
    "taxRate": "5",
    "yearlyInflationRate": "6",
	"startDate": "1970-01-01"
}
```

### Request Fields

| Field | Type | Required | Description |
|---|---|---|---|
| Name | string | Yes | The name of the savings plan. Doesn't need to be unique. |
| startingCapital | integer | Yes | How much money was deposited at the start of the term, in cents. For example, startingCapital = 100 would mean $1. |
| yearlyInterestRate | string | Yes | The yearly interest rate for the savings plan. Send as a percent. For example, "6.25" would be a 6.25% interest rate. |
| startDate | integer | Yes | The start date of the savings plan. |
| durationYears | integer | Yes | The term you want to calculate in years. 1 means "calculate the savings plan for 1 year".  ISO 8601|
| interestRateType | string | No | Send "APR" or "APY", depending on the type of interest rate. If empty, it defaults to APY. |
| monthlyContribution | integer | No | The monthly deposits that will be made (if any). Defaults to 0 if not in the request. The amount is in cents (e.g. 15000 means $150) |
| taxRate | string | No | The tax rate paid on returns. Send as a percent. For example, "5" means a 5% tax rate. Defaults to 0% if not in the request. |
| yearlyInflationRate | string | No | The yearly inflation rate, used for rate of return calculations. Send as a percent. For example, "6" means a 6% yearly inflation rate. Defaults to 0% if not in the request..  |

---

### Response

### Success Response

**Status Code:** `201 Created`

```json
{
    "id": "bd790a03-aabe-4a9d-861e-148fcd5adb46",
    "name": "test",
    "startingCapital": 700000,
    "yearlyInterestRate": "4.75",
    "interestRateType": "APY",
    "monthlyContribution": 15000,
    "durationYears": 1,
    "taxRate": "0",
    "yearlyInflationRate": "0",
    "startDate": "2026-01-31T18:00:00-06:00",
    "monthlyInterestRate": "0.3874684992129274856453709025059820372",
    "totalDeposited": 880000,
    "totalEarnings": 37136,
    "rateOfReturn": "4.22",
    "inflationAdjustedROR": "4.22"
}
```

### Response Fields

| Field | Type | Description |
|---|---|---|
| id | string | The ID of the created Savings Plan. UUID v4. |
| name | string | The name of the created loan. |
| startingCapital | int | Confirmation of the starting capital saved in the database |
| yearlyInterestRate | string | Confirmation of the yearly interest rate saved in the database |
| interestRateType | string | Confirmation of the interest rate type saved in the database. |
| monthlyContribution | integer | Confirmation of the monthly contribution saved in the database. |
| durationYears | integer | Confirmation of the duration saved in the database. |
| taxRate | string | Confirmation of the tax rate saved in the database. |
| yearlyInflationRate | string | Confirmation of the yearly inflation rate saved in the database. |
| startDate | string | Confirmation of the start date saved in the database. ISO 8601 timestamp. |
| monthlyInterestRate | string | The monthly interest rate calculated from the yearly one, used for calculations. |
| totalDeposited | integer | The total money deposited on the account in cents. |
| totalEarnings | string | The total money earned in interest after the term. |
| rateOfReturn | string | The total savings in the account divided by what was actually deposited. Represents a percent (4.22 means 4.22% rate of return) |
| inflationAdjustedROR | string | The rate of retrun divided by the inflation at the end of the term. Represents a percent (4.22 means 4.22% rate of return). |

---

### Error Responses

### `400 Bad Request`

```json
{
  "error": "error message depending on the invalid or missing field"
}
```

### `401 Unauthorized`

```json
{
    "error": "error message dependimg on authentication error"
}
```

### `500 Internal Server Error`

```json
{
  "error": "internal server error"
}
```
## `POST /app/loans/save`

Calculate and then save a Loan Payment Plan in database.

---

### URL

```http
POST /app/loans/save
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer token `Bearer [access token here]` |
| Content-Type | Yes | Usually `application/json` |

### Request Body

```json
{
    "name": "Test",
	"startingPrincipal": 10000000,
	"yearlyInterestRate": "5",
	"monthlyPayment": 1500000,
	"escrowPayment": 10000,
	"startDate": "1970-01-01"
}
```

### Request Fields

| Field | Type | Required | Description |
|---|---|---|---|
| name | string | Yes | Name for the new Loan Payment Plan. |
| startingPrincipal | integer | Yes | Starting principal of the loan in cents. |
| yearlyInterestRate | string | Yes | Yearly interest rate of the loan (APR). Send as a percent ("5" would mean 5% yearly interest rate) |
| monthlyPayment | integer | Yes | Monthly payments made to the loan, in cents. |
| escrowPayment | integer | Yes | Additional, fixed payments that you make every month. For example, insurance, that are part of the monthly payments. |
| startDate | string | Yes | Start date of the loan. |

---

### Response

### Success Response

**Status Code:** `201 Created`

```json
{
    "id": "7632d773-4846-4b47-b3b9-3f3f06330072",
    "name": "Test",
    "startingPrincipal": 10000000,
    "yearlyInterestRate": "5",
    "monthlyPayment": 900076,
    "escrowPayment": 10000,
    "startDate": "1969-12-31T18:00:00-06:00",
    "durationMonths": 12,
    "totalExpenditure": 383416,
    "totalPaid": 10383416,
    "costOfCreditPercent": "3.83416398261762"
}
```



### Response Fields

| Field | Type | Description |
|---|---|---|
| id | string | The ID of the created Loan Payment Plan. UUID v4. |
| name | string | The name of the created Loan Payment Plan. |
| startingPrincipal | integer | Confirmation of the starting principal saved on database. |
| yearlyInterestRate | integer | Confirmation of the yearly interest rate saved on database. |
| monthlyPayment | integer | Confirmation of the monthly payment saved on database. |
| escrowPayment | integer | Confirmation of the escrow payment saved on database. |
| startDate | integer | Confirmation of the start date saved on database. ISO 8601 timestamp. |
| durationMonths | integer | Calculated duration of the loan in months. |
| totalExpenditure | integer | The total non-principal expenditures at the end of the loan. Includes interest and escrow payments. |
| totalPaid | integer | The toal money paid on the loan. Includes interest, escrow and principal payments. |
| costOfCreditPercent | string | Calculated cost of credit as a percent. For example, "3.33" means a 3.33% cost of credit. |

---

### Error Responses

### `400 Bad Request`

```json
{
  "error": "error message depending on the invalid or missing field"
}
```

### `401 Unauthorized`

```json
{
    "error": "error message depending on authentication error"
}
```

### `500 Internal Server Error`

```json
{
  "error": "internal server error"
}
```
## `GET /app/savings/list`

List all the Savings Plans associated with a user.

---

### URL

```http
GET /app/savings/list
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer token `Bearer [access token here]` |

### Request Body

No body required. The user information is extracted from the Access Token.

---

### Response

### Success Response

**Status Code:** `200 OK`

```json
{
    "plans": [
        {
            "id": "a5b140fd-583c-4edb-8668-e7ee986d2a37",
            "name": "test",
            "startingCapital": 700000
        },
        {
            "id": "479d382b-2129-4d91-a77f-9ea1f90acd5b",
            "name": "test",
            "startingCapital": 700000
        }
    ]
}
```

### Response Fields

| Field | Type | Description |
|---|---|---|
| plans | Array of plans data. | The list of savings plans associated with the user. |

**Plan data**

| Field | Type | Description |
|---|---|---|
| id | string | The plan's ID. UUID v4. |
| name | string | The plan's name. |
| startingCapital | integer | The plan's starting capital in cents. |

---

### Error Responses

### `401 Unauthorized`

```json
{
    "error": "error message depending on authentication error"
}
```

### `500 Internal Server Error`

```json
{
  "error": "internal server error"
}
```
## `GET /app/loans/list`

List all the loans associated with a user.

---

### URL

```http
GET /app/loans/list
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer token `Bearer [access token here]` |

### Request Body

No body required. The user information is extracted from the Access Token.

---

### Response

### Success Response

**Status Code:** `200 OK`

```json
{
    "loans": [
        {
            "id": "d25d1653-6ac5-4b97-b5db-b164246c4b88",
            "name": "Test 2",
            "loanAmount": 10000000
        },
        {
            "id": "ef375465-08ad-4947-9541-e192367797f7",
            "name": "Test 2",
            "loanAmount": 10000000
        }
    ]
}
```

### Response Fields

| Field | Type | Description |
|---|---|---|
| plans | Array of plans data. | The list of Loan Payment Plans associated with the user. |

**Plan data**

| Field | Type | Description |
|---|---|---|
| id | string | The plan's ID. UUID v4. |
| name | string | The plan's name. |
| loanAmount | integer | The plan's starting principal in cents. |

---

### Error Responses

### `401 Unauthorized`

```json
{
    "error": "error message depending on authentication error"
}
```

### `500 Internal Server Error`

```json
{
  "error": "internal server error"
}
```
## `GET /app/savings/{id}`

Retrieve the data of the specified Savings Plan

---

### URL

```http
GET /app/savings/{id}
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer token `Bearer [access token here]` |

### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| id | string | Yes | UUID of the requested Savings Plan |

### Request Body

None.

---

### Response

### Success Response

**Status Code:** `200 OK`

```json
{
    "id": "a5b140fd-583c-4edb-8668-e7ee986d2a37",
    "name": "test",
    "originalData": {
        "startingCapital": 700000,
        "yearlyInterestRate": "4.75",
        "interestRateType": "APY",
        "monthlyContribution": 15000,
        "durationYears": 1,
        "taxRate": "0",
        "yearlyInflationRate": "0",
        "startDate": "2026-01-31T18:00:00-06:00"
    },
    "calculatedData": {
        "monthlyInterestRate": "0.3874684992",
        "totalInterestEarnings": 37136,
        "totalDeposited": 0,
        "rateOfReturn": "4.22",
        "inflationAdjustedROR": "4.22"
    },
    "plan": [
        {
            "date": "2026-02-28T18:00:00-06:00",
            "interest": 2712,
            "tax": 0,
            "contribution": 15000,
            "increase": 17712,
            "capital": 717712
        }
    ]
}
```

### Response Fields

| Field | Type | Description |
|---|---|---|
| id | string | Resource ID |
| originalData | Object | The original plan data, received when it was created. Consult the `/app/savings/calculate` documentation for more information on the fields. |
| calculatedData | Object | The data calculated from the initial data. Consult the `/app/savings/calculate` documentation for more information on the fields. |
| plan | Array of Objects | The monthly statuses of the savings plan. Consult the `/app/savings/calculate` documentation for more information on the fields. |

---

### Error Responses

### `401 Unauthorized`

Returned either when the resource is not found, or the user was unauthorized.

### `500 Internal Server Error`

```json
{
  "error": "internal server error"
}
```
## `GET /app/loans/{id}`

Retrieve the data of the specified Loan Payment Plan.

---

### URL

```http
GET /app/loans/{id}
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer token `Bearer [access token here]` |

### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| id | string | Yes | UUID of the requested Loan Payment Plan. |

### Request Body

None.

---

### Response

### Success Response

**Status Code:** `200 OK`

```json
{
    "id": "d25d1653-6ac5-4b97-b5db-b164246c4b88",
    "name": "Test 2",
    "originalData": {
        "startingPrincipal": 10000000,
        "yearlyInterestRate": "5",
        "monthlyPayment": 900076,
        "escrowPayment": 10000,
        "startDate": "1969-12-31T18:00:00-06:00"
    },
    "calculatedData": {
        "monthlyInterestRate": "0.4166666667",
        "durationMonths": 12,
        "totalExpenditure": 383416,
        "totalPaid": 10383416,
        "costOfCredit": "1.0383416398261762"
    },
    "paymentPlan": [
        {
            "date": "1970-01-31T18:00:00-06:00",
            "payment": 900076,
            "interest": 41667,
            "otherPayments": 10000,
            "paydown": 848409,
            "principal": 9151591
        }
    ]
}
```

### Response Fields

| Field | Type | Description |
|---|---|---|
| id | string | Resource ID |
| originalData | Object | The original plan data, received when it was created. Consult the `/app/loans/calculate` documentation for more information on the fields. |
| calculatedData | Object | The data calculated from the initial data. Consult the `/app/loans/calculate` documentation for more information on the fields. |
| plan | Array of Objects | The monthly statuses of the savings plan. Consult the `/app/loans/calculate` documentation for more information on the fields. |

---

### Error Responses

### `401 Unauthorized`

Returned either when the resource is not found, or the user was unauthorized.

### `500 Internal Server Error`

```json
{
  "error": "internal server error"
}
```
## `PATCH /app/savings/{id}`

Update the initial data of a Savings Plan, recalculate and save changes to database.

---

### URL

```http
PATCH /app/savings/{id}
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer token `Bearer [access token here]` |
| Content-Type | Yes |  `application/json` |

### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| id | string | Yes | UUID of the requested Savings Plan |

### Request Body

```json
{
  "name": "Test",
  "startingCapital": 700000,
  "yearlyInterestRate": "4.75",
  "interestRateType": "APR",
  "monthlyContribution": 15000,
  "durationYears": 1,
  "taxRate": "5",
  "yearlyInflationRate": "6",
  "startDate": "1970-01-01"
}
```

### Request Fields

| Field | Type | Required | Description |
|---|---|---|---|
| Name | string | No | Updated name for the savings plan. |
| startingCapital | integer | No | Updated starting capital for the savings plan, in cents. |
| yearlyInterestRate | string | No | Updated yearly interest rate for the savings plan, as a percent. Example: "5" means 5%. |
| startDate | integer | No | Updated start date for the savings plan. YYYY-MM-DD format. |
| durationYears | integer | No | Updated duration for the savings plan, in years. |
| interestRateType | string | No | Updated interest rate type for the savings plan. "APR" or "APY". |
| monthlyContribution | integer | No | Updated monthly contribution, in cents. |
| taxRate | string | No | Updated tax rate, as a percent. Example: "5" means 5%. |
| yearlyInflationRate | string | No | Updated yearly inflation rate, as a percent. Example: "5" means 5%.  |

---

### Response

### Success Response

**Status Code:** `200 OK`

```json
{
    "id": "bea4b58a-86a0-4eee-8380-37929b599216",
    "name": "test",
    "startingCapital": 700000,
    "yearlyInterestRate": "4.2",
    "interestRateType": "APY",
    "monthlyContribution": 15000,
    "durationYears": 1,
    "taxRate": "0",
    "yearlyInflationRate": "0",
    "startDate": "2026-01-23T18:00:00-06:00",
    "monthlyInterestRate": "0.3434379290046821080713036287486481524",
    "totalDeposited": 880000,
    "totalEarnings": 32839,
    "rateOfReturn": "3.73",
    "inflationAdjustedROR": "3.73"
}
```
### Response Fields

| Field | Type | Description |
|---|---|---|
| id | string | The ID of the updated Savings Plan. UUID v4. |
| name | string | The name of the updated savins plan. |
| startingCapital | integer | Current starting capital of the savings plan. |
| yearlyInterestRate | string | Current yearly interest rate of the savings plan. |
| interestRateType | string | Current interest rate type of the savings plan. |
| monthlyContribution | integer | Current monthly contribution of the savings plan. |
| durationYears | integer | Current duration of the savings plan. |
| taxRate | string | Current tax rate of the savings plan. |
| yearlyInflationRate | string | Current yearly inflation rate of the savings plan. |
| startDate | string | Current start date of the savings plan.. ISO 8601 timestamp. |
| monthlyInterestRate | string | Recalculated monthly interest rate. |
| totalDeposited | integer | Recalculated total deposited. |
| totalEarnings | string | Recalculated total earnings. |
| rateOfReturn | string | Recalculated rate of return. |
| inflationAdjustedROR | string | Recalculated inflation adjusted rate of return. |

---

### Error Responses

### `400 Bad Request`

```json
{
  "error": "error message depending on the invalid or missing field"
}
```

### `401 Unauthorized`

```json
{
    "error": "error message dependimg on authentication error"
}
```

### `500 Internal Server Error`

```json
{
  "error": "internal server error"
}
```
## `PATCH /app/loans/{id}`

Update the initial data of a Loan Payment Plan, recalculate and save changes to database.

---

### URL

```http
PATCH /app/loans/{id}
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer token `Bearer [access token here]` |
| Content-Type | Yes |  `application/json` |

### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| id | string | Yes | UUID of the requested Loan Payment Plan |

### Request Body

```json
{
  "name": "Test",
  "startingPrincipal": 10000000,
  "yearlyInterestRate": "5",
  "monthlyPayment": 1500000,
  "escrowPayment": 10000,
  "startDate": "1970-01-01"
}
```

### Request Fields

| Field | Type | Required | Description |
|---|---|---|---|
| Name | string | No | Updated name for the savings plan. |
| startingPrincipal | integer | No | Updated starting principal for the loan payment plan, in cents. |
| yearlyInterestRate | string | No | Updated yearly interest rate for the loan payment plan, as a percent. Example: "5" means 5%. |
| monthlyPayment | integer | No | Updated monthly payment for the loan payment plan, in cents. |
| escrowPayment | integer | No | Updated amount of other payments (not interest or principal) for the loan payment plan. |
| startDate | string | No | Updated start date for the loan payment plan. YYYY-MM-DD format. |

---

### Response

### Success Response

**Status Code:** `200 OK`

```json
{
    "id": "ef375465-08ad-4947-9541-e192367797f7",
    "name": "Test 2",
    "startingPrincipal": 10000000,
    "yearlyInterestRate": "4",
    "monthlyPayment": 900076,
    "escrowPayment": 10000,
    "startDate": "1969-12-23T18:00:00-06:00",
    "durationMonths": 12,
    "totalExpenditure": 329407,
    "totalPaid": 10329407,
    "costOfCreditPercent": "3.2940656243549"
}

```

### Response Fields

| Field | Type | Description |
|---|---|---|
| id | string | The ID of the updated Savings Plan. UUID v4. |
| name | string | The name of the updated savins plan. |
| startingPrincipal | integer | Current starting starting principal of the loan payment plan. |
| yearlyInterestRate | string | Current yearly interest rate of the loan payment plan. |
| monthlyPayment | integer | Current monthly payment rate of the loan payment plan. |
| escrowPayment | integer | Current amount of other payments of the loan payment plan. |
| startDate | string | Current start date of the loan payment plan. YYYY-MM-DD format. |
| durationMonths | integer | Recalculated duration of the loan payment plan. |
| totalExpenditure | integer | Recalculated total expenditure of the loan payment plan. |
| totalPaid | integer | Recalculated total paid of the loan payment plan. |
| costOfCreditPercent | string | Recalculated cost of credit of the loan payment plan. |

---

### Error Responses

### `400 Bad Request`

```json
{
  "error": "error message depending on the invalid or missing field"
}
```

### `401 Unauthorized`

```json
{
    "error": "error message depending on authentication error"
}
```

### `500 Internal Server Error`

```json
{
  "error": "internal server error"
}
```
## `DELETE /app/savings/{id}`

Delete the specified Savings Plan.

---

### URL

```http
DELETE /app/savings/{id}
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer token `Bearer [access token here]` |
| Content-Type | Yes |  `application/json` |

### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| id | string | Yes | UUID of the Savings Plan to be deleted. |

### Request Body

None.

### Response

### Success Response

**Status Code:** `204 No Content`

```json
No body.
```

---

### Error Responses

### `401 Unauthorized`

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
## `DELETE /app/loans/{id}`

Delete the specified Loan Payment Plan.

---

### URL

```http
DELETE /app/loans/{id}
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer token `Bearer [access token here]` |
| Content-Type | Yes |  `application/json` |

### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| id | string | Yes | UUID of the Loan Payment Plan to be deleted. |

### Request Body

None.

### Response

### Success Response

**Status Code:** `204 No Content`

```json
No body.
```

---

### Error Responses

### `401 Unauthorized`

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

## Contributing

### Clone the repo

```bash
git clone https://github.com/Mr-Rafael/finance-calculator.git
cd finance-calculator
```

### Run the unit test suite

```bash
go test ./internal/...
```

### Build the compiled binary

```bash
go build
```

### Submit a pull request

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.

# Collaborators</h2>
<table>
  <tr>
    <td align="center">
      <a href="#">
        <img src="https://avatars.githubusercontent.com/u/35672719?s=48&v=4" width="100px;" alt="Rafael Mazariegos picture"/><br>
        <sub>
          <b>Rafael Mazariegos (Mr-Rafael)</b>
        </sub>
      </a>
    </td>
  </tr>
</table>
