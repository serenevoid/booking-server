# A dockerized Booking API server

 - Includes Golang API docker image + Postgres DB docker image
 - Both containers are linked with Docker compose
 - Includes Postman Collection file to test which requires base_url

## Requirements
 - Docker should be installed on the device
 - Twilio account with SMS service

## Set up server
1. Clone the repository
2. Fill .env file with the three necessary env variables
3. Go inside the repository where docker-compose.yml exists
4. Run `docker compose up -d` to start both API and DB containers
