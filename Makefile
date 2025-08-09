codegen:
	oapi-codegen -generate chi-server -package api api/schema.yaml > internal/generated/server.gen.go
	oapi-codegen -generate types -package api api/schema.yaml > internal/generated/models.gen.go