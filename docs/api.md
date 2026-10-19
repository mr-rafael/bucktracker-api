# API Endpoints

| route | description |
|---|---|
| <kbd>GET /api/healthz</kbd> | Check if server is running. |
| <kbd>POST /app/users/create</kbd> | Create a new user. |
| <kbd>POST /app/login</kbd> | Login user. |
| <kbd>POST /app/refresh</kbd> | Refresh the access token. |
| <kbd>POST /app/savings/calculate</kbd> | Generate a Savings Plan without saving. |
| <kbd>POST /app/loans/calculate</kbd> | Calculate a Loan Payment Plan without saving. |
| <kbd>POST /app/savings/save</kbd> | Calculate and save a Savings Plan. |
| <kbd>POST /app/loans/save</kbd> | Calculate and save a Loan. |
| <kbd>GET /app/savings/list</kbd> | List a User's saved Savings Plans. |
| <kbd>GET /app/loans/list</kbd> | List a User's saved Loans. |
| <kbd>GET /app/savings/{id}</kbd> | Get a previously saved Savings Plan. |
| <kbd>GET /app/loans/{id}</kbd> | Get a previously saved Loan. |
| <kbd>GET /app/loans/{loanId}/payment-plans/{paymentPlanId}</kbd> | Get a payment plan for a Loan. |
| <kbd>POST /app/loans/{loanId}/payment-plans</kbd> | Create a payment plan for a Loan. |
| <kbd>PATCH /app/loans/{loanId}/payment-plans/{paymentPlanId}</kbd> | Update a payment plan for a Loan. |
| <kbd>DELETE /app/loans/{loanId}/payment-plans/{paymentPlanId}</kbd> | Delete a payment plan from a Loan. |
| <kbd>PATCH /app/savings/{id}</kbd> | Update and recalculate a Savings Plan. |
| <kbd>PATCH /app/loans/{id}</kbd> | Update a Loan and recalculate its payment plans. |
| <kbd>DELETE /app/savings/{id}</kbd> | Delete a Savings Plan. |
| <kbd>DELETE /app/loans/{id}</kbd> | Delete a Loan. |

Money amounts are integers in **cents**. Interest rates and similar percentages are sent as **strings** (for example, `"5"` means 5%).

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
## `POST /app/login`

Login to receive an access and refresh token. Access tokens are necessary to access some endpoints. The refresh token is set as an HTTP-only cookie named `refresh_token`.

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

### `400 Bad Request`

Returned when the request body cannot be decoded.

### `401 Unauthorized`

Returned when the email or password is invalid.
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
| Cookie | Yes | The cookie should contain `refresh_token`, set by the Login endpoint. |

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
| Content-Type | Yes | `application/json` |

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
| startDate | string | Yes | The start date of the savings plan. YYYY-MM-DD. |
| durationYears | integer | Yes | The term you want to calculate in years. 1 means "calculate the savings plan for 1 year". |
| interestRateType | string | No | Send "APR" or "APY", depending on the type of interest rate. If empty, it defaults to APY. |
| monthlyContribution | integer | No | The monthly deposits that will be made (if any). Defaults to 0 if not in the request. The amount is in cents (e.g. 15000 means $150) |
| taxRate | string | No | The tax rate paid on returns. Send as a percent. For example, "5" means a 5% tax rate. Defaults to 0% if not in the request. |
| yearlyInflationRate | string | No | The yearly inflation rate, used for rate of return calculations. Send as a percent. For example, "6" means a 6% yearly inflation rate. Defaults to 0% if not in the request.  |

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
| Content-Type | Yes | `application/json` |

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
| startDate | string | Yes | Start date of the loan. YYYY-MM-DD. |

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
| otherPayments | integer | The amount that went into other payments this month (escrow payments like insurance). |
| paydown | integer | The amount that was paid to principal this month. |
| principal | integer | The loan principal at the end of this month. |

---

### Error Responses

### `400 Bad Request`

