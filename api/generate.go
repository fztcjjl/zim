package api

//go:generate protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative protocol/*.proto
//go:generate protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative logic/*.proto
