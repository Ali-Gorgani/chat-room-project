# Build the Go app with Linux and 64-bit architecture
build_binary:
	env GOOS=linux CGO_ENABLED=0 go build -o bin/auth-service cmd/main.go
	chmod +x bin/auth-service

# Run the service
run_binary:
	./auth-service

# Clean up the binary
clean_binary:
	rm -f auth-service

swagger:
	swag init -g cmd/main.go -d .

generate_initialization:
	@[ -n "$(name)" ] || (echo "Error: name variable is required. Usage: make generate_initialization name=Auth"; exit 1)
	go run -mod=mod entgo.io/ent/cmd/ent new "$(name)" --target ./utils/ent/schema

generate_migration:
	@[ -n "$(name)" ] || (echo "Error: name variable is required. Usage: make generate_migration name=my_custom_migration"; exit 1)
	atlas migrate diff "$(name)" \
		--dir "file://utils/ent/migrate/migrations" \
		--to "ent://utils/ent/schema" \
		--dev-url "docker://postgres/15/test?search_path=public"

generate_go_ent:
	go generate ./utils/ent

apply_migration:
	atlas migrate apply \
  		--dir "file://utils/ent/migrate/migrations" \
  		--url "postgres://root:secret@localhost:5433/auth-db?search_path=public&sslmode=disable"

status_migration:
	atlas migrate status \
		--dir "file://utils/ent/migrate/migrations" \
		--url "postgres://root:secret@localhost:5433/auth-db?search_path=public&sslmode=disable"

rollback_migration:
	atlas migrate down \
		--dir "file://utils/ent/migrate/migrations" \
		--url "postgres://root:secret@localhost:5433/auth-db?search_path=public&sslmode=disable" \
		--dev-url "docker://postgres/15/test?search_path=public"

proto-auth:
	@protoc \
		--proto_path=grpc/pb "grpc/pb/auth.proto" \
		--go_out=grpc/pkg/auth --go_opt=paths=source_relative \
		--go-grpc_out=grpc/pkg/auth --go-grpc_opt=paths=source_relative

.PHONY: swagger generate_initialization generate_migration generate_go_ent apply_migration status_migration rollback_migration proto-auth