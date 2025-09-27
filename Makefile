.PHONY: test test_api init run

init:
	cd backend && go mod tidy && go mod download

run:
	docker compose down --volumes
	docker compose build
	docker compose up -d

logs:
	docker compose logs -f

test:
	cd backend && go test ./pkg/... -v

test_api:
	cd backend && go test ./tests/api_test.go -v