# BuckTracker

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

## Clone the repository

```bash
git clone git@github.com:Mr-Rafael/bucktracker-api.git
cd bucktracker-api
```

## Run with Docker Compose

From the project root, start the stack:

```bash
docker compose up --build -d
```

This command builds the images and starts the following containers:

| Container | Role |
|---|---|
| `db` | PostgreSQL 15 database on port `5432`. Data is stored in the `bucktracker-api_postgres_data` volume. |
| `migrate` | One-shot Goose job that applies the SQL migrations in `backend/internal/db/migrations`, then exits. |
| `backend` | The BuckTracker REST API server. |
| `frontend` | The BuckTracker web UI, served by nginx. |

## Access the app

Once the containers are up:

- **API:** [http://localhost:8080](http://localhost:8080) — health check at [http://localhost:8080/api/healthz](http://localhost:8080/api/healthz)
- **Frontend:** [http://localhost:3000](http://localhost:3000)

To stop the stack:

```bash
docker compose down
```

# Using the App

Access the app at: [http://localhost:3000](http://localhost:3000). You will be shown a **Welcome screen**:

<figure style="text-align: center;">
  <img src="docs/media/welcome_screen.png" alt="Welcome screen" width="500">
  <figcaption style="font-style: italic; color: #555; margin-top: 8px;">
    Welcome screen
  </figcaption>
</figure>

You can check if the backend is running at the Welcome Page.

## Using the Calculators

You can navigate to the loan and savings calculators without having a user, from the Welcome screen.

### Savings Calculator

<figure style="text-align: center;">
  <img src="docs/media/savings_calculator_1.png" alt="Savings Calculator Input" width="500">
  <figcaption style="font-style: italic; color: #555; margin-top: 8px;">
    Savings Calculator Input
  </figcaption>
</figure>

The **savings calculator** allows you to calculate your savings in an account after a determined period. You can specify the following fields:

| Field | Description |
| ----------- | ----------- |
| Starting Capital | The amount in your account at the beginning of the term. |
| Monthly Contribution | The monthly amount you will deposit in the account (can be 0). |
| Yearly Interest Rate | The yearly interest rate of your account. |
| Interest Rate Type | Banks normally provide the yearly interest rate as APY (Annual Percentage Yield), but in case it gives an APR (Annual Percentage Rate), you can set it here. |
| Tax Rate | You can include an Income Tax Rate in the calculation. You can leave it as 0 to ignore it. |
| Yearly Inflation Rate | You can include a Yearly Inflation Rate in the calculation, to see the real growth of your savings. You can leave it as 0 to ignore it. |
| Duration Years | The duration of the term you wish to calculate for. |
| Start Date | The start date of the savings investment. |

Clicking on **Calculate Savings Plan** will run calculations and return the following information:

<figure style="text-align: center;">
  <img src="docs/media/savings_calculator_2.png" alt="Savings Calculator Output" width="500">
  <figcaption style="font-style: italic; color: #555; margin-top: 8px;">
    Savings Calculator Output
  </figcaption>
</figure>

| Field | Description |
| ----------- | ----------- |
| Monthly Interest Rate | The yearly interest rate you entered, converted to a monthly one. This is used in the actual month-to-month interest calculations. |
| Total Interest Earnings | The sum of all the interest amounts generated during the period. |
| Total Deposited | The sum of all the deposits you made during the period, including the initial deposit. |
| Rate of Return | Your total savings at the end of the period divided by the total deposited as a percent. |
| Inflation Adjusted Rate of Return | How much your savings really grew, when considering the inflation you entered. |

Additionally, it will generate a month-to-month status report of your savings account, with the following fields:

| Field | Description |
| ----------- | ----------- |
| Date | The date of the status. |
| Interest | The interest earned that month. |
| Tax | The taxes deducted from the interest profits that month. |
| Contribution | The amount you will deposit that month. |
| Increase | By how much your savings increased that month. |
| Capital | The total on your account at the end of that month. |

### Loans Calculator

<figure style="text-align: center;">
  <img src="docs/media/loan_calculator_1.png" alt="Loans Calculator Input" width="500">
  <figcaption style="font-style: italic; color: #555; margin-top: 8px;">
    Loans Calculator Input
  </figcaption>
</figure>

The loans calculator allows you to calculate the duration of a loan, and some other details, based on the following fields:

| Field | Description |
| ----------- | ----------- |
| Starting Principal | The principal of the loan at the start of the loan. |
| Monthly Payment | The monthly payments you  will make. |
| Other Payments | Total amount included in your monthly payments that doesn't go to interest or principal. For example, insurances, taxes, etc. |
| Yearly Interest Rate | The yearly interest rate (APR) of your loan. |
| Start Date | The start date of the loan. |

Clicking on **Calculate loan** will run calculations and return the following information:

<figure style="text-align: center;">
  <img src="docs/media/loan_calculator_2.png" alt="Loans Calculator Output" width="500">
  <figcaption style="font-style: italic; color: #555; margin-top: 8px;">
    Loans Calculator Output
  </figcaption>
</figure>

| Field | Description |
| ----------- | ----------- |
| Starting Principal | The principal of the loan at the start of the loan. |
| Monthly Payment | The monthly payments you  will make. |
| Other Payments | Total amount included in your monthly payments that doesn't go to interest or principal. For example, insurances, taxes, etc. |
| Yearly Interest Rate | The yearly interest rate (APR) of your loan. |
| Start Date | The start date of the loan. |

Additionally, it will generate a month-to-month status report of your loan, with the following fields:

| Field | Description |
| ----------- | ----------- |
| Date | The date of the status. |
| Payment | The payment made that month. |
| Interest | The interest charged that month. |
| Other Payments | The amount of the payment that went into other payments (not interest or principal). |
| Paydown | The amount paid directly to principal. |
| Paydown | The remaining principal at the end of the month. |

## Applying Principal Payments

Bucktracker can also help you **apply extraordinary payments** to principal, to see how much faster you will pay off the loan, how much less interest you will pay, etc.

To do this, you need to **create a user** at:

```Welcome Screen -> Log in -> Create Account```

And then log in. Once logged in, **create a new loan** with the button at the bottom right:

<figure style="text-align: center;">
  <img src="docs/media/loans_menu_1.png" alt="Loans Menu" width="500">
  <figcaption style="font-style: italic; color: #555; margin-top: 8px;">
    Loans Menu
  </figcaption>
</figure>

Enter your loan data and **create the loan**. It works the same as the calculator, with an additional name field.

After creating the loan, go to view its details. You will find a _Default Payment Plan_ on the _Payment Plans_ section. This is the normal payment plan without extraordinary payments. Add a new Payment Plan using the button:

<figure style="text-align: center;">
  <img src="docs/media/loan_details_1.png" alt="Loan Details" width="500">
  <figcaption style="font-style: italic; color: #555; margin-top: 8px;">
    Loan Details
  </figcaption>
</figure>

A payment plan consists on a list of amounts and dates, representing the extraordinary payments you will make. You can add new rows to add more payments:

<figure style="text-align: center;">
  <img src="docs/media/loan_payment_plan_creation_1.png" alt="Loans Calculator Output" width="500">
  <figcaption style="font-style: italic; color: #555; margin-top: 8px;">
    Loans Calculator Output
  </figcaption>
</figure>


After saving the new payment plan, it will appear in the _Payment Plans_ section. There you can compare the loan duration, or view the details of the Payment Plan:

<figure style="text-align: center;">
  <img src="docs/media/payment_plan_comparison.png" alt="Payment Plan Comparison" width="500">
  <figcaption style="font-style: italic; color: #555; margin-top: 8px;">
    Payment Plan Comparison
  </figcaption>
</figure>


<figure style="text-align: center;">
  <img src="docs/media/payment_plan_details.png" alt="Payment Plan Details" width="500">
  <figcaption style="font-style: italic; color: #555; margin-top: 8px;">
    Payment Plan Details
  </figcaption>
</figure>


## Contributing

### Clone the repo

```bash
git clone https://github.com/Mr-Rafael/bucktracker-api.git
cd bucktracker-api
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

# API Documentation

Endpoint request and response details live in [docs/api.md](docs/api.md).

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
