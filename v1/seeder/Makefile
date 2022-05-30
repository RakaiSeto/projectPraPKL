migrateup:
	migrate -path db/migration -database "postgresql://postgres:password@localhost/prepkl?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://postgres:password@localhost/prepkl?sslmode=disable" -verbose down