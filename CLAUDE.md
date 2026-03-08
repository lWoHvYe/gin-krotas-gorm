# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Kratos-based microservices project written in Go that implements a distributed system with multiple services (user, product, order) using gRPC and HTTP protocols. The project follows Domain-Driven Design (DDD) principles with clean architecture.

## Essential Commands

### Development Setup
```bash
make init          # Install required dependencies (protoc plugins, kratos CLI, wire)
make api           # Generate API code from proto files
make config        # Generate internal proto configs
make all           # Run api, config, and generate commands
make build         # Build all services
```

### Running Services
```bash
# Build and run a specific service
go build -o ./bin/ ./cmd/product
./bin/product -conf ./configs/config.yaml

# Available services: product, order, user-go, helloworld-go
```

### Code Generation
```bash
# Generate from proto files
kratos proto server api/[service]/v1/[service].proto -t internal/[service]/service

# Generate database models (GORM)
go run cmd/gorm_gen.go

# Generate wire dependencies (after adding new providers)
cd cmd/[service] && wire
```

## Architecture

The project follows clean architecture with DDD:
- **api/** - Proto definitions and generated code
- **cmd/** - Service entry points (main.go files)
- **internal/** - Core business logic
  - **biz/** - Business logic layer (domain)
  - **data/** - Data access layer (infrastructure)
  - **service/** - Application service layer
  - **server/** - Transport layer (HTTP/gRPC)

## Key Technologies
- **Framework**: Kratos v2 (microservices framework)
- **Database**: MySQL with GORM
- **Cache**: Redis
- **Service Discovery**: ETCD
- **Distributed Transactions**: DTM
- **Authentication**: JWT
- **Authorization**: Casbin (RBAC)
- **DI**: Google Wire

## Adding a New Service

1. Create proto definition in `api/[service]/v1/[service].proto`
2. Generate code: `make api`
3. Generate service: `kratos proto server api/[service]/v1/[service].proto -t internal/[service]/service`
4. Implement biz layer in `internal/[service]/biz/`
5. Implement data layer in `internal/[service]/data/`
6. Add wire providers in each layer
7. Create main.go in `cmd/[service]/`
8. Generate wire: `cd cmd/[service] && wire`

## Configuration

Main configuration file: `configs/config.yaml`
- Database: MySQL connection string
- Redis: Cache configuration
- ETCD: Service discovery
- DTM: Distributed transaction manager

## Important Notes

- The project uses proto files as the source of truth for APIs
- All services expose both gRPC (port 9000) and HTTP (port 8000) endpoints
- Use dependency injection via Wire for clean architecture
- Database operations use GORM with generated type-safe query code
- Distributed transactions are handled via DTM with compensation logic