```json
{
  "error": "error message depending on the invalid or missing field"
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
| name | string | Yes | The name of the savings plan. Doesn't need to be unique. |
| startingCapital | integer | Yes | How much money was deposited at the start of the term, in cents. For example, startingCapital = 100 would mean $1. |
| yearlyInterestRate | string | Yes | The yearly interest rate for the savings plan. Send as a percent. For example, "6.25" would be a 6.25% interest rate. |
| startDate | string | Yes | The start date of the savings plan. YYYY-MM-DD. |
| durationYears | integer | Yes | The term you want to calculate in years. 1 means "calculate the savings plan for 1 year". |
| interestRateType | string | No | Send "APR" or "APY", depending on the type of interest rate. If empty, it defaults to APY. |
| monthlyContribution | integer | No | The monthly deposits that will be made (if any). Defaults to 0 if not in the request. The amount is in cents (e.g. 15000 means $150) |
| taxRate | string | No | The tax rate paid on returns. Send as a percent. For example, "5" means a 5% tax rate. Defaults to 0% if not in the request. |
| yearlyInflationRate | string | No | The yearly inflation rate, used for rate of return calculations. Send as a percent. For example, "6" means a 6% yearly inflation rate. Defaults to 0% if not in the request.  |

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
    "startDate": "2026-01-31",
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
| name | string | The name of the created savings plan. |
| startingCapital | integer | Confirmation of the starting capital saved in the database |
| yearlyInterestRate | string | Confirmation of the yearly interest rate saved in the database |
| interestRateType | string | Confirmation of the interest rate type saved in the database. |
| monthlyContribution | integer | Confirmation of the monthly contribution saved in the database. |
| durationYears | integer | Confirmation of the duration saved in the database. |
| taxRate | string | Confirmation of the tax rate saved in the database. |
| yearlyInflationRate | string | Confirmation of the yearly inflation rate saved in the database. |
| startDate | string | Confirmation of the start date saved in the database. YYYY-MM-DD. |
| monthlyInterestRate | string | The monthly interest rate calculated from the yearly one, used for calculations. |
| totalDeposited | integer | The total money deposited on the account in cents. |
| totalEarnings | integer | The total money earned in interest after the term. |
| rateOfReturn | string | The total savings in the account divided by what was actually deposited. Represents a percent (4.22 means 4.22% rate of return) |
| inflationAdjustedROR | string | The rate of return divided by the inflation at the end of the term. Represents a percent (4.22 means 4.22% rate of return). |

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
## `POST /app/loans/save`

Calculate and then save a Loan. A newly created loan starts with a default payment plan named `Default Payment Plan` and no extra principal payments.

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
| name | string | Yes | Name for the new Loan. |
| startingPrincipal | integer | Yes | Starting principal of the loan in cents. |
| yearlyInterestRate | string | Yes | Yearly interest rate of the loan (APR). Send as a percent ("5" would mean 5% yearly interest rate) |
| monthlyPayment | integer | Yes | Monthly payments made to the loan, in cents. |
| escrowPayment | integer | Yes | Additional, fixed payments that you make every month. For example, insurance, that are part of the monthly payments. |
| startDate | string | Yes | Start date of the loan. YYYY-MM-DD. |

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
    "defaultPaymentPlan": {
        "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "name": "Default Payment Plan",
        "durationMonths": 12
    },
    "paymentPlans": [
        {
            "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
            "name": "Default Payment Plan",
            "durationMonths": 12
        }
    ]
}
```



### Response Fields

| Field | Type | Description |
|---|---|---|
| id | string | The ID of the created Loan. UUID v4. |
| name | string | The name of the created Loan. |
| startingPrincipal | integer | Confirmation of the starting principal saved on database. |
| yearlyInterestRate | string | Confirmation of the yearly interest rate saved on database. |
| defaultPaymentPlan | object | Summary of the default payment plan created with the loan. |
| paymentPlans | array of payment plan summaries | All payment plans on the loan. A newly created loan contains only the default plan. |

**Payment Plan Summary Fields**

