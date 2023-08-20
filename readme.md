# Test assignement for Soundproof

## 1. ASSIGNMENT

Overview
- Build simple restful backend with golang (Gin): 15 points
- Database: PostgreSQL: 5 points
- Build Login/Register/Profile APIs: 25 points

**Endpoints:**

- Register: POST /auth/register: 3 point --> Register api with email & password, store on User table with DB
- Login: POST /auth/login: 5 point --> User Login API, return user object & jwt auth token
- Get Profile API: GET /user/profile: 7 point --> Get User profile API, must authenticate with jwt token on middleware
- Update Profile API: PUT /user/profile: 10 point --> Update User Profile API, From signed str which from Metamask, convert it to
public Address and update it on DB, auth with jwt on middleware.

Integrate with Swagger: 5 points

Total: 50 points

## 2. PROJECT ARCHITECTURE

The application features 3-layer hexagonal architecture:

- **Domain/core** layer with two ports/interfaces: Storate interface and Service interface
- **Service layer**
- **Transport layer**

Dependencies are directed outwards and inner layers are not aware of what is happening on outer layer(s).

## 3. The project represents a RESTful API with the following endpoints:
- "/auth/register" (POST): user registration form which helps to create users with unique email addresses. No users with duplicated emails are allowed.
- "/auth/login" (POST): this is a user login form which checks login (email) and password of an user, and if successful, it returns an access and a refresh JWT tokens for operations described below;
- "/user/profile" (GET): retrieves user info in case of a valid JWT token;
- "/user/profile" (PUT): parses Metamask signature and a string signed by Metamask and saves an Etherium public address into the Postgres database.

## 4. DEFAULT ENVIRONMENT VARIABLES

GIN_MODE=debug

LOG_LEVEL=debug

SERVER_HOST=localhost

SERVER_PORT=8080

DB_DRIVER=postgres

DB_USER=root

DB_PASSWORD=postgres

DB_HOST=localhost

DB_PORT=5432

DB_TABLE=soundproof_db

ACCESS_SECRET=1234567890

REFRESH_SECRET=0987654321

All environment variables can be overwritten in Docker-compose file, Dockerfile, or during the deployment phase of the application.

## 5. TESTING

TO RUN THE SERVER AND TEST IT WITH REQUESTS YOU CAN USE THE FOLLOWING:

### Clone the repository
```bash
git clone https://github.com/PavelKaratkevich/soundproof
```

### Start the server
With Docker-compose
```bash
Docker-compose up 
```

OR

With the following set of commands:

```bash
make postgres
make createdb
make migrateup
make start-server
```

### Run Register User command
```golang
curl --header "Content-Type: application/json" \
--request POST  \
--data '{
    "first_name": "Adam",
    "last_name": "Smith",
    "password": "123456789",
    "email": "economy@ac.gov.uk"
}' http://localhost:8080/auth/register
```

### Run Login command
```golang
curl --header "Content-Type: application/json" \
--request POST  \
--data '{
    "password": "123456789",
    "email": "economy@ac.gov.uk"}' http://localhost:8080/auth/login
```

### Run Get User command: replace '{token}' with JWT token received in the previous step
```golang
curl --header "Content-Type: application/json" \
--request GET  \
-H "Authorization: Bearer {token}" \
--data '{
    "password": "123456789",
    "email": "economy@ac.gov.uk"}' http://localhost:8080/user/profile
```

### Run Update User command which passes login info, signed string from Metamask and a signature from Metamask, parses them and saved public Eth address into the database. Replace '{token}' with JWT token received in the previous step
```golang
curl --header "Content-Type: application/json" \
--request PUT  \
-H "Authorization: Bearer {token}" \
--data '{
    "password": "123456789",
    "email": "economy@ac.gov.uk",
    "signed_message": "0x54686973206973206vd79206e6577206d657373616765",
    "signature": "4cc9b720ca28cc5660d0c8b713fd5016e84f102ba242e86cd891ba17d88db8cc4a6257f0807a6fd3503df9a6f3c18d13d2de2efa88b143e079db4bf121012eac1b"}' http://localhost:8080/user/profile
```
the same can be done with Postman.
