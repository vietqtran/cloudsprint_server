// Package swagger Code generated by swaggo/swag. DO NOT EDIT
package swagger

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "your.email@example.com"
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
        "/auth/forgot-password": {
            "post": {
                "description": "Send a password reset email",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Request password reset",
                "parameters": [
                    {
                        "description": "Forgot password request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.ForgotPasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.BaseResponse"
                        }
                    }
                }
            }
        },
        "/auth/github/auth": {
            "get": {
                "description": "Redirect to GitHub for authentication",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Initiate GitHub OAuth",
                "responses": {}
            }
        },
        "/auth/github/callback": {
            "get": {
                "description": "Process the callback from GitHub OAuth",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "GitHub OAuth callback",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization code",
                        "name": "code",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "State for CSRF protection",
                        "name": "state",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/auth/refresh": {
            "post": {
                "description": "Refresh access token using refresh token or session ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Refresh token",
                "parameters": [
                    {
                        "description": "Refresh token request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.RefreshTokenRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.RefreshTokenResponse"
                        }
                    }
                }
            }
        },
        "/auth/reset-password": {
            "post": {
                "description": "Reset password with a valid token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Reset password",
                "parameters": [
                    {
                        "description": "Reset password request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.ResetPasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.BaseResponse"
                        }
                    }
                }
            }
        },
        "/auth/sign-in": {
            "post": {
                "description": "SignIn with username and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "SignIn a user",
                "parameters": [
                    {
                        "description": "SignIn request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.SignInRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.SignInResponse"
                        }
                    }
                }
            }
        },
        "/auth/sign-up": {
            "post": {
                "description": "SignUp a new user with username, email, and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "SignUp a new user",
                "parameters": [
                    {
                        "description": "SignUp request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.SignUpRequest"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/auth/verify-email/send-otp": {
            "post": {
                "description": "Send a one-time password to verify email",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Send email verification OTP",
                "parameters": [
                    {
                        "description": "Send OTP request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.SendEmailOTPRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.BaseResponse"
                        }
                    }
                }
            }
        },
        "/auth/verify-email/status": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Check if the authenticated user's email is verified",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Check email verification status",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.EmailVerificationResponse"
                        }
                    }
                }
            }
        },
        "/auth/verify-email/verify": {
            "post": {
                "description": "Verify email address using OTP",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Verify email with OTP",
                "parameters": [
                    {
                        "description": "Verify OTP request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.VerifyEmailOTPRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.EmailVerificationResponse"
                        }
                    }
                }
            }
        },
        "/auth/verify-reset-token": {
            "post": {
                "description": "Verify if a password reset token is valid",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Verify reset token",
                "parameters": [
                    {
                        "description": "Verify token request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.VerifyResetTokenRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.BaseResponse"
                        }
                    }
                }
            }
        },
        "/github/repositories": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get all GitHub repositories for the authenticated user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "github"
                ],
                "summary": "List GitHub repositories",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/response.GitHubRepositoryResponse"
                            }
                        }
                    }
                }
            }
        },
        "/github/repositories/{repo_name}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Get a specific GitHub repository by name",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "github"
                ],
                "summary": "Get GitHub repository",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Repository name",
                        "name": "repo_name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.GitHubRepositoryResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "constants.ErrorCode": {
            "type": "string",
            "enum": [
                "000001",
                "000002"
            ],
            "x-enum-varnames": [
                "COMMON_ERROR",
                "EMAIL_UNVERIFIED"
            ]
        },
        "constants.HttpStatusCode": {
            "type": "integer",
            "enum": [
                100,
                101,
                102,
                103,
                200,
                201,
                202,
                203,
                204,
                205,
                206,
                207,
                208,
                226,
                300,
                301,
                302,
                303,
                304,
                305,
                307,
                308,
                400,
                401,
                402,
                403,
                404,
                405,
                406,
                407,
                408,
                409,
                410,
                411,
                412,
                413,
                414,
                415,
                416,
                417,
                418,
                421,
                422,
                423,
                424,
                425,
                426,
                428,
                429,
                431,
                451,
                500,
                501,
                502,
                503,
                504,
                505,
                506,
                507,
                508,
                510,
                511
            ],
            "x-enum-varnames": [
                "StatusContinue",
                "StatusSwitchingProtocols",
                "StatusProcessing",
                "StatusEarlyHints",
                "StatusOK",
                "StatusCreated",
                "StatusAccepted",
                "StatusNonAuthoritativeInformation",
                "StatusNoContent",
                "StatusResetContent",
                "StatusPartialContent",
                "StatusMultiStatus",
                "StatusAlreadyReported",
                "StatusIMUsed",
                "StatusMultipleChoices",
                "StatusMovedPermanently",
                "StatusFound",
                "StatusSeeOther",
                "StatusNotModified",
                "StatusUseProxy",
                "StatusTemporaryRedirect",
                "StatusPermanentRedirect",
                "StatusBadRequest",
                "StatusUnauthorized",
                "StatusPaymentRequired",
                "StatusForbidden",
                "StatusNotFound",
                "StatusMethodNotAllowed",
                "StatusNotAcceptable",
                "StatusProxyAuthRequired",
                "StatusRequestTimeout",
                "StatusConflict",
                "StatusGone",
                "StatusLengthRequired",
                "StatusPreconditionFailed",
                "StatusRequestEntityTooLarge",
                "StatusRequestURITooLong",
                "StatusUnsupportedMediaType",
                "StatusRequestedRangeNotSatisfiable",
                "StatusExpectationFailed",
                "StatusTeapot",
                "StatusMisdirectedRequest",
                "StatusUnprocessableEntity",
                "StatusLocked",
                "StatusFailedDependency",
                "StatusTooEarly",
                "StatusUpgradeRequired",
                "StatusPreconditionRequired",
                "StatusTooManyRequests",
                "StatusRequestHeaderFieldsTooLarge",
                "StatusUnavailableForLegalReasons",
                "StatusInternalServerError",
                "StatusNotImplemented",
                "StatusBadGateway",
                "StatusServiceUnavailable",
                "StatusGatewayTimeout",
                "StatusHTTPVersionNotSupported",
                "StatusVariantAlsoNegotiates",
                "StatusInsufficientStorage",
                "StatusLoopDetected",
                "StatusNotExtended",
                "StatusNetworkAuthenticationRequired"
            ]
        },
        "request.ForgotPasswordRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "request.RefreshTokenRequest": {
            "type": "object",
            "properties": {
                "sessionId": {
                    "type": "string"
                }
            }
        },
        "request.ResetPasswordRequest": {
            "type": "object",
            "properties": {
                "confirmPassword": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "request.SendEmailOTPRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "request.SignInRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "request.SignUpRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "request.VerifyEmailOTPRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "otp": {
                    "type": "string"
                }
            }
        },
        "request.VerifyResetTokenRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "response.BaseResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "$ref": "#/definitions/constants.HttpStatusCode"
                },
                "data": {},
                "error_code": {
                    "$ref": "#/definitions/constants.ErrorCode"
                },
                "message": {
                    "type": "string"
                },
                "pagination": {
                    "$ref": "#/definitions/response.Pagination"
                },
                "status": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "string"
                },
                "trace": {}
            }
        },
        "response.EmailVerificationResponse": {
            "type": "object",
            "properties": {
                "emailVerified": {
                    "type": "boolean"
                }
            }
        },
        "response.GitHubRepositoryResponse": {
            "type": "object",
            "properties": {
                "clone_url": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "fork": {
                    "type": "boolean"
                },
                "full_name": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "language": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "private": {
                    "type": "boolean"
                },
                "updated_at": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "response.Pagination": {
            "type": "object",
            "properties": {
                "page": {
                    "type": "integer"
                },
                "pages": {
                    "type": "integer"
                },
                "per_page": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "response.RefreshTokenResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                }
            }
        },
        "response.SignInResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                },
                "session_id": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/response.UserResponse"
                }
            }
        },
        "response.UserResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "first_name": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "last_name": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Go Postgres API",
	Description:      "A RESTful API built with Go, Fiber, and PostgreSQL",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
