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
