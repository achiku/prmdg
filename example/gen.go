package taskyapi

//go:generate prmdg struct --file=./doc/schema/schema.json --package=taskyapi --output=./struct.go
//go:generate prmdg jsval --file=./doc/schema/schema.json --package=taskyapi --output=./validator.go
