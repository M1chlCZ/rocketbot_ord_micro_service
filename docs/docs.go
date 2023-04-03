// Code generated by swaggo/swag. DO NOT EDIT.

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
            "name": "Michal Žídek",
            "url": "http://app.rocketbot.pro",
            "email": "m1chlcz18@gmail.com"
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
        "/address": {
            "get": {
                "description": "Mint Inscription",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Daemon"
                ],
                "summary": "Mint Inscription",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.NewAddressRequest"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorHTTP"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorHTTP"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorHTTP"
                        }
                    }
                }
            }
        },
        "/inscriptions": {
            "get": {
                "description": "List Inscription",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Inscriptions"
                ],
                "summary": "List Inscription",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ListInscriptionsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorHTTP"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorHTTP"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorHTTP"
                        }
                    }
                }
            }
        },
        "/mint": {
            "post": {
                "description": "Mint Inscription",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Inscriptions"
                ],
                "summary": "Mint Inscription",
                "parameters": [
                    {
                        "description": "File in base64 and file type",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.MintRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Inscribe"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorHTTP"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorHTTP"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorHTTP"
                        }
                    }
                }
            }
        },
        "/send": {
            "post": {
                "description": "Send Inscription",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Inscriptions"
                ],
                "summary": "Send Inscription",
                "parameters": [
                    {
                        "description": "File in base64 and file type",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SendRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Inscribe"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorHTTP"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorHTTP"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorHTTP"
                        }
                    }
                }
            }
        },
        "/transactions": {
            "get": {
                "description": "List transactions from BTC Core",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Transactions"
                ],
                "summary": "List transactions from BTC Core",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Number of items per page",
                        "name": "pageSize",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "object",
                                "properties": {
                                    "abandoned": {
                                        "type": "boolean"
                                    },
                                    "address": {
                                        "type": "string"
                                    },
                                    "amount": {
                                        "type": "number"
                                    },
                                    "bip125-replaceable": {
                                        "type": "string"
                                    },
                                    "blockhash": {
                                        "type": "string"
                                    },
                                    "blockheight": {
                                        "type": "integer"
                                    },
                                    "blockindex": {
                                        "type": "integer"
                                    },
                                    "blocktime": {
                                        "type": "integer"
                                    },
                                    "category": {
                                        "type": "string"
                                    },
                                    "confirmations": {
                                        "type": "integer"
                                    },
                                    "fee": {
                                        "type": "number"
                                    },
                                    "label": {
                                        "type": "string"
                                    },
                                    "parent_descs": {
                                        "type": "array",
                                        "items": {
                                            "type": "string"
                                        }
                                    },
                                    "time": {
                                        "type": "integer"
                                    },
                                    "timereceived": {
                                        "type": "integer"
                                    },
                                    "trusted": {
                                        "type": "boolean"
                                    },
                                    "txid": {
                                        "type": "string"
                                    },
                                    "vout": {
                                        "type": "integer"
                                    },
                                    "walletconflicts": {
                                        "type": "array",
                                        "items": {}
                                    },
                                    "wtxid": {
                                        "type": "string"
                                    }
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorHTTP"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorHTTP"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorHTTP"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.ErrorHTTP": {
            "type": "object",
            "properties": {
                "errorMessage": {
                    "type": "string"
                },
                "hasError": {
                    "type": "boolean"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "models.Inscribe": {
            "type": "object",
            "properties": {
                "commit": {
                    "type": "string"
                },
                "fees": {
                    "type": "integer"
                },
                "inscription": {
                    "type": "string"
                },
                "reveal": {
                    "type": "string"
                }
            }
        },
        "models.ListInscriptionsResponse": {
            "type": "object",
            "properties": {
                "hasError": {
                    "type": "boolean"
                },
                "inscriptions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.TxTable"
                    }
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "models.MintRequest": {
            "type": "object",
            "properties": {
                "base64": {
                    "type": "string"
                },
                "feeRate": {
                    "type": "integer"
                },
                "format": {
                    "type": "string"
                }
            }
        },
        "models.NewAddressRequest": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "privKey": {
                    "type": "string"
                }
            }
        },
        "models.SendRequest": {
            "type": "object",
            "properties": {
                "Address": {
                    "type": "string"
                },
                "InscriptionID": {
                    "type": "string"
                },
                "feeRate": {
                    "type": "integer"
                }
            }
        },
        "models.TxTable": {
            "type": "object",
            "properties": {
                "bcAddress": {
                    "type": "string"
                },
                "contentLink": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "link": {
                    "type": "string"
                },
                "ordID": {
                    "type": "string"
                },
                "txID": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "89.116.25.234:7500",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Rocketbot ORD API",
	Description:      "Private API for ORD",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
