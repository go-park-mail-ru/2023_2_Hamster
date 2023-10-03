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
                            "$ref": "#/definitions/http.loginResponse"
                        }
                    },
                    "400": {
                        "description": "Incorrect Input",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
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
                            "$ref": "#/definitions/http.Error"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
                        }
                    }
                }
            }
        },
        "/api/auth/validateAuth": {
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
                            "$ref": "#/definitions/http.Response"
                        }
                    },
                    "400": {
                        "description": "Invalid cookie",
                        "schema": {
                            "$ref": "#/definitions/http.Response"
                        }
                    },
                    "500": {
                        "description": "Server error: cookie read fail",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
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
                            "$ref": "#/definitions/http.Response"
                        }
                    },
                    "400": {
                        "description": "Client error",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
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
                            "$ref": "#/definitions/http.budgetActualResponse"
                        }
                    },
                    "400": {
                        "description": "Client error",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
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
                            "$ref": "#/definitions/http.Response"
                        }
                    },
                    "400": {
                        "description": "Client error",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
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
                            "$ref": "#/definitions/http.budgetPlannedResponse"
                        }
                    },
                    "400": {
                        "description": "Client error",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
                        }
                    },
                    "500": {
                        "description": "Server error",
                        "schema": {
                            "$ref": "#/definitions/http.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "http.Error": {
            "type": "object",
            "properties": {
                "errmsg": {
                    "type": "string"
                }
            }
        },
        "http.Response": {
            "type": "object",
            "properties": {
                "body": {},
                "status": {
                    "type": "string"
                }
            }
        },
        "http.budgetActualResponse": {
            "type": "object",
            "properties": {
                "actual_balance": {
                    "type": "number"
                }
            }
        },
        "http.budgetPlannedResponse": {
            "type": "object",
            "properties": {
                "planned_balance": {
                    "type": "number"
                }
            }
        },
        "http.loginResponse": {
            "type": "object",
            "properties": {
                "jwt": {
                    "type": "string"
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
