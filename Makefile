postgres:
	 docker run --name likesDB -p 5432:5432 -d -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root postgres:15-alpine
	
redis:
	docker run -d --name redis-stack -p 6379:6379 -p 8001:8001 redis/redis-stack:latest

createdb: 
	docker exec -it likesDB createdb --username=root --owner=root likesApi 

dropdb:
	docker exec -it likesDB dropdb likesApi

migrateup:
	migrate -path db/migration -database 'postgresql://root:root@localhost:5432/likesApi?sslmode=disable' -verbose up
	
migratedown:
	migrate -path db/migration -database 'postgresql://root:root@localhost:5432/likesApi?sslmode=disable' -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test
