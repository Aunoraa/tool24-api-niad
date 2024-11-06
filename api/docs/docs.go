// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

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
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
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
        "/todo": {
            "get": {
                "description": "Lấy danh sách tất cả các Todo",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Todos"
                ],
                "summary": "Lấy tất cả các Todo",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.Todo"
                            }
                        }
                    }
                }
            }
        },
        "/todo/create": {
            "post": {
                "description": "Tạo một Todo mới với thông tin được cung cấp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Todos"
                ],
                "summary": "Tạo một Todo mới",
                "parameters": [
                    {
                        "description": "Thông tin của Todo",
                        "name": "todo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.Todo"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/main.Todo"
                        }
                    },
                    "400": {
                        "description": "Request body không hợp lệ",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/todo/delete/{id}": {
            "delete": {
                "description": "Xóa một Todo theo ID",
                "tags": [
                    "Todos"
                ],
                "summary": "Xóa một Todo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID của Todo",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Xóa Todo thành công",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "ID không hợp lệ",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Không tìm thấy Todo",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/todo/getuser/{id}": {
            "get": {
                "description": "Lấy thông tin chi tiết của một Todo theo ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Todos"
                ],
                "summary": "Lấy một Todo theo ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID của Todo",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Todo"
                        }
                    },
                    "400": {
                        "description": "ID không hợp lệ",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Không tìm thấy Todo",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/todo/update-status/{id}": {
            "patch": {
                "description": "Cập nhật trạng thái của một Todo theo ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Todos"
                ],
                "summary": "Cập nhật trạng thái của một Todo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID của Todo",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Todo"
                        }
                    },
                    "404": {
                        "description": "Không tìm thấy Todo",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/todo/update/{id}": {
            "patch": {
                "description": "Cập nhật thông tin của một Todo theo ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Todos"
                ],
                "summary": "Cập nhật một Todo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID của Todo",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Thông tin cập nhật của Todo",
                        "name": "todo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.Todo"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Todo"
                        }
                    },
                    "400": {
                        "description": "Request body không hợp lệ",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Không tìm thấy Todo",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.Todo": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "desc": {
                    "type": "string"
                },
                "done": {
                    "type": "boolean"
                },
                "done_at": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "title": {
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
	BasePath:         "/todo",
	Schemes:          []string{},
	Title:            "Todo API",
	Description:      "This is a sample API for managing todos.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}