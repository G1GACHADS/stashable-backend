# Stashable Backend Repository

REST API Service using Go, PostgreSQL database and Redis for caching

## **Architecture**

> NOTE: This architecture was made in consideration of time (1 month deadline). Thus, this is a minified uncle bob's clean architecture without the repository/entity layer (Business is mixed with repository).

```sh
├── api # Presentation layer
├── backend # Business layer
├── build # Dockerfiles
├── clients # Infrastructure layer
├── cmd # Application entrypoint
│   ├── api # REST API entrypoint
│   └── populate # Populate db script
├── config # Application configuration
├── core # Application core libraries
│   ├── logger
│   ├── mime
│   ├── nanoid
│   └── token
├── db # Database migration scripts
└── hack # Bash scripts
    └── db
```

## Getting Started

Run the dependencies containers (PostgreSQL, Redis, Adminer)

```bash
make compose
```

Start local development environment

```bash
make apiserver
```

Migrate database

```bash
chmod +x ./hack/db/migrate_db.sh
./hack/db/migrate_db.sh
```

Populating DB (only warehouses & categories for now)

```bash
make populate
```

#### Extras

Create migration

```bash
$ chmod +x ./hack/db/migrate_new.sh
$ ./hack/db/migrate_new.sh
Usage: ./hack/db/migrate_new.sh NAME

$ ./hack/db/migrate_new.sh init
/home/.../stashable-backend/db/000001_init.up.sql
/home/.../stashable-backend/db/000001_init.down.sql

```

HTTP Load test

```bash
chmod +x ./hack/http_load_test.sh
./hack/http_load_test.sh
```
