run:
	ENV=dev go run main.go

local-redis:
	docker run -d --name election-redis -p 6379:6379 -v $(shell pwd)/redis/data:/bitnami/redis/data -e ALLOW_EMPTY_PASSWORD=yes bitnami/redis:latest || docker start election-redis
