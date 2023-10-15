# auth-service

**OAuth2 Authentication Service for Seamless User Authentication**

## Getting Started

### 1. Set Up Environment Variables

Create a `.env` file with the following configurations:

```bash
CLIENT_ID={{YOUR-CLIENT-ID}}
CLIENT_SECRET={{YOUR-CLIENT-SECRET}}
REDIRECT_URI="http://localhost:8081/v1/callback"
AUTH_URL="https://accounts.google.com/o/oauth2/auth"
TOKEN_URL="https://accounts.google.com/o/oauth2/token"
USER_INFO_URL="https://www.googleapis.com/oauth2/v3/userinfo"

DB_HOST=localhost
DB_PORT=3306
DB_USER={{YOUR-DB-USER}}
DB_PASSWORD={{YOUR-DB-PASSWORD}}
DB_NAME=auth_svc
```

### 2. Create a new database:

```sql
CREATE DATABASE auth_svc;
```

### 3. Run Database Migrations

Execute the SQL commands in `db/migrations` to set up the database schema.

### 4. Run the Application

Launch the application using the following command:

```bash
go run main.go
```

### 5. Access the Server

The server will be accessible at [http://localhost:8081](http://localhost:8081).

## Contributing

We welcome contributions! Feel free to submit issues, feature requests, or pull requests.

## License

This project is licensed under the [MIT License](LICENSE).
