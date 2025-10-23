# KataGenzi API Backend

## Project Description

This is the API backend for the KataGenzi article platform. Built with Go (Golang) and the Fiber framework, this project provides a set of RESTful endpoints to manage users, articles (posts), tags, and image uploads. The architecture is designed to be decoupled, allowing a frontend or any other client application to interact with it independently.

## Tech Stack

- **Language:** Go (Golang)
- **Web Framework:** Fiber
- **Database ORM:** GORM
- **Database:** PostgreSQL
- **Authentication:** JSON Web Tokens (JWT)
- **Image Storage:** Cloudinary
- **Environment Variables:** godotenv

## Project Structure

```
backend/
├── api/              # Vercel serverless function entry point
├── database/         # Database configuration and connection (GORM)
├── handlers/         # Logic for handling API requests (Controllers)
├── middleware/       # Middleware for requests (e.g., JWT authentication)
├── models/           # Database table representations (GORM structs)
├── utils/            # Utility functions (e.g., Cloudinary)
├── go.mod            # Go dependency management
├── go.sum            # Dependency checksums
├── main.go           # Application entry point and route definitions (for local dev)
└── .env.example      # Example file for environment variables
```

## API Endpoints

### Authentication
- `POST /api/register`: Register a new user.
- `POST /api/login`: Log in a user and receive a JWT.

### Posts
- `GET /api/posts`: Get a paginated list of all published posts.
- `GET /api/posts/my`: Get posts belonging to the authenticated user (protected).
- `GET /api/posts/:id`: Get a single post by its ID.
- `POST /api/posts`: Create a new post (protected).
- `PUT /api/posts/:id`: Update an existing post (protected).
- `DELETE /api/posts/:id`: Move a post to trash (soft delete) (protected).

### Admin
- `GET /api/admin/posts`: Get all posts with any status (admin, protected).

### User
- `GET /api/profile`: Get the profile of the authenticated user (protected).

### Media
- `POST /api/upload`: Upload an image to Cloudinary (protected).


## Installation and Setup Guide

### 1. Prerequisites
Ensure the following software is installed on your machine:
- Go (v1.21+ recommended)
- Git
- A PostgreSQL Server

### 2. Initial Setup
1.  Clone this repository:
    ```bash
    git clone https://github.com/mohamadsolkhannawawi/article-project.git
    ```
2.  Navigate into the backend project directory:
    ```bash
    cd article-project/backend
    ```

### 3. Backend Configuration (Go)
1.  Install Go dependencies:
    ```bash
    go mod tidy
    ```
2.  Create the environment file by copying the example:
    ```bash
    cp .env.example .env
    ```
3.  Open the `.env` file and configure your database connection and other credentials.

    ```env
    # --- DATABASE ---
    # PostgreSQL connection string
    DATABASE_URL="host=localhost user=your_user password=your_password dbname=katagenzi_db port=5432 sslmode=disable"

    # --- JWT ---
    JWT_SECRET="your_super_secret_key"

    # --- CLOUDINARY ---
    CLOUDINARY_CLOUD_NAME="your_cloud_name"
    CLOUDINARY_API_KEY="your_api_key"
    CLOUDINARY_API_SECRET="your_api_secret"
    ```

4.  **Database Setup:**
    -   Start your PostgreSQL server.
    -   Create a new database with the name you specified in the `DATABASE_URL` (e.g., `katagenzi_db`).

5.  Run the backend development server:
    ```bash
    go run main.go
    ```
    The application will start, and it will automatically run database migrations on the first launch. The API will be available at `http://localhost:3000`.

### 4. Accessing the Application
-   Once the backend is running, you can start using it with a frontend application or an API testing tool like Postman.
-   For a guide on testing with Postman, please refer to the [ https://.postman.co/workspace/E-Commerce~8de014ec-8d43-44ac-be07-c0c5e10c2d87/collection/36177362-f12254c5-a0b0-4e38-935b-ee83d4457c97?action=share&creator=36177362&active-environment=36177362-22361039-cb13-40ea-ae31-706bb958fc85 ]