// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Vince",
            "url": "https://vincent.serpoul.com",
            "email": "v@po.com"
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
        "/happycatfact": {
            "get": {
                "description": "listHandler returns a list of cat fact",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "happyCatFact"
                ],
                "summary": "Get list of happy cat facts",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/db.HappycatFact"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/http.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "createHandler does not return an article",
                "tags": [
                    "HappyCat"
                ],
                "summary": "saves a happy cat fact",
                "parameters": [
                    {
                        "description": "happy cat fact",
                        "name": "happycatfact",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.SaveHappycatFactParams"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": ""
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/http.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/happycatfact/{happyCatFactID}": {
            "get": {
                "description": "getHandler returns a single cat fact by id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "happyCatFact"
                ],
                "summary": "Get happy cat fact by id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "happy cat fact id",
                        "name": "happyCatFactID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/db.HappycatFact"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/http.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "db.HappycatFact": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "fact": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                }
            }
        },
        "db.SaveHappycatFactParams": {
            "type": "object",
            "properties": {
                "fact": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                }
            }
        },
        "http.ErrorResponse": {
            "type": "object",
            "properties": {
                "error_code": {
                    "type": "string"
                },
                "error_msg": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "archiver.orchestration.dev",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "Swagger archiver API",
	Description:      "This is a sample server db save service.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
