# User Management Service

It is a Go-based service that provides functionalities related to user creation, user deletion, password resets, etc.

## Features

- **Integration with MongoDB:** Utilizes [MongoDB](https://github.com/mongodb/mongo) to store user details such as credentials, addresses, contact details, etc. 
- **Gin Web Framework:** Uses the [Gin](https://github.com/gin-gonic/gin) web framework for handling HTTP requests and responses.

## Prerequisites

Before running the service, make sure you have the following dependencies installed:

- Go (version 1.20 or higher)
- Docker 

## Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/user-management-service.git
   cd user-management-service

3. Add a `.env` file in the root directory with following fields:
   
   ```
   PORT=<add host port>
   SESSION_KEY=<random private session key to secure user sessions>
   FRONTEND_URL=<url of fronend service / reverse proxy>
   MONGO_URL=<url of mongo cluster on cloud>
   MONGO_DB_NAME=<name of db on mongo cloud>
   MONGO_COLLECTION_CUSTOMERS=<collection name for customer data>
   MONGO_COLLECTION_MERCHANTS=<collection name for merchant data>
   ```

5. Build the docker image for the **api-service**. Run the following command in the root directory.

   ```
   docker build -t user-management-service-api:latest .
   ```

6. Start the services using **docker-compose**. Note that it is expected your *Mongo* is hosted on cloud.

   ```
   sudo docker-compose compose build && sudo docker-compose up
   ```