| Field | Type | Description |
|---|---|---|
| id | string | The payment plan's ID. UUID v4. |
| name | string | The payment plan's name. |
| durationMonths | integer | Calculated duration of this payment plan in months. |

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
| loans | Array of loan data. | The list of Loans associated with the user. |

**Loan data**

| Field | Type | Description |
|---|---|---|
| id | string | The loan's ID. UUID v4. |
| name | string | The loan's name. |
| loanAmount | integer | The loan's starting principal in cents. |

---

### Error Responses

### `401 Unauthorized`

```json
{
    "error": "error message depending on authentication error"
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
        "startDate": "2026-01-31"
    },
    "calculatedData": {
        "monthlyInterestRate": "0.3874684992",
        "totalInterestEarnings": 37136,
        "totalDeposited": 880000,
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
| name | string | The savings plan's name. |
| originalData | Object | The original plan data, received when it was created. Consult the `/app/savings/calculate` documentation for more information on the fields. |
| calculatedData | Object | The data calculated from the initial data. Consult the `/app/savings/calculate` documentation for more information on the fields. `totalInterestEarnings` is the same value as `totalEarnings` on the calculate endpoint. |
| plan | Array of Objects | The monthly statuses of the savings plan. Consult the `/app/savings/calculate` documentation for more information on the fields. |

---

### Error Responses

### `401 Unauthorized`

Returned either when the resource is not found, or the user was unauthorized.
## `GET /app/loans/{id}`

Retrieve the data of the specified Loan, including summaries of its payment plans.

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
| id | string | Yes | UUID of the requested Loan. |

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
        "startDate": "1970-01-01"
    },
    "defaultPaymentPlan": {
        "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "name": "Default Payment Plan",
        "durationMonths": 12
    },
    "paymentPlans": [
        {
            "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
            "name": "Default Payment Plan",
            "durationMonths": 12
        }
    ]
}
```

### Response Fields

| Field | Type | Description |
|---|---|---|
| id | string | Resource ID |
| name | string | The loan's name. |
| originalData | Object | The original loan data, received when it was created. Consult the `/app/loans/calculate` documentation for more information on the fields. |
| defaultPaymentPlan | Object | Summary of the loan's default payment plan. Consult the `/app/loans/save` documentation for the summary fields. |
| paymentPlans | Array of Objects | Summaries of every payment plan on the loan. Use `GET /app/loans/{loanId}/payment-plans/{paymentPlanId}` to retrieve a plan's monthly breakdown. |

---

### Error Responses

### `401 Unauthorized`

Returned either when the resource is not found, or the user was unauthorized.
## `GET /app/loans/{loanId}/payment-plans/{paymentPlanId}`

Retrieve a specific payment plan for a Loan, including extra principal payments and the monthly breakdown.

---

### URL

