# Task Board API

## Overview
This project provides a task management system that allows users to create tasks on boards, assign them to users, and manage files associated with tasks. Users can be granted access to boards, retrieve their boards, and upload/download various types of files. The system uses JWT authentication stored in Redis, preventing unauthorized access if the token is not found in Redis. It also integrates RabbitMQ for queue processing and provides API documentation using Swagger.

## Technologies Used
- **Golang** (Backend Development)
- **PostgreSQL** (Database)
- **Redis** (Token Storage)
- **MinIO** (File Storage)
- **RabbitMQ** (Message Queue)
- **JWT Token Authentication** (Security)
- **Swagger** (API Documentation)
- **Mux Router** (HTTP Routing)

## Features
- **User Authentication & Management**
    - User registration with activation email
    - JWT-based authentication stored in Redis
    - Logout functionality that removes JWT from Redis
    - Account deletion with a 30-day reactivation period

- **Task Management**
    - Create tasks under specific boards
    - Retrieve details of a task

- **Board Management**
    - Create boards
    - Grant access to boards for specific users
    - Retrieve boards assigned to a user

- **File Management**
    - Upload and download task-related attachments
    - Upload and retrieve task images
    - Upload and stream task videos

- **Pagination Support**
- **Inter-service Communication using RabbitMQ**

## API Endpoints

### User Endpoints
| Method | Endpoint | Description |
|--------|---------|-------------|
| POST | `/users/login` | Authenticate a user |
| POST | `/users/register` | Register a new user |
| GET | `/users/active` | Activate a user account |
| GET | `/users/logout` | Logout and remove JWT from Redis |

### Task Endpoints
| Method | Endpoint | Description |
|--------|---------|-------------|
| POST | `/tasks/{boardId}` | Create a task under a board |
| GET | `/tasks/{id}` | Retrieve task details |

### Board Endpoints
| Method | Endpoint | Description |
|--------|---------|-------------|
| POST | `/boards` | Create a board |
| POST | `/boards/{id}/access` | Grant access to a board |
| GET | `/boards/{userId}` | Get boards assigned to a user |

### File Management Endpoints
| Method | Endpoint | Description |
|--------|---------|-------------|
| POST | `/files/upload/attachment/{taskId}` | Upload an attachment file |
| DELETE | `/files/delete/attachment/{attachmentFileId}` | Delete an attachment file |
| GET | `/files/download/attachment/{attachmentFileId}` | Download an attachment file |
| POST | `/files/upload/task-image/{taskId}` | Upload a task image |
| GET | `/files/get/task-image/{taskId}` | Retrieve a task image |
| GET | `/files/stream/task-video/{taskVideoId}` | Stream a task video |
| POST | `/files/upload/task-video/{taskId}` | Upload a task video |

## Installation & Setup
### Prerequisites
- **Go** installed on your system
- **Docker & Docker Compose** (for PostgreSQL, Redis, MinIO, and RabbitMQ services)
- **Postman or Curl** (for API testing)

### Running the Project
1. Clone the repository:
   ```sh
   git clone <repository-url>
   cd task-board-api
   ```

2. Start services using Docker Compose:
   ```sh
   docker-compose up -d
   ```

3. Run the application:
   ```sh
   go run main.go
   ```

## Authentication
- All endpoints (except registration and activation) require authentication.
- A JWT token is issued upon successful login.
- The token must be included in the `Authorization` header as `Bearer <token>`.
- The token is stored in Redis for session management.
- If the token is not found in Redis, requests are denied.

## File Storage (MinIO)
- Files uploaded are stored in MinIO.
- Attachments, images, and videos related to tasks are managed via MinIO storage.

## Message Queue (RabbitMQ)
- RabbitMQ is used to listen to queues for background task processing.

## Swagger Documentation
- API documentation is available at:
  ```
  http://localhost:1991/swagger-ui/
  ```

## Notes
- If a user deletes their account and does not reactivate it within 30 days, it is permanently removed.
- Pagination is supported for relevant endpoints.
- The application supports inter-service function calls for seamless data retrieval.

## License
This project is open-source and licensed under the MIT License.

