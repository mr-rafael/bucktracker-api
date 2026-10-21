[<- Back to README.md](/README.md#configure-the-app)

# App Configurations

## Environment Variables

| Variable | What it is for |
|---|---|
| `ALLOWED_ORIGIN` | Used for CORS. Set it to `http://your-frontend-domain.com:3000` |
| `POSTGRES_CONNECTION_STRING` | Used for database connections. Set it to `postgres://<your user>:<your password>@localhost:<database port>/finance_calculator?sslmode=disable`. |
| `ACCESS_SECRET` | The secret used for signing Access Tokens. |
| `REFRESH_SECRET` | The secret used for signing Refresh Tokens. |
| `ENV` | In what environment you're running. It affects your Refresh Token cookie. If in "production", it will be set as `Secure: true`. |
