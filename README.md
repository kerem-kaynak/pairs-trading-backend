# Web Application to Visualize the Performance of Machine Learning Aided Pairs Trading Strategies

![Archutecture Diagram](https://i.ibb.co/F4GWF1g/Pairs-Trading-Architecture.png)

This repository is part of a project that aims to showcase the performance of a machine learning assisted pairs trading strategy developed for a thesis. The project consists of 4 separate services in different repositories.

- [Orchestrator](https://github.com/kerem-kaynak/pairs-trading-orchestrator): Orchestrator of data pipelines and processing workflows. Runs scheduled ETL jobs.
- [Quant / ML Service](https://github.com/kerem-kaynak/pairs-trading-quant-service): Web server exposing endpoints to perform machine learning tasks.
- [Backend API](https://github.com/kerem-kaynak/pairs-trading-backend): Backend API serving data to the client.
- [Frontend](https://github.com/kerem-kaynak/pairs-trading-frontend): Frontend application for web access.

The research in the thesis leading to this project can be found [here](https://github.com/kerem-kaynak/pairs-trading-with-ml) with deeper explanations of the financial and statistical concepts.

## Pairs Trading Backend API

This service is the bridge between the data and the end user. The API provides useful endpoints to interface with the vast amounts of data in the database efficiently. It interacts with the frontend, the database and the ML service for certain endpoints.

# Technologies

The service is built using Golang Gin. It also utilizes GORM for models, migrations and database interactions. Utilizes OAuth flow for Google Sign In. The whole service is containerized and can be run locally easily with the help of a Makefile.

# Project Structure

Project is made up of two main directories. The `cmd` directory is the entrypoint for the API. The `internal` folder is where everything else lies. There are multiple folders for certain abstractions such as auth, config, database, handlers, http and models.

- `auth` handles OAuth and authentication middleware.
- `config` loads certain configurations for the initialization of the service.
- `database` handles database connections and migrations.
- `handlers` defines and generates route handlers and contains all the business logic.
- `http` handles the initialization of the webserver and the router configuration.
- `models` keeps interfaces for all the data models the backend interacts with.

The whole project is Dockerized and has a docker-compose file for ease of development.

# Requirements

- Go 1.21+
- Docker / Docker Compose
- Make

# Local Development

Install dependencies:
```
go mod tidy
```

Create & populate a .env file:
```
PORT=
DB_USER=
DB_PASS=
DB_NAME=
INSTANCE_CONNECTION_NAME=
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
JWT_SECRET=
GOOGLE_REDIRECT_URL_HOST=
POLYGON_API_KEY=
QUANT_SERVICE_HOST=
QUANT_SERVICE_API_KEY=
```

Build and run the binary locally:
```
make run
```

Optionally, run using Docker:
```
make docker-compose-up
```
