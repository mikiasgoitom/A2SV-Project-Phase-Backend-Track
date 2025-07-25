# ✅ Task 6 - Task Management API

This document provides an overview of the Task Management API, including its base URL, endpoints, request and response format, error handling, status codes, example cURL commands, technologies used, and versioning.

---

## 🌐 Base URL

```
localhost:8080/tasks
```

---

## 📌 Endpoints

### 1. `GET /` – Get All Tasks

- **Description:** Retrieves a list of all items in the task management system.
- **Response:**

```json
[
  {
    "id": 1,
    "title": "Sample Task",
    "description": "This is a task.",
    "due_date": "2025-08-01",
    "status": "pending"
  }
]
```

- **Status Codes:**

  - `200 OK`: Successfully retrieved the list of items.
  - `404 Not Found`: No tasks found.

---

### 2. `GET /tasks/{id}` – Get Task by ID

- **Description:** Retrieves a specific task by ID.
- **Response:**

```json
{
  "id": 1,
  "title": "Sample Task",
  "description": "This is a task.",
  "due_date": "2025-08-01",
  "status": "pending"
}
```

- **Status Codes:**

  - `200 OK`: Successfully retrieved the task.
  - `404 Not Found`: Task with the specified ID does not exist.

---

### 3. `PUT /{id}` – Update Task

- **Description:** Updates an existing task by ID.
- **Request:**

```json
{
  "title": "Updated Task",
  "description": "Updated description",
  "due_date": "2025-08-15",
  "status": "completed"
}
```

- **Response:**

```json
{
  "id": 1,
  "title": "Updated Task",
  "description": "Updated description",
  "due_date": "2025-08-15",
  "status": "completed"
}
```

- **Status Codes:**

  - `200 OK`: Successfully updated the task.
  - `400 Bad Request`: Invalid request format or missing required fields.
  - `404 Not Found`: Task with the specified ID does not exist.

---

### 4. `POST /` – Create a New Task

- **Description:** Creates a new task in the task management system.
- **Request:**

```json
{
  "title": "New Task",
  "description": "This is a new task.",
  "due_date": "2025-08-20",
  "status": "pending"
}
```

- **Response:**

```json
{
  "id": 2,
  "title": "New Task",
  "description": "This is a new task.",
  "due_date": "2025-08-20",
  "status": "pending"
}
```

- **Status Codes:**

  - `201 Created`: Successfully created the task.
  - `400 Bad Request`: Invalid request format or missing required fields.

---

### 5. `DELETE /{id}` – Delete Task by ID

- **Description:** Deletes a specific task by ID.
- **Response:**

```json
{
  "message": "Task deleted successfully"
}
```

- **Status Codes:**

  - `200 OK`: Successfully deleted the task.
  - `404 Not Found`: Task with the specified ID does not exist.

---

## 🧾 Request & Response Format

### ✅ Request Format

- `Content-Type: application/json`
- Body: A JSON object containing the task details (for POST and PUT)

### ✅ Response Format

- `Content-Type: application/json`
- Body: A JSON object or array containing task details or a success message

---

## ⚠️ Error Handling

The API handles errors by returning appropriate HTTP status codes and error messages in the response body. Common responses include:

- `400 Bad Request`: Invalid request format or missing required fields.
- `404 Not Found`: Task with the specified ID does not exist.

---

## 📊 Status Codes Summary

- `200 OK`: Successful request
- `201 Created`: Resource successfully created
- `204 No Content`: Successful request with no content (e.g., on deletion)
- `400 Bad Request`: Invalid input or malformed request
- `404 Not Found`: Resource not found

---

## 💻 Example cURL Commands

### Get All Tasks

```sh
curl -X GET http://localhost:8080/tasks/
```

### Get Task by ID

```sh
curl -X GET http://localhost:8080/tasks/1
```

### Create a New Task

