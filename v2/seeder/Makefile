migrateup:
	migrate -path db/migration -database "postgresql://postgres:password@localhost/projectpkl?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://postgres:password@localhost/projectpkl?sslmode=disable" -verbose down