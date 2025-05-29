# Migrations

Use the golang [migrate](https://github.com/golang-migrate/migrate) package to manage database migrations.

Create the database by running the docker-compose file in the root of the project. Then run the migrations:

```bash
migrate -path ./db/migrations/ -database "postgresql://smartsplit:smartsplit@localhost:5432/smartsplit?sslmode=disable" up
```