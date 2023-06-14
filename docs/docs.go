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
                "description": "Get new BTC Address",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Daemon"
                ],
                "summary": "Get new BTC Address",
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
        "/estimate": {
            "post": {
                "description": "Estimate inscription cost !!!Don't use this method!!!",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Inscriptions"
                ],
                "summary": "Estimate inscription cost !!!Don't use this method!!!",
                "parameters": [
                    {
                        "description": "Image URL from hosting service and number of blocks",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.EstimateRequest"
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
        "/feerate": {
            "get": {
                "description": "Get fee rate for transaction",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Fees"
                ],
                "summary": "Get fee rate for transaction",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.FeeResponse"
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
        "/inscription/image": {
            "get": {
                "description": "Get Base64 image from Inscription ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Inscriptions"
                ],
                "summary": "Get Base64 image from Inscription ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Inscription ID",
                        "name": "idInscription",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.InscriptionPicResponse"
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
                "description": "List of Inscriptions in the wallet",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Inscriptions"
                ],
                "summary": "List of Inscriptions in the wallet",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of items per page",
                        "name": "pageSize",
                        "in": "query"
                    }
                ],
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
        "/inscriptions/nsfw": {
            "get": {
                "description": "List of Inscriptions in the wallet waiting to be approved",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Inscriptions"
                ],
                "summary": "List of Inscriptions in the wallet waiting to be approved",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of items per page",
                        "name": "pageSize",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.NSFWInscriptionsResponse"
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
        "/inscriptions/nsfw/approve": {
            "get": {
                "description": "Approve Inscription from NSFW list",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Inscriptions"
                ],
                "summary": "Approve Inscription from NSFW list",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ORD id",
                        "name": "ord",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.HttpSuccess"
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
                "description": "Mint an Inscription",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Inscriptions"
                ],
                "summary": "Mint an Inscription",
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
        "/nsfw": {
            "post": {
                "description": "Test picture for NSFW content",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "NSFW"
                ],
                "summary": "Test picture for NSFW content",
                "parameters": [
                    {
                        "description": "File in base64 and filename",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.TestPicReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.TestPicResponse"
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
                "description": "Send an Inscription",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Inscriptions"
                ],
                "summary": "Send an Inscription",
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
        "/transaction/raw": {
            "get": {
                "description": "Get Raw transaction from BTC code",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Transactions"
                ],
                "summary": "Get Raw transaction from BTC code",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Transaction ID",
                        "name": "tx",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.RawTransaction"
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
                "description": "List of transactions in the BTC Core",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Transactions"
                ],
                "summary": "List of transactions in the BTC Core",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of items per page",
                        "name": "pageSize",
                        "in": "query"
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
        "models.EstimateRequest": {
            "type": "object",
            "properties": {
                "blocks": {
                    "type": "integer"
                },
                "imageUrl": {
                    "type": "string"
                }
            }
        },
        "models.FeeResponse": {
            "type": "object",
            "properties": {
                "feeRate": {
                    "type": "integer"
                },
                "hasError": {
                    "type": "boolean"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "models.HttpSuccess": {
            "type": "object",
            "properties": {
                "hasError": {
                    "type": "boolean"
                },
                "message": {
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
        "models.InscriptionPicResponse": {
            "type": "object",
            "properties": {
                "base64": {
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
        "models.NSFWInscriptionsResponse": {
            "type": "object",
            "properties": {
                "hasError": {
                    "type": "boolean"
                },
                "inscriptions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.NSFWTable"
                    }
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "models.NSFWTable": {
            "type": "object",
            "properties": {
                "approved": {
                    "type": "boolean"
                },
                "bc_address": {
                    "type": "string"
                },
                "content_link": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "link": {
                    "type": "string"
                },
                "ord_id": {
                    "type": "string"
                },
                "tx_id": {
                    "type": "string"
                }
            }
        },
        "models.NewAddressRequest": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                }
            }
        },
        "models.RawTransaction": {
            "type": "object",
            "properties": {
                "blockhash": {
                    "type": "string"
                },
                "blocktime": {
                    "type": "integer"
                },
                "confirmations": {
                    "type": "integer"
                },
                "hash": {
                    "type": "string"
                },
                "hex": {
                    "type": "string"
                },
                "locktime": {
                    "type": "integer"
                },
                "size": {
                    "type": "integer"
                },
                "time": {
                    "type": "integer"
                },
                "txid": {
                    "type": "string"
                },
                "version": {
                    "type": "integer"
                },
                "vin": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Vin"
                    }
                },
                "vout": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Vout"
                    }
                },
                "vsize": {
                    "type": "integer"
                },
                "weight": {
                    "type": "integer"
                }
            }
        },
        "models.ScriptPubKey": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "asm": {
                    "type": "string"
                },
                "desc": {
                    "type": "string"
                },
                "hex": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "models.ScriptSig": {
            "type": "object",
            "properties": {
                "asm": {
                    "type": "string"
                },
                "hex": {
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
        "models.TestPicReq": {
            "type": "object",
            "properties": {
                "base64": {
                    "type": "string"
                },
                "filename": {
                    "type": "string"
                }
            }
        },
        "models.TestPicResponse": {
            "type": "object",
            "properties": {
                "hasError": {
                    "type": "boolean"
                },
                "nsfwPicture": {
                    "type": "boolean"
                },
                "nsfwText": {
                    "type": "boolean"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "models.TxTable": {
            "type": "object",
            "properties": {
                "base64": {
                    "type": "string"
                },
                "bc_address": {
                    "type": "string"
                },
                "content_link": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "link": {
                    "type": "string"
                },
                "ord_id": {
                    "type": "string"
                },
                "tx_id": {
                    "type": "string"
                }
            }
        },
        "models.Vin": {
            "type": "object",
            "properties": {
                "scriptSig": {
                    "$ref": "#/definitions/models.ScriptSig"
                },
                "sequence": {
                    "type": "integer"
                },
                "txid": {
                    "type": "string"
                },
                "txinwitness": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "vout": {
                    "type": "integer"
                }
            }
        },
        "models.Vout": {
            "type": "object",
            "properties": {
                "n": {
                    "type": "integer"
                },
                "scriptPubKey": {
                    "$ref": "#/definitions/models.ScriptPubKey"
                },
                "value": {
                    "type": "number"
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
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
