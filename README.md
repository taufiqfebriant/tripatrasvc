# Tripatra Backend Service

Mohon maaf, karena keterbatasan waktu, autentikasi saat ini hanya mencakup fitur login. Idealnya, harus ada mekanisme refresh token menggunakan httpOnly cookie dan blacklist token saat logout.

## Technologies

- **Go**: The primary programming language for the backend.
- **MongoDB**: The database used for storing data.
- **GraphQL**: The query language used for API communication.
- **Echo**: The web framework used for building the API.
- **JWT**: JSON Web Tokens for authentication.
- **Docker**: For containerization and easy setup.

## Setup and Run

1. Clone the repository

   ```bash
   git clone https://github.com/taufiqfebriant/tripatrasvc.git
   ```

2. Go to the project directory

   ```bash
   cd tripatrasvc
   ```

3. Setup environment variables

   ```bash
   cp .env.example .env
   ```

4. Spin up the development environment using Docker Compose. Wait until you see "Echo. http server started on ..."

   ```bash
   docker-compose up --build
   ```

5. Open new terminal and seed the database with sample data

   ```bash
   docker exec tripatrasvc go run scripts/seed.go
   ```

6. Access the GraphQL Playground at [http://localhost:1323](http://localhost:1323)
