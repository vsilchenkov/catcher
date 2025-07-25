// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/prj/:id/sendEvent": {
            "post": {
                "description": "Отправка Event в Sentry",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Event"
                ],
                "summary": "Send Event",
                "operationId": "sendEvent",
                "parameters": [
                    {
                        "description": "Данные события",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/catcher_app_internal_models.Event"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/catcher_app_internal_models.SendEventResult"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/app_internal_handler.errorResponse"
                        }
                    }
                }
            }
        },
        "/api/reg": {
            "get": {
                "description": "Проверка работы метода getInfo",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "info"
                ],
                "summary": "Get Info",
                "operationId": "getInfo",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/catcher_app_internal_models.RegistryInfo"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/app_internal_handler.errorResponse"
                        }
                    }
                }
            }
        },
        "/api/reg/getInfo": {
            "post": {
                "description": "Получение информации для отчета об ошибки",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Info"
                ],
                "summary": "Get Info Post",
                "operationId": "getInfoPost",
                "parameters": [
                    {
                        "description": "Значения для отчета об ошибке",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/catcher_app_internal_models.RegistryInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/catcher_app_internal_models.RegistryInfo"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/app_internal_handler.errorResponse"
                        }
                    }
                }
            }
        },
        "/api/reg/pushReport": {
            "post": {
                "description": "Отправка отчета об ошибки",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Report"
                ],
                "summary": "Push Report",
                "operationId": "pushReport",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Файл в архиве формата https://its.1c.ru/db/v8327doc#bookmark:dev:TI000002558",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/catcher_app_internal_models.RegistryPushReportResult"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/app_internal_handler.errorResponse"
                        }
                    }
                }
            }
        },
        "/api/service/clearCache": {
            "get": {
                "description": "Очищает все данные, сохранённые в кэше",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Cache"
                ],
                "summary": "Clear cache",
                "operationId": "clearCache",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/app_internal_handler.errorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "app_internal_handler.errorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "catcher_app_internal_models.Event": {
            "type": "object",
            "properties": {
                "breadcrumbs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/sentry.Breadcrumb"
                    }
                },
                "check_in": {
                    "$ref": "#/definitions/sentry.CheckIn"
                },
                "contexts": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/sentry.Context"
                    }
                },
                "debug_meta": {
                    "$ref": "#/definitions/sentry.DebugMeta"
                },
                "dist": {
                    "type": "string"
                },
                "environment": {
                    "type": "string"
                },
                "event_id": {
                    "type": "string"
                },
                "exception": {
                    "$ref": "#/definitions/catcher_app_internal_models.Exception"
                },
                "extra": {
                    "type": "object",
                    "additionalProperties": true
                },
                "fingerprint": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "level": {
                    "$ref": "#/definitions/sentry.Level"
                },
                "logger": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "modules": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "monitor_config": {
                    "$ref": "#/definitions/sentry.MonitorConfig"
                },
                "platform": {
                    "type": "string"
                },
                "release": {
                    "type": "string"
                },
                "request": {},
                "sdk": {
                    "$ref": "#/definitions/sentry.SdkInfo"
                },
                "server_name": {
                    "type": "string"
                },
                "spans": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/sentry.Span"
                    }
                },
                "start_timestamp": {
                    "type": "string"
                },
                "tags": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "threads": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/sentry.Thread"
                    }
                },
                "timestamp": {
                    "type": "string"
                },
                "transaction": {
                    "type": "string"
                },
                "transaction_info": {
                    "$ref": "#/definitions/sentry.TransactionInfo"
                },
                "type": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/sentry.User"
                }
            }
        },
        "catcher_app_internal_models.Exception": {
            "type": "object",
            "properties": {
                "values": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/catcher_app_internal_models.ExceptionValue"
                    }
                }
            }
        },
        "catcher_app_internal_models.ExceptionValue": {
            "type": "object",
            "properties": {
                "stacktrace": {
                    "$ref": "#/definitions/catcher_app_internal_models.Stacktrace"
                },
                "type": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "catcher_app_internal_models.Frame": {
            "type": "object",
            "properties": {
                "abs_path": {
                    "type": "string"
                },
                "context_line": {
                    "type": "string"
                },
                "filename": {
                    "type": "string"
                },
                "function": {
                    "type": "string"
                },
                "in_app": {
                    "type": "boolean"
                },
                "lineno": {
                    "type": "integer"
                },
                "module": {
                    "type": "string"
                },
                "module_abs": {
                    "type": "string"
                },
                "stack_start": {
                    "type": "boolean"
                }
            }
        },
        "catcher_app_internal_models.RegistryInfo": {
            "type": "object",
            "properties": {
                "dumpType": {
                    "type": "integer"
                },
                "needSendReport": {
                    "type": "boolean"
                },
                "userMessage": {
                    "type": "string"
                }
            }
        },
        "catcher_app_internal_models.RegistryInput": {
            "type": "object",
            "required": [
                "configName"
            ],
            "properties": {
                "ErrorCategories": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "appName": {
                    "type": "string"
                },
                "appStackHash": {
                    "type": "string"
                },
                "appVersion": {
                    "type": "string"
                },
                "clientID": {
                    "type": "string"
                },
                "clientStackHash": {
                    "type": "string"
                },
                "configHash": {
                    "type": "string"
                },
                "configName": {
                    "type": "string"
                },
                "configVersion": {
                    "type": "string"
                },
                "configurationInterfaceLanguageCode": {
                    "type": "string"
                },
                "platformInterfaceLanguageCode": {
                    "type": "string"
                },
                "platformType": {
                    "type": "string"
                },
                "reportID": {
                    "type": "string"
                }
            }
        },
        "catcher_app_internal_models.RegistryPushReportResult": {
            "type": "object",
            "properties": {
                "eventID": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                }
            }
        },
        "catcher_app_internal_models.SendEventResult": {
            "type": "object",
            "properties": {
                "eventID": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                }
            }
        },
        "catcher_app_internal_models.Stacktrace": {
            "type": "object",
            "properties": {
                "frames": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/catcher_app_internal_models.Frame"
                    }
                }
            }
        },
        "sentry.Breadcrumb": {
            "type": "object",
            "properties": {
                "category": {
                    "type": "string"
                },
                "data": {
                    "type": "object",
                    "additionalProperties": true
                },
                "level": {
                    "$ref": "#/definitions/sentry.Level"
                },
                "message": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "sentry.CheckIn": {
            "type": "object",
            "properties": {
                "check_in_id": {
                    "description": "Check-In ID (unique and client generated)",
                    "type": "string"
                },
                "duration": {
                    "description": "The duration of the check-in. Will only take effect if the status is ok or error.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/time.Duration"
                        }
                    ]
                },
                "monitor_slug": {
                    "description": "The distinct slug of the monitor.",
                    "type": "string"
                },
                "status": {
                    "description": "The status of the check-in.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/sentry.CheckInStatus"
                        }
                    ]
                }
            }
        },
        "sentry.CheckInStatus": {
            "type": "string",
            "enum": [
                "in_progress",
                "ok",
                "error"
            ],
            "x-enum-varnames": [
                "CheckInStatusInProgress",
                "CheckInStatusOK",
                "CheckInStatusError"
            ]
        },
        "sentry.Context": {
            "type": "object",
            "additionalProperties": true
        },
        "sentry.DebugMeta": {
            "type": "object",
            "properties": {
                "images": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/sentry.DebugMetaImage"
                    }
                },
                "sdk_info": {
                    "$ref": "#/definitions/sentry.DebugMetaSdkInfo"
                }
            }
        },
        "sentry.DebugMetaImage": {
            "type": "object",
            "properties": {
                "arch": {
                    "description": "macho,elf,pe",
                    "type": "string"
                },
                "code_file": {
                    "description": "macho,elf,pe,wasm,sourcemap",
                    "type": "string"
                },
                "code_id": {
                    "description": "macho,elf,pe,wasm",
                    "type": "string"
                },
                "debug_file": {
                    "description": "macho,elf,pe,wasm",
                    "type": "string"
                },
                "debug_id": {
                    "description": "macho,elf,pe,wasm,sourcemap",
                    "type": "string"
                },
                "image_addr": {
                    "description": "macho,elf,pe",
                    "type": "string"
                },
                "image_size": {
                    "description": "macho,elf,pe",
                    "type": "integer"
                },
                "image_vmaddr": {
                    "description": "macho,elf,pe",
                    "type": "string"
                },
                "type": {
                    "description": "all",
                    "type": "string"
                },
                "uuid": {
                    "description": "proguard",
                    "type": "string"
                }
            }
        },
        "sentry.DebugMetaSdkInfo": {
            "type": "object",
            "properties": {
                "sdk_name": {
                    "type": "string"
                },
                "version_major": {
                    "type": "integer"
                },
                "version_minor": {
                    "type": "integer"
                },
                "version_patchlevel": {
                    "type": "integer"
                }
            }
        },
        "sentry.Exception": {
            "type": "object",
            "properties": {
                "mechanism": {
                    "$ref": "#/definitions/sentry.Mechanism"
                },
                "module": {
                    "type": "string"
                },
                "stacktrace": {
                    "$ref": "#/definitions/sentry.Stacktrace"
                },
                "thread_id": {
                    "type": "integer"
                },
                "type": {
                    "description": "used as the main issue title",
                    "type": "string"
                },
                "value": {
                    "description": "used as the main issue subtitle",
                    "type": "string"
                }
            }
        },
        "sentry.Frame": {
            "type": "object",
            "properties": {
                "abs_path": {
                    "type": "string"
                },
                "addr_mode": {
                    "type": "string"
                },
                "colno": {
                    "type": "integer"
                },
                "context_line": {
                    "type": "string"
                },
                "filename": {
                    "type": "string"
                },
                "function": {
                    "type": "string"
                },
                "image_addr": {
                    "type": "string"
                },
                "in_app": {
                    "type": "boolean"
                },
                "instruction_addr": {
                    "type": "string"
                },
                "lineno": {
                    "type": "integer"
                },
                "module": {
                    "description": "Module is, despite the name, the Sentry protocol equivalent of a Go\npackage's import path.",
                    "type": "string"
                },
                "package": {
                    "description": "Package and the below are not used for Go stack trace frames.  In\nother platforms it refers to a container where the Module can be\nfound.  For example, a Java JAR, a .NET Assembly, or a native\ndynamic library.  They exists for completeness, allowing the\nconstruction and reporting of custom event payloads.",
                    "type": "string"
                },
                "platform": {
                    "type": "string"
                },
                "post_context": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "pre_context": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "stack_start": {
                    "type": "boolean"
                },
                "symbol": {
                    "type": "string"
                },
                "symbol_addr": {
                    "type": "string"
                },
                "vars": {
                    "type": "object",
                    "additionalProperties": true
                }
            }
        },
        "sentry.Level": {
            "type": "string",
            "enum": [
                "debug",
                "info",
                "warning",
                "error",
                "fatal"
            ],
            "x-enum-varnames": [
                "LevelDebug",
                "LevelInfo",
                "LevelWarning",
                "LevelError",
                "LevelFatal"
            ]
        },
        "sentry.Mechanism": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object",
                    "additionalProperties": {}
                },
                "description": {
                    "type": "string"
                },
                "exception_id": {
                    "type": "integer"
                },
                "handled": {
                    "type": "boolean"
                },
                "help_link": {
                    "type": "string"
                },
                "is_exception_group": {
                    "type": "boolean"
                },
                "parent_id": {
                    "type": "integer"
                },
                "source": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "sentry.MonitorConfig": {
            "type": "object",
            "properties": {
                "checkin_margin": {
                    "description": "The allowed margin of minutes after the expected check-in time that\nthe monitor will not be considered missed for.",
                    "type": "integer"
                },
                "failure_issue_threshold": {
                    "description": "The number of consecutive failed check-ins it takes before an issue is created.",
                    "type": "integer"
                },
                "max_runtime": {
                    "description": "The allowed duration in minutes that the monitor may be ` + "`" + `in_progress` + "`" + `\nfor before being considered failed due to timeout.",
                    "type": "integer"
                },
                "recovery_threshold": {
                    "description": "The number of consecutive OK check-ins it takes before an issue is resolved.",
                    "type": "integer"
                },
                "schedule": {},
                "timezone": {
                    "description": "A tz database string representing the timezone which the monitor's execution schedule is in.\nSee: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones",
                    "type": "string"
                }
            }
        },
        "sentry.Request": {
            "type": "object",
            "properties": {
                "cookies": {
                    "type": "string"
                },
                "data": {
                    "type": "string"
                },
                "env": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "headers": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "method": {
                    "type": "string"
                },
                "query_string": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "sentry.SdkInfo": {
            "type": "object",
            "properties": {
                "integrations": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "packages": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/sentry.SdkPackage"
                    }
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "sentry.SdkPackage": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "sentry.Span": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object",
                    "additionalProperties": true
                },
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "op": {
                    "type": "string"
                },
                "origin": {
                    "type": "string"
                },
                "parent_span_id": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "span_id": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "start_timestamp": {
                    "type": "string"
                },
                "status": {
                    "$ref": "#/definitions/sentry.SpanStatus"
                },
                "tags": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "timestamp": {
                    "type": "string"
                },
                "trace_id": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "sentry.SpanStatus": {
            "type": "integer",
            "enum": [
                0,
                1,
                2,
                3,
                4,
                5,
                6,
                7,
                8,
                9,
                10,
                11,
                12,
                13,
                14,
                15,
                16,
                17,
                18
            ],
            "x-enum-varnames": [
                "SpanStatusUndefined",
                "SpanStatusOK",
                "SpanStatusCanceled",
                "SpanStatusUnknown",
                "SpanStatusInvalidArgument",
                "SpanStatusDeadlineExceeded",
                "SpanStatusNotFound",
                "SpanStatusAlreadyExists",
                "SpanStatusPermissionDenied",
                "SpanStatusResourceExhausted",
                "SpanStatusFailedPrecondition",
                "SpanStatusAborted",
                "SpanStatusOutOfRange",
                "SpanStatusUnimplemented",
                "SpanStatusInternalError",
                "SpanStatusUnavailable",
                "SpanStatusDataLoss",
                "SpanStatusUnauthenticated",
                "maxSpanStatus"
            ]
        },
        "sentry.Stacktrace": {
            "type": "object",
            "properties": {
                "frames": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/sentry.Frame"
                    }
                },
                "frames_omitted": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "sentry.Thread": {
            "type": "object",
            "properties": {
                "crashed": {
                    "type": "boolean"
                },
                "current": {
                    "type": "boolean"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "stacktrace": {
                    "$ref": "#/definitions/sentry.Stacktrace"
                }
            }
        },
        "sentry.TransactionInfo": {
            "type": "object",
            "properties": {
                "source": {
                    "$ref": "#/definitions/sentry.TransactionSource"
                }
            }
        },
        "sentry.TransactionSource": {
            "type": "string",
            "enum": [
                "custom",
                "url",
                "route",
                "view",
                "component",
                "task"
            ],
            "x-enum-varnames": [
                "SourceCustom",
                "SourceURL",
                "SourceRoute",
                "SourceView",
                "SourceComponent",
                "SourceTask"
            ]
        },
        "sentry.User": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "ip_address": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "time.Duration": {
            "type": "integer",
            "enum": [
                -9223372036854775808,
                9223372036854775807,
                1,
                1000,
                1000000,
                1000000000,
                60000000000,
                3600000000000,
                -9223372036854775808,
                9223372036854775807,
                1,
                1000,
                1000000,
                1000000000,
                60000000000,
                3600000000000
            ],
            "x-enum-varnames": [
                "minDuration",
                "maxDuration",
                "Nanosecond",
                "Microsecond",
                "Millisecond",
                "Second",
                "Minute",
                "Hour",
                "minDuration",
                "maxDuration",
                "Nanosecond",
                "Microsecond",
                "Millisecond",
                "Second",
                "Minute",
                "Hour"
            ]
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8000",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Catcher",
	Description:      "Catcher API Service",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
