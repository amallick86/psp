images:
	sudo docker pull postgres
container:
	sudo docker run --name psps -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root123 -d postgres
createdb:
	sudo docker exec -it psps createdb --username=root --owner=root psp
dropdb:
	sudo docker exec -it psps dropdb psp
migrateup:
	migrate -path db/migration -database "postgresql://root:root123@localhost:5432/psp?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration -database "postgresql://root:root123@localhost:5432/psp?sslmode=disable" -verbose up 1
migratedown:
	migrate -path db/migration -database "postgresql://root:root123@localhost:5432/psp?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration -database "postgresql://root:root123@localhost:5432/psp?sslmode=disable" -verbose down 1
migrationcreate:
	migrate create -ext sql -dir db/migration -seq add_users
sqlcinit:
	sqlc init
sqlc:
	sqlc generate	
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/amallick86/psp/db/sqlc Store

.PHONY: images createcontainer createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlcinit sqlc test server mockgen migrationcreate