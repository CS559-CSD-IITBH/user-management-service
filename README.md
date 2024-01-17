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

2. Add your personal **MongoDB** URL to the `config.toml` file in the `monstache` directory.

   ```
   mongo-url = <url for your personal mongo instance>
   ```
   
4. Add a `.env` file in the root directory. It should contain the following fields.

   ```
   PORT=<port for the api-service>
   ELASTICSEARCH_URL=<url for your personal elasticsearch instance>
   MONGODB_URL=<url for your personal mongo instance>
   MONGODB_COLLECTION=<database-name.collection-name in your mongo instance> 
   ```

5. Build the docker image for the **api-service**. Run the following command in the root directory.

   ```
   docker build -t user-management-service-api:latest .
   ```

6. Start the services using **docker-compose**. Note that it is expected your *Mongo* is hosted on cloud.

   ```
   sudo docker-compose compose build && sudo docker-compose up
   ```