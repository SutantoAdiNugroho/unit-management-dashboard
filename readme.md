# Unit Management System
<p>
Unit Management Dashboard is web-based application for managing unit data (such as capsules and cabins) part of technical test of Bobobox
</p>

## Technologies
- **Backend**: REST API with Go (Gin) + MySQL
- **Frontend**: Next.js + Tailwind CSS + shadcn/ui
- **Database**: MySQL 8

## Core Features
- CRUD Unit
- Unit status with validation rules (Available, Occupied, Cleaning In Progress, Maintenance Needed)
- Pagination & filter
- Testing (backend unit test and API test)
- Docker Compose for fullstack running

## Prerequisites
<p>
Before begin, ensure you have the following installed:
</p>

- [Go 1.23.4 or higher](https://go.dev/dl/)
- [Nodejs 20 or higher](https://nodejs.org/en)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (WSL integration for Windows)
- [Make](https://www.gnu.org/software/make/) (optional, for command easily)

## Run Applications
1. Clone the repository:

   ```bash
   $ git clone https://github.com/SutantoAdiNugroho/unit-management-dashboard.git
   $ cd unit-management-dashboard
   ```

2. Running with Makefile
    - Initiate dependencies
        ```bash
        $ make init
        ```

    - Run backend functionality tests
        ```bash
        $ make test
        ```

    - Run all tests including API test
        ```bash
        $ ./e2e.sh
        ```

    - Run all aplications (frontend + backend)
        ```bash
        $ make run
        ```
        <p>
        After application is running, you can access frontend at `http://localhost:3000` in your browser, and backend will be served at `http://localhost:5000`
        </p>

3. Running without Makefile
    <p>
    If make is not available, you can run it manually:
    </p>

    - Initiate dependencies
        ```bash
        $ cd backend
        $ go mod tidy && go mod download
        $ cd ..
        $ cd frontend
        $ npm install
        ```
    - Run backend functionality tests
        ```bash
        $ cd backend
        $ go test ./pkg/... -v
        $ cd ..
        ``` 
    - Run all aplications (frontend + backend)
        ```bash
        $ docker compose down --volumes
        $ docker compose build
        $ docker compose up -d
        ``` 
        <p>
        Access frontend at `http://localhost:3000` in your browser, and backend will be served at `http://localhost:5000`
        </p>

## Testing
1. Unit test
    <p>Unit test is written for the backend package:</p>

    ```bash
    $ cd backend && go test ./pkg/... -v
    ```

2. API test
    <p>API test is available in backend/tests/api_test.go:</p>

    ```bash
    $ cd backend && go test ./tests/api_test.go -v
    ```

## Running with Docker Compose
```bash
$ docker compose up --build
```
<p> 
Available services:
</p>

- **MySQL** -> `localhost:3307`
- **Backend** -> `http://localhost:5000`
- **Frontend** -> `http://localhost:3000`

## API Documentation (Swagger)

<p>
After backend service is running properly, you can access the API documentation at `http://localhost:5000/swagger/index.html`
</p>