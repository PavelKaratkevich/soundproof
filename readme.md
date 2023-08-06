This is a test assignement

The project is represented by a RESTful API with the following endpoints:
    - "/auth/register" (POST): this is a user registration form which helps creates users with unique email addresses. No users with duplicated emails are allowed.
	- "/auth/login" (POST): this is a user login form which checks login (email) and password of an user, and if successful, it returns an access and a refresh token for operations described below;
	- "/user/profile" (GET): this operation retrieves a user info provided there is a valid JWT token;
	- "/user/profile" (PUT): this operation parses Metamask signature and a string signed by Metamask and saves a Eth public address into the database.

Architecture:
    The application features 3-layer hexagonal architecture:
    - domain/core layer
    - service layer
    - transport layer
Dependencies are directed outwards and inner layers are not aware of what is happening on the outer layer(s).

TO RUN THE SERVER AND TEST IT WITH REQUESTS YOU CAN USE THE FOLLOWING:

with the help of curl:
<!-- Srart the server -->
1.1. With Docker-compose
    $ Docker-compose up // the easieast way to run it

1.2 set of commands:
    $ make postgres
    $ make createdb
    $ make migrateup
    $ make start-server

<!-- Run Register User command -->
2. curl --header "Content-Type: application/json" \
--request POST  \
--data '{
    "first_name": "Igor",
    "password": "123456789",
    "email": "abc@abcd.com",
    "last_name": "Karatkevich"

}' http://localhost:8080/auth/register

<!-- Run Login command -->
3. curl --header "Content-Type: application/json" \
--request POST  \
--data '{
    "password": "123456789",
    "email": "abc@abcd.com"}' http://localhost:8080/auth/login

<!-- Run Get User command: replace '{token}' with JWT token received in the step 3 -->
4. curl --header "Content-Type: application/json" \
--request GET  \
-H "Authorization: Bearer {token}" \
--data '{
    "password": "123456789",
    "email": "abc@abcd.com"}' http://localhost:8080/user/profile

<!-- Run Update User command which passes login info, signed string from Metamask and a signature from Metamask, parses them and saved public Eth address into the database. Replace '{token}' with JWT token received in the step 3 -->
5. curl --header "Content-Type: application/json" \
--request PUT  \
-H "Authorization: Bearer {token}" \
--data '{
    "password": "123456789",
    "email": "abc@abcd.com",
    "signed_message": "0x54686973206973206vd79206e6577206d657373616765",
    "signature": "4cc9b720ca28cc5660d0c8b713fd5016e84f102ba242e86cd891ba17d88db8cc4a6257f0807a6fd3503df9a6f3c18d13d2de2efa88b143e079db4bf121012eac1b"}' http://localhost:8080/user/profile

the same can be done with Postman.