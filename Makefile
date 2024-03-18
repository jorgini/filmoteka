build:
	docker-compose build filmoteka

run:
	docker-compose up filmoteka

test:
	go test ./handlers -v -cover

migrate:
	migrate -path ./database/migration -database 'postgres://postgres:qwerty@0.0.0.0:5432/postgres?sslmode=disable' up
