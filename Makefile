migrations-up:
	go run ./cmd/migrator/main.go \
	--env-path=.env \
	--migrations-path=migrations

migrations-down:
	go run ./cmd/migrator/main.go \
	--action=down \
	--env-path=./.env \
	--migrations-path=./migrations

run:
	docker compose --env-file=docker.env up -d --build

test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -count 5
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html coverage.out -o coverage.html
	go tool cover -func=./coverage.out | grep "total"

e2e-stand:
	docker compose -f=./docker-compose-e2e.yaml --env-file=e2e.env up -d
	go run ./cmd/migrator/main.go --env-path=e2e.env --migrations-path=migrations
	go run ./cmd/api --env-path=e2e.env