```http
GET /app/loans/{loanId}/payment-plans/{paymentPlanId}
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer token `Bearer [access token here]` |

### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| loanId | string | Yes | UUID of the Loan that owns the payment plan. |
| paymentPlanId | string | Yes | UUID of the requested payment plan. |

### Request Body

None.

---

### Response

### Success Response

**Status Code:** `200 OK`

```json
{
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "name": "Default Payment Plan",
    "durationMonths": 12,
    "totalExpenditure": 120000,
    "totalPaid": 10800000,
    "costOfCredit": "8",
    "principalPayments": [
        {
            "date": "1970-02-01",
            "amount": 500000
        }
    ],
    "paymentPlanBreakdown": [
        {
            "date": "1970-02-01T00:00:00Z",
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
| id | string | The payment plan's ID. UUID v4. |
| name | string | The payment plan's name. |
| durationMonths | integer | Calculated duration of the payment plan in months. |
| totalExpenditure | integer | Total money paid to non-principal payments at the end of the plan. This includes interest and escrow payments. |
| totalPaid | integer | Total money paid at the end of the plan, including interest, principal and escrow payments. |
| costOfCredit | string | How much more was paid than what was loaned. |
| principalPayments | array of principal payments | Extra principal payments applied to this plan. |
| paymentPlanBreakdown | array of plan statuses | An array of monthly statuses of the payment plan. Consult the `/app/loans/calculate` documentation for the status fields. |

**Principal Payment Fields**

| Field | Type | Description |
|---|---|---|
| date | string | Date of the extra principal payment. YYYY-MM-DD. |
| amount | integer | Extra principal payment amount, in cents. |

---

### Error Responses

### `401 Unauthorized`

Returned either when the resource is not found, or the user was unauthorized.
## `POST /app/loans/{loanId}/payment-plans`

Create a new payment plan for an existing Loan. The plan is calculated from the loan's current terms plus any extra principal payments in the request.

---

### URL

```http
POST /app/loans/{loanId}/payment-plans
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer token `Bearer [access token here]` |
| Content-Type | Yes | `application/json` |

### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| loanId | string | Yes | UUID of the Loan that will own the new payment plan. |

### Request Body

```json
{
    "name": "Extra Principal Plan",
    "principalPayments": [
        {
            "date": "1970-02-01",
            "amount": 500000
        }
    ]
}
```

### Request Fields

| Field | Type | Required | Description |
|---|---|---|---|
| name | string | Yes | Name for the new payment plan. |
| principalPayments | array of principal payments | No | Extra principal payments to apply while calculating the plan. Dates must be YYYY-MM-DD. Amounts are in cents. |

---

### Response

### Success Response

**Status Code:** `201 Created`

The response body uses the same shape as `GET /app/loans/{loanId}/payment-plans/{paymentPlanId}`.

```json
{
    "id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "name": "Extra Principal Plan",
    "durationMonths": 10,
    "totalExpenditure": 95000,
    "totalPaid": 10500000,
    "costOfCredit": "5",
    "principalPayments": [
        {
            "date": "1970-02-01",
            "amount": 500000
        }
    ],
    "paymentPlanBreakdown": [
        {
            "date": "1970-02-01T00:00:00Z",
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

Consult the `GET /app/loans/{loanId}/payment-plans/{paymentPlanId}` documentation for the response fields.

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
## `PATCH /app/loans/{loanId}/payment-plans/{paymentPlanId}`

Update a payment plan's name and/or replace its extra principal payments. If `principalPayments` is sent, the plan is recalculated. Sending an empty body is rejected.

---

### URL

```http
PATCH /app/loans/{loanId}/payment-plans/{paymentPlanId}
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer token `Bearer [access token here]` |
| Content-Type | Yes | `application/json` |

### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| loanId | string | Yes | UUID of the Loan that owns the payment plan. |
| paymentPlanId | string | Yes | UUID of the payment plan to update. |

### Request Body

```json
{
    "name": "New Name",
    "principalPayments": [
        {
            "date": "1970-02-01",
            "amount": 500000
        }
    ]
}
```

### Request Fields

| Field | Type | Required | Description |
|---|---|---|---|
| name | string | No | Updated name for the payment plan. |
| principalPayments | array of principal payments | No | Replaces the plan's extra principal payments, then recalculates the plan. Dates must be YYYY-MM-DD. Amounts are in cents. |

At least one of `name` or `principalPayments` must be present.

---

### Response

### Success Response

**Status Code:** `200 OK`

The response body uses the same shape as `GET /app/loans/{loanId}/payment-plans/{paymentPlanId}`.

```json
{
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "name": "New Name",
    "durationMonths": 12,
    "totalExpenditure": 120000,
    "totalPaid": 10800000,
    "costOfCredit": "8",
    "principalPayments": [
        {
            "date": "1970-02-01",
            "amount": 500000
        }
    ],
    "paymentPlanBreakdown": [
        {
            "date": "1970-02-01T00:00:00Z",
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

Consult the `GET /app/loans/{loanId}/payment-plans/{paymentPlanId}` documentation for the response fields.

---

### Error Responses

### `400 Bad Request`

```json
{
  "error": "empty update request"
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
## `DELETE /app/loans/{loanId}/payment-plans/{paymentPlanId}`

Delete the specified payment plan.

---

### URL

```http
DELETE /app/loans/{loanId}/payment-plans/{paymentPlanId}
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer token `Bearer [access token here]` |

### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| loanId | string | Yes | UUID of the Loan that owns the payment plan. |
| paymentPlanId | string | Yes | UUID of the payment plan to delete. |

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

Returned when the user is unauthorized or a path ID is invalid.

### `404 Not Found`

Returned when the payment plan could not be deleted.
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
| name | string | No | Updated name for the savings plan. |
| startingCapital | integer | No | Updated starting capital for the savings plan, in cents. |
| yearlyInterestRate | string | No | Updated yearly interest rate for the savings plan, as a percent. Example: "5" means 5%. |
| startDate | string | No | Updated start date for the savings plan. YYYY-MM-DD format. |
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
    "startDate": "2026-01-23",
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
| name | string | The name of the updated savings plan. |
| startingCapital | integer | Current starting capital of the savings plan. |
| yearlyInterestRate | string | Current yearly interest rate of the savings plan. |
| interestRateType | string | Current interest rate type of the savings plan. |
| monthlyContribution | integer | Current monthly contribution of the savings plan. |
| durationYears | integer | Current duration of the savings plan. |
| taxRate | string | Current tax rate of the savings plan. |
| yearlyInflationRate | string | Current yearly inflation rate of the savings plan. |
| startDate | string | Current start date of the savings plan. YYYY-MM-DD. |
| monthlyInterestRate | string | Recalculated monthly interest rate. |
| totalDeposited | integer | Recalculated total deposited. |
| totalEarnings | integer | Recalculated total earnings. |
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
    "error": "error message depending on authentication error"
}
```

### `500 Internal Server Error`

```json
{
  "error": "internal server error"
}
```
## `PATCH /app/loans/{id}`

Update the initial data of a Loan and save the changes. If loan terms change (principal, interest rate, monthly payment, escrow, or start date), every payment plan on the loan is recalculated.

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
| id | string | Yes | UUID of the requested Loan |

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
| name | string | No | Updated name for the loan. |
| startingPrincipal | integer | No | Updated starting principal for the loan, in cents. |
| yearlyInterestRate | string | No | Updated yearly interest rate for the loan, as a percent. Example: "5" means 5%. |
| monthlyPayment | integer | No | Updated monthly payment for the loan, in cents. |
| escrowPayment | integer | No | Updated amount of other payments (not interest or principal) for the loan. |
| startDate | string | No | Updated start date for the loan. YYYY-MM-DD format. |

---

### Response

### Success Response

**Status Code:** `200 OK`

The response body uses the same shape as `GET /app/loans/{id}`.

```json
{
    "id": "ef375465-08ad-4947-9541-e192367797f7",
    "name": "Test 2",
    "originalData": {
        "startingPrincipal": 10000000,
        "yearlyInterestRate": "4",
        "monthlyPayment": 900076,
        "escrowPayment": 10000,
        "startDate": "1970-01-01"
    },
    "defaultPaymentPlan": {
        "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "name": "Default Payment Plan",
        "durationMonths": 12
    },
    "paymentPlans": [
        {
            "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
            "name": "Default Payment Plan",
            "durationMonths": 12
        }
    ]
}
```

### Response Fields

Consult the `GET /app/loans/{id}` documentation for the response fields.

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

### `404 Not Found`

Returned when the loan ID in the URL is not a valid UUID.

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

Returned when the user is unauthorized or the plan ID is invalid.

### `404 Not Found`

Returned when the savings plan could not be deleted.
## `DELETE /app/loans/{id}`

Delete the specified Loan.

---

### URL

```http
DELETE /app/loans/{id}
```

### Headers

| Header | Required | Description |
|---|---|---|
| Authorization | Yes | Bearer token `Bearer [access token here]` |

### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| id | string | Yes | UUID of the Loan to be deleted. |

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

Returned when the user is unauthorized or the loan ID is invalid.

### `404 Not Found`

Returned when the loan could not be deleted.
