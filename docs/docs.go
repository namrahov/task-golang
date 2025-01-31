// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://your.terms.of.service.url",
        "contact": {
            "name": "API Support",
            "url": "http://www.your-support-url.com",
            "email": "support@your-email.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/boards": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Creates a new board based on the provided data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Boards"
                ],
                "summary": "Create a new board",
                "parameters": [
                    {
                        "description": "Board Request Data",
                        "name": "BoardRequestDto",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.BoardRequestDto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Board created successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/boards/{id}/access": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Assign access to a specific board for a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Boards"
                ],
                "summary": "Give access to a board",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Board ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "userId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Access granted successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/boards/{userId}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieve all boards associated with a specific user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Boards"
                ],
                "summary": "Get user boards",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "userId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of user boards",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Board"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/files/delete/attachment/{attachmentFileId}": {
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Deletes an attachment file associated with a specific task",
                "tags": [
                    "Files"
                ],
                "summary": "Delete an attachment file",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Attachment File ID",
                        "name": "attachmentFileId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "File deleted successfully"
                    },
                    "400": {
                        "description": "Invalid request or file ID",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/files/download/attachment/{attachmentFileId}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Downloads an attachment file associated with a specific task",
                "tags": [
                    "Files"
                ],
                "summary": "Download an attachment file",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Attachment File ID",
                        "name": "attachmentFileId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "File downloaded successfully"
                    },
                    "400": {
                        "description": "Invalid request or file ID",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/files/get/task-image/{taskId}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieves an image file associated with a specific task",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "image/png",
                    " image/jpeg",
                    " image/webp"
                ],
                "tags": [
                    "Files"
                ],
                "summary": "Get task image",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Task ID",
                        "name": "taskId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Image file retrieved successfully",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "400": {
                        "description": "Invalid request or task ID",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Image not found",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/files/stream/task-video/{taskVideoId}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Streams a video file associated with a specific task",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "video/mp4",
                    " video/webm",
                    " video/ogg"
                ],
                "tags": [
                    "Files"
                ],
                "summary": "Stream task video",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Task Video ID",
                        "name": "taskVideoId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Video file streamed successfully",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "400": {
                        "description": "Invalid request or task ID",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Video not found",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/files/upload/attachment/{taskId}": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Uploads a file as an attachment for a specific task",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Files"
                ],
                "summary": "Upload an attachment file",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Task ID",
                        "name": "taskId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "File to upload",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "File uploaded successfully",
                        "schema": {
                            "$ref": "#/definitions/model.FileResponseDto"
                        }
                    },
                    "400": {
                        "description": "Invalid request or file",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/files/upload/task-image/{taskId}": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Uploads an image file associated with a specific task",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Files"
                ],
                "summary": "Upload an image for a task",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Task ID",
                        "name": "taskId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Image file to upload",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "File uploaded successfully",
                        "schema": {
                            "$ref": "#/definitions/model.FileResponseDto"
                        }
                    },
                    "400": {
                        "description": "Invalid request or task ID",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/files/upload/task-video/{taskId}": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Uploads a video file associated with a specific task",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Files"
                ],
                "summary": "Upload a video for a task",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Task ID",
                        "name": "taskId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Video file to upload",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "File uploaded successfully",
                        "schema": {
                            "$ref": "#/definitions/model.FileResponseDto"
                        }
                    },
                    "400": {
                        "description": "Invalid request or task ID",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/tasks/{boardId}": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Create a new task for a specific board",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tasks"
                ],
                "summary": "Create a task",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Board ID",
                        "name": "boardId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Task details",
                        "name": "task",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.TaskRequestDto"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Task created successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/tasks/{id}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieves the details of a specific task by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tasks"
                ],
                "summary": "Get task details",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Task ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Task details retrieved successfully",
                        "schema": {
                            "$ref": "#/definitions/model.TaskResponseDto"
                        }
                    },
                    "400": {
                        "description": "Invalid request or task ID",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Task not found",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/users/activate": {
            "get": {
                "description": "Activates a user account using the provided activation token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Activate user account",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Activation token",
                        "name": "token",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "User successfully activated",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Token is required",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/users/login": {
            "post": {
                "description": "Authenticates a user by validating their credentials and returns a JWT token upon success.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Authenticate user",
                "parameters": [
                    {
                        "description": "authRequestDto",
                        "name": "authRequestDto",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.AuthRequestDto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "JWT Token response",
                        "schema": {
                            "$ref": "#/definitions/model.JwtToken"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/users/logout": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Logs out the user by clearing all cookies and jwt token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Logout user",
                "responses": {
                    "204": {
                        "description": "User successfully logged out",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v1/users/register": {
            "post": {
                "description": "Registers a new user with the provided information.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "userRegistrationDto",
                        "name": "userRegistrationDto",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.UserRegistrationDto"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "User successfully registered",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request data",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.AuthRequestDto": {
            "type": "object",
            "properties": {
                "emailOrNickname": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "rememberMe": {
                    "type": "boolean"
                }
            }
        },
        "model.Board": {
            "type": "object",
            "properties": {
                "createdBy": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "users": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.User"
                    }
                }
            }
        },
        "model.BoardRequestDto": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "model.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "model.FileResponseDto": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "model.JwtToken": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "model.Permission": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "httpMethod": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "model.Priority": {
            "type": "string",
            "enum": [
                "HIGH",
                "MEDIUM",
                "LOW"
            ],
            "x-enum-varnames": [
                "High",
                "Medium",
                "Low"
            ]
        },
        "model.Role": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "permissions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Permission"
                    }
                }
            }
        },
        "model.Status": {
            "type": "string",
            "enum": [
                "NOT_STARTED",
                "IN_PROGRESS",
                "READY_FOR_TEST",
                "DONE"
            ],
            "x-enum-varnames": [
                "NotStarted",
                "InProgress",
                "ReadyForTest",
                "Done"
            ]
        },
        "model.TaskRequestDto": {
            "type": "object",
            "properties": {
                "deadline": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "priority": {
                    "$ref": "#/definitions/model.Priority"
                }
            }
        },
        "model.TaskResponseDto": {
            "type": "object",
            "properties": {
                "assignedBy": {
                    "$ref": "#/definitions/model.User"
                },
                "assignedTo": {
                    "$ref": "#/definitions/model.User"
                },
                "attachmentFileId": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "board": {
                    "$ref": "#/definitions/model.Board"
                },
                "changedBy": {
                    "$ref": "#/definitions/model.User"
                },
                "createdBy": {
                    "$ref": "#/definitions/model.User"
                },
                "deadline": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "priority": {
                    "$ref": "#/definitions/model.Priority"
                },
                "status": {
                    "$ref": "#/definitions/model.Status"
                },
                "taskImageUrl": {
                    "type": "string"
                },
                "taskVideoId": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "model.User": {
            "type": "object",
            "properties": {
                "acceptNotification": {
                    "type": "boolean"
                },
                "boards": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Board"
                    }
                },
                "description": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "fullName": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "inactivatedDate": {
                    "type": "string"
                },
                "isActive": {
                    "type": "boolean"
                },
                "phoneNumber": {
                    "type": "string"
                },
                "roles": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Role"
                    }
                },
                "userName": {
                    "type": "string"
                }
            }
        },
        "model.UserRegistrationDto": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "{swaggerHost}  // Dynamic host placeholder",
	BasePath:         "{swaggerBasePath}  // Dynamic base path placeholder",
	Schemes:          []string{"http", "https"},
	Title:            "Your API Title",
	Description:      "This is a sample API server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
