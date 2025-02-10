DSN=postgres://postgres:mysecretpassword@localhost:5432/avito_tech_backend?sslmode=disable
PORT=:4000
SECRET_KEY=avito-backend-spring-secret-key


run:
	@echo "Building the application..."
	go run ./cmd/web --secret-key=$(SECRET_KEY) --port=$(PORT) --dsn=$(DSN)

test:
	@echo "Building the application..."
	go test ./cmd/web -coverprofile=.coverage.html
