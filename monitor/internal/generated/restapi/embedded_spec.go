// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "title": "Monitor Service",
    "license": {
      "name": "MIT",
      "url": "https://opensource.org/license/mit"
    },
    "version": "v1"
  },
  "host": "localhost:8080",
  "basePath": "/api/v1",
  "paths": {
    "/monitor/process": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "monitor"
        ],
        "summary": "The method is used to process request.",
        "operationId": "process",
        "parameters": [
          {
            "description": "Information required to accpet a transaction.",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/ProcessRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Request successfully processed.",
            "schema": {
              "$ref": "#/definitions/ProcessResponse"
            }
          },
          "400": {
            "description": "Validation error.",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "403": {
            "description": "Forbidden error.",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "500": {
            "description": "Internal server error.",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "ErrorResponse": {
      "type": "object",
      "properties": {
        "code": {
          "type": "integer",
          "format": "int32"
        },
        "message": {
          "type": "string"
        }
      }
    },
    "ProcessRequest": {
      "type": "object",
      "required": [
        "from",
        "to",
        "method"
      ],
      "properties": {
        "from": {
          "type": "string"
        },
        "method": {
          "type": "string"
        },
        "payload": {
          "type": "object"
        },
        "to": {
          "type": "string"
        }
      }
    },
    "ProcessResponse": {
      "type": "object"
    }
  },
  "tags": [
    {
      "description": "Methods for monitor management.",
      "name": "monitor"
    }
  ]
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "title": "Monitor Service",
    "license": {
      "name": "MIT",
      "url": "https://opensource.org/license/mit"
    },
    "version": "v1"
  },
  "host": "localhost:8080",
  "basePath": "/api/v1",
  "paths": {
    "/monitor/process": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "monitor"
        ],
        "summary": "The method is used to process request.",
        "operationId": "process",
        "parameters": [
          {
            "description": "Information required to accpet a transaction.",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/ProcessRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Request successfully processed.",
            "schema": {
              "$ref": "#/definitions/ProcessResponse"
            }
          },
          "400": {
            "description": "Validation error.",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "403": {
            "description": "Forbidden error.",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "500": {
            "description": "Internal server error.",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "ErrorResponse": {
      "type": "object",
      "properties": {
        "code": {
          "type": "integer",
          "format": "int32"
        },
        "message": {
          "type": "string"
        }
      }
    },
    "ProcessRequest": {
      "type": "object",
      "required": [
        "from",
        "to",
        "method"
      ],
      "properties": {
        "from": {
          "type": "string"
        },
        "method": {
          "type": "string"
        },
        "payload": {
          "type": "object"
        },
        "to": {
          "type": "string"
        }
      }
    },
    "ProcessResponse": {
      "type": "object"
    }
  },
  "tags": [
    {
      "description": "Methods for monitor management.",
      "name": "monitor"
    }
  ]
}`))
}
