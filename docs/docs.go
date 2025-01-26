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
        "model.JwtToken": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:9093",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "Your API Title",
	Description:      "This is a sample API server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}