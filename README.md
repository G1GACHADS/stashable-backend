# Stashable Backend Repository

Stashable REST API service in Go.

## Architecture

```sh
├── build
├── cmd
│   ├── api # Rest API endpoint
│   └── populate # Database populate script
├── db # Database Migrations
├── hack # Scripts
├── internal # Application source code
│   ├── api # Presentation layer (HTTP handlers)
│   ├── backend # Business layer
│   ├── clients # Infrastructure layer
│   └── config # Configurations
├── public
│   └── uploads # Uploaded images
├── logger
├── nanoid
└── token
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