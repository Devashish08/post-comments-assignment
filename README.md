# Post-Comments Service

A simple RESTful API service built with Go (Golang) that allows users to create text-based posts and comment on them. This project was created as a backend coding assignment.

## Features

- **Create and View Posts**: Endpoints to create, retrieve a single post, and list all posts.
- **Commenting System**: Add comments to any existing post and view all comments for a post.
- **Rich Text Support (Bonus)**: Comments are submitted in Markdown and the server converts and stores them as HTML, making them ready for any frontend client.
- **RESTful API Design**: Clean, predictable API following REST principles.
- **In-Memory Storage**: Runs out-of-the-box with no database dependency.

---

## Architecture & Technology Choices

- **Language**: **Go (Golang)**
  - Chosen for its performance, simplicity, strong standard library, and excellent support for building concurrent web services.

- **Router**: **`chi` (`github.com/go-chi/chi`)**
  - A lightweight, idiomatic, and highly performant router that simplifies handling URL parameters, middleware, and route grouping.

- **Data Storage**: **In-Memory Store**
  - For this assignment, an in-memory store (using maps protected by a `sync.RWMutex`) was chosen for its simplicity. It requires **zero external setup**, allowing the application to be run instantly with a single command. The storage logic is abstracted behind a `Store` interface, so this could easily be swapped for a persistent database like PostgreSQL or SQLite without changing the API handlers.

- **Rich Text / Markdown**: **`goldmark` (`github.com/yuin/goldmark`)**
  - A fast and extensible Markdown parser used to implement the bonus rich text feature. It securely converts user-submitted Markdown into HTML on the server.

- **UUIDs**: **`google/uuid`**
  - Used for generating unique, non-sequential IDs for posts and comments.

---

## Setup and Running the Application

### Prerequisites

- Go (version 1.18 or higher) installed on your system.
- `curl` or an API client like Postman to test the endpoints.

### Instructions

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/Devashish08/post-comments-assignment.git
    cd post-comments-assignment
    ```

2.  **Install dependencies:**
    Go modules will handle this automatically when you run the application, but you can also run this command explicitly:
    ```bash
    go mod tidy
    ```

3.  **Run the server:**
    ```bash
    go run ./cmd/api/main.go
    ```

The server will start and listen on `http://localhost:8080`.

### Health Check

You can verify the server is running by accessing the health endpoint:
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"status":"service is up running"}
```

---

## API Documentation

All request and response bodies are in JSON format.

### Posts

#### 1. Create a Post

- **Endpoint**: `POST /posts`
- **Description**: Creates a new post.
- **Request Body**:
  ```json
  {
    "content": "This is the content of my first post!"
  }
  ```
- **Example `curl`**:
  ```bash
  curl -X POST http://localhost:8080/posts -H "Content-Type: application/json" -d '{"content": "This is my first post!"}'
  ```
- **Success Response (`201 Created`)**:
  ```json
  {
    "id": "a1b2c3d4-...",
    "content": "This is my first post!",
    "created_at": "2023-10-27T10:00:00Z"
  }
  ```

#### 2. Get All Posts

- **Endpoint**: `GET /posts`
- **Example `curl`**: 
  ```bash
  curl http://localhost:8080/posts
  ```
- **Success Response (`200 OK`)**:
  ```json
  [
    {
      "id": "a1b2c3d4-...",
      "content": "This is my first post!",
      "created_at": "2023-10-27T10:00:00Z"
    }
  ]
  ```

#### 3. Get a Single Post (with Comments)

- **Endpoint**: `GET /posts/{postId}`
- **Example `curl`**: 
  ```bash
  curl http://localhost:8080/posts/a1b2c3d4-...
  ```
- **Success Response (`200 OK`)**:
  ```json
  {
    "id": "a1b2c3d4-...",
    "content": "This is my first post!",
    "created_at": "2023-10-27T10:00:00Z",
    "comments": [
      {
        "id": "e5f6g7h8-...",
        "post_id": "a1b2c3d4-...",
        "content": "**Awesome** post!",
        "content_html": "<p><strong>Awesome</strong> post!</p>\n",
        "created_at": "2023-10-27T10:05:00Z"
      }
    ]
  }
  ```

### Comments

#### 1. Create a Comment

- **Endpoint**: `POST /posts/{postId}/comments`
- **Description**: Adds a comment to the post specified by `{postId}`. Accepts Markdown in the content field.
- **Request Body**:
  ```json
  {
    "content": "This is a comment with **Markdown**."
  }
  ```
- **Example `curl`**:
  ```bash
  curl -X POST http://localhost:8080/posts/a1b2c3d4-.../comments \
    -H "Content-Type: application/json" \
    -d '{"content": "This is a comment with **Markdown**."}'
  ```
- **Success Response (`201 Created`)**:
  ```json
  {
      "id": "e5f6g7h8-...",
      "post_id": "a1b2c3d4-...",
      "content": "This is a comment with **Markdown**.",
      "content_html": "<p>This is a comment with <strong>Markdown</strong>.</p>\n",
      "created_at": "2023-10-27T10:05:00Z"
  }
  ```

#### 2. Get All Comments for a Post

- **Endpoint**: `GET /posts/{postId}/comments`
- **Example `curl`**: 
  ```bash
  curl http://localhost:8080/posts/a1b2c3d4-.../comments
  ```
- **Success Response (`200 OK`)**: Returns an array of comment objects, same format as above.

---

## Project Structure

```
post-comments-assignment/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── handler/                 # HTTP request handlers
│   │   ├── comment_handler.go   # Comment-related endpoints
│   │   ├── post_handler.go      # Post-related endpoints
│   │   └── utils.go             # Utility functions for handlers
│   ├── model/                   # Data models
│   │   ├── comment.go           # Comment struct definition
│   │   └── post.go              # Post struct definition
│   └── store/                   # Data persistence layer
│       ├── store.go             # Store interface definition
│       └── in_memory_store.go   # In-memory store implementation
├── go.mod                       # Go module definition
├── go.sum                       # Go dependencies checksum
└── README.md                    # This file
```

---

## Error Responses

All error responses follow this format:
```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- `400 Bad Request`: Invalid request body or missing required fields
- `404 Not Found`: Requested resource (post/comment) doesn't exist
- `500 Internal Server Error`: Server-side error occurred

---

## Development Notes

- **Thread Safety**: The in-memory store uses `sync.RWMutex` for concurrent access
- **UUID Generation**: Uses Google's UUID library for unique identifiers
- **Markdown Processing**: Comments support full Markdown syntax via goldmark
- **Middleware**: Includes logging and panic recovery middleware
- **No Authentication**: This is a demo service without authentication/authorization

---
