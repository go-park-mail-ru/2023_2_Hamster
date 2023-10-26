// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Hamster API Support",
            "email": "dimka.komarov@bk.ru"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/auth/checkAuth": {
            "post": {
                "description": "Validate auth",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Validate Auth",
                "parameters": [
                    {
                        "description": "user info",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User status",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "400": {
                        "description": "Invalid cookie",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Server error: cookie read fail",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    }
                }
            }
        },
        "/api/auth/signin": {
            "post": {
                "description": "Login account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Sign In",
                "parameters": [
                    {
                        "description": "username \u0026\u0026 password",
                        "name": "userInput",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/http.signInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User logedin",
                        "schema": {
                            "$ref": "#/definitions/http.signUpResponse"
                        }
                    },
                    "400": {
                        "description": "Incorrect Input",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    }
                }
            }
        },
        "/api/auth/signup": {
            "post": {
                "description": "Create Account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Sign Up",
                "parameters": [
                    {
                        "description": "user info",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User Created",
                        "schema": {
                            "$ref": "#/definitions/http.signUpResponse"
                        }
                    },
                    "400": {
                        "description": "Incorrect Input",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    }
                }
            }
        },
        "/api/transaction/{userID}/all": {
            "get": {
                "description": "Get User all transaction",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Transaction"
                ],
                "summary": "Get all transaction",
                "responses": {
                    "200": {
                        "description": "Show actual budget",
                        "schema": {
                            "$ref": "#/definitions/http.Response-models_TransactionTransfer"
                        }
                    },
                    "204": {
                        "description": "Show actual accounts",
                        "schema": {
                            "$ref": "#/definitions/http.Response-string"
                        }
                    },
                    "400": {
                        "description": "Client error",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized user",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "403": {
                        "description": "Forbidden user",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    }
                }
            }
        },
        "/api/user/{userID}/": {
            "get": {
                "description": "Get user with chosen ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get User",
                "responses": {
                    "200": {
                        "description": "Show user",
                        "schema": {
                            "$ref": "#/definitions/http.Response-transfer_models_UserTransfer"
                        }
                    },
                    "400": {
                        "description": "Client error",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized user",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "403": {
                        "description": "Forbidden user",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    }
                }
            }
        },
        "/api/user/{userID}/accounts/all": {
            "get": {
                "description": "Get User accounts",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get User Accounts",
                "responses": {
                    "200": {
                        "description": "Show actual accounts",
                        "schema": {
                            "$ref": "#/definitions/http.Response-transfer_models_Account"
                        }
                    },
                    "204": {
                        "description": "Show actual accounts",
                        "schema": {
                            "$ref": "#/definitions/http.Response-string"
                        }
                    },
                    "400": {
                        "description": "Client error",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized user",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "403": {
                        "description": "Forbidden user",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    }
                }
            }
        },
        "/api/user/{userID}/actualBudget": {
            "get": {
                "description": "Get User actual budget",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get Actual Budget",
                "responses": {
                    "200": {
                        "description": "Show actual budget",
                        "schema": {
                            "$ref": "#/definitions/http.Response-transfer_models_BudgetActualResponse"
                        }
                    },
                    "400": {
                        "description": "Client error",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized user",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "403": {
                        "description": "Forbidden user",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    }
                }
            }
        },
        "/api/user/{userID}/balance": {
            "get": {
                "description": "Get User balance",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get Balance",
                "responses": {
                    "200": {
                        "description": "Show balance",
                        "schema": {
                            "$ref": "#/definitions/http.Response-transfer_models_BalanceResponse"
                        }
                    },
                    "400": {
                        "description": "Client error",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized user",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "403": {
                        "description": "Forbidden user",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    }
                }
            }
        },
        "/api/user/{userID}/feed": {
            "get": {
                "description": "Get Feed user info",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get Feed",
                "responses": {
                    "200": {
                        "description": "Show actual accounts",
                        "schema": {
                            "$ref": "#/definitions/http.Response-transfer_models_UserFeed"
                        }
                    },
                    "400": {
                        "description": "Client error",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized user",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "403": {
                        "description": "Forbidden user",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    }
                }
            }
        },
        "/api/user/{userID}/plannedBudget": {
            "get": {
                "description": "Get User planned budget",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get Planned Budget",
                "responses": {
                    "200": {
                        "description": "Show planned budget",
                        "schema": {
                            "$ref": "#/definitions/http.Response-transfer_models_BudgetPlannedResponse"
                        }
                    },
                    "400": {
                        "description": "Client error",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized user",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "403": {
                        "description": "Forbidden user",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    }
                }
            }
        },
        "/api/user/{userID}/update": {
            "put": {
                "description": "Update user info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "PUT Update",
                "parameters": [
                    {
                        "description": "user info update",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/transfer_models.UserUdate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Update user info",
                        "schema": {
                            "$ref": "#/definitions/http.Response-http_NilBody"
                        }
                    },
                    "400": {
                        "description": "Client error",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized user",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "403": {
                        "description": "Forbidden user",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    }
                }
            }
        },
        "/api/user/{userID}/updatePhoto": {
            "put": {
                "description": "Update user photo",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "PUT Update Photo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "userID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "New photo to upload",
                        "name": "upload",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Path to old photo",
                        "name": "path",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Photo updated successfully",
                        "schema": {
                            "$ref": "#/definitions/http.Response-transfer_models_PhotoUpdate"
                        }
                    },
                    "400": {
                        "description": "Client error",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized user",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "403": {
                        "description": "Forbidden user",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/http.ResponseError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "http.NilBody": {
            "type": "object"
        },
        "http.Response-http_NilBody": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/http.NilBody"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "http.Response-models_TransactionTransfer": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/models.TransactionTransfer"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "http.Response-string": {
            "type": "object",
            "properties": {
                "body": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "http.Response-transfer_models_Account": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/transfer_models.Account"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "http.Response-transfer_models_BalanceResponse": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/transfer_models.BalanceResponse"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "http.Response-transfer_models_BudgetActualResponse": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/transfer_models.BudgetActualResponse"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "http.Response-transfer_models_BudgetPlannedResponse": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/transfer_models.BudgetPlannedResponse"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "http.Response-transfer_models_PhotoUpdate": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/transfer_models.PhotoUpdate"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "http.Response-transfer_models_UserFeed": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/transfer_models.UserFeed"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "http.Response-transfer_models_UserTransfer": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/transfer_models.UserTransfer"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "http.ResponseError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "http.signInput": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "http.signUpResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.Accounts": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "number"
                },
                "id": {
                    "type": "string"
                },
                "mean_payment": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "models.TransactionTransfer": {
            "type": "object",
            "properties": {
                "account_id": {
                    "type": "string"
                },
                "category_id": {
                    "type": "string"
                },
                "date": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "is_income": {
                    "type": "boolean"
                },
                "payer": {
                    "type": "string"
                },
                "total": {
                    "type": "number"
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "avatar_url": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "planned_budget": {
                    "type": "number"
                },
                "salt": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "transfer_models.Account": {
            "type": "object",
            "properties": {
                "account": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Accounts"
                    }
                }
            }
        },
        "transfer_models.BalanceResponse": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "number"
                }
            }
        },
        "transfer_models.BudgetActualResponse": {
            "type": "object",
            "properties": {
                "actual_balance": {
                    "type": "number"
                }
            }
        },
        "transfer_models.BudgetPlannedResponse": {
            "type": "object",
            "properties": {
                "planned_balance": {
                    "type": "number"
                }
            }
        },
        "transfer_models.PhotoUpdate": {
            "type": "object",
            "properties": {
                "path": {
                    "type": "string"
                }
            }
        },
        "transfer_models.UserFeed": {
            "type": "object",
            "properties": {
                "account": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Accounts"
                    }
                },
                "actual_balance": {
                    "type": "number"
                },
                "balance": {
                    "type": "number"
                },
                "planned_balance": {
                    "type": "number"
                }
            }
        },
        "transfer_models.UserTransfer": {
            "type": "object",
            "properties": {
                "avatar_url": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "login": {
                    "type": "string"
                },
                "planned_budget": {
                    "type": "number"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "transfer_models.UserUdate": {
            "type": "object",
            "properties": {
                "planned_budget": {
                    "type": "number"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.1",
	Host:             "localhost:8090",
	BasePath:         "/user/{userID}/account/feed",
	Schemes:          []string{},
	Title:            "Hamster API",
	Description:      "Server API for Hamster Money Service Application",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