```sh
curl -X POST http://localhost:8080/tasks/ \
    -H "Content-Type: application/json" \
    -d '{
        "title": "New Task",
        "description": "This is a new task.",
        "due_date": "2025-08-20",
        "status": "pending"
    }'
```

### Update a Task

```sh
curl -X PUT http://localhost:8080/tasks/1 \
    -H "Content-Type: application/json" \
    -d '{
        "title": "Updated Task",
        "description": "Updated description",
        "due_date": "2025-08-15",
        "status": "completed"
    }'
```

### Delete a Task

```sh
curl -X DELETE http://localhost:8080/tasks/1
```

### 1. `POST /register` – Register User

- **Description**: Creates a new user account with a unique username and password.
- **Request**:

```json
{
  "username": "new_user",
  "password": "strong_password",
  "user_type": "USER"
}
```

- **Response**:

```json
{
  "message": "User registered successfully",
  "user_id": "abc123"
}
```

- **Status Codes**:

  - `201 Created`: User successfully registered.

  - `400 Bad Request`: Validation error or username already taken.

### 2. POST /login – Login User

- **Description**: Authenticates a user and returns an access token and refresh token.
- **Request**:

```json
{
  "username": "example_user",
  "password": "example_password"
}
```

- **Response**:

```json
{
  "access_token": "jwt-access-token",
  "refresh_token": "jwt-refresh-token",
  "user_type": "USER"
}
```

- **Status Codes**:

  - `200 OK`: Login successful.

  - `400 Bad Request`: Invalid credentials or missing fields.

### 3. GET /users – Get All Users

- **Description**: Retrieves a list of all users. Only accessible by users with the ADMIN role.

- **Authorization**: Bearer Token with ADMIN role required.

- **Response**:

```json
[
  {
    "id": "abc123",
    "username": "admin_user",
    "user_type": "ADMIN",
    "created_at": "2025-07-21T10:00:00Z"
  },
  {
    "id": "def456",
    "username": "regular_user",
    "user_type": "USER",
    "created_at": "2025-07-20T10:00:00Z"
  }
]
```

- **Status Codes**:

  - `200 OK`: Successfully retrieved users.

  - `403 Forbidden`: User lacks required permissions.
  - `401 Unauthorized`: Missing or invalid token.

### 4. GET /user/{id} – Get User by ID

- **Description**: Retrieves user information by their unique ID.
- **Response**:

```json
{
  "id": "abc123",
  "username": "sample_user",
  "user_type": "USER",
  "created_at": "2025-07-21T10:00:00Z"
}
```

- **Status Codes**:
  - `200 OK`: Successfully retrieved user.
  - `404 Not Found`: User with the specified ID does not exist.

### 5. PUT /user/{id} – Update User

- **Description**: Updates a user's information. Accessible to the user themselves or an admin.
- **Authorization**: Bearer Token with USER or ADMIN role required.

- **Request**:

```json
{
  "username": "updated_user",
  "password": "new_password"
}
```

- **Response**:

```json
{
  "message": "User updated successfully"
}
```

- **Status Codes**:

  - `200 OK`: Successfully updated.
  - `400 Bad Request`: Invalid input.
  - `403 Forbidden`: Insufficient permissions.
  - `404 Not Found`: User not found.

## 6. DELETE /user/{id} – Delete User

- **Description**: Deletes a user by ID. Accessible to the user themselves or an admin.
- **Authorization**: Bearer Token with USER or ADMIN role required.
- **Response**:

```json
{
  "message": "User deleted successfully"
}
```

- **Status Codes**:

  - `200 OK`: Successfully deleted user.
  - `403 Forbidden`: Insufficient permissions.
  - `404 Not Found`: User not found.

---

## 🛠 Technologies Used

- **Go**: Programming language
- **Gin**: Web framework for routing and middleware

---

## 🧩 Versioning

The API follows [Semantic Versioning](https://semver.org/).
Current version: `v1.0.0`
Future updates will be versioned accordingly.

---
