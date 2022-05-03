run:
	ENV=dev go run main.go

gqlgen:
	gqlgen generate

local-redis:
	docker run -d --name election-redis -p 6379:6379 -e ALLOW_EMPTY_PASSWORD=yes bitnami/redis:latest || docker start election-redis