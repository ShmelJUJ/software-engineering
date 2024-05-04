package user

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --package user --target gen --clean docs/yaml/api.yaml
//go:generate mockgen -package mocks -destination mocks/user_client_mocks.go github.com/ShmelJUJ/software-engineering/user/gen Handler
