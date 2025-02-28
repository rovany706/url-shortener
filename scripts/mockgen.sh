#!/bin/bash
mockgen -source=internal/repository/repository.go -destination=internal/repository/mock/repository.go -package mock Repository
mockgen -source=internal/database/database.go -destination=internal/database/mock/database.go -package mock Database
mockgen -source=internal/auth/jwt.go -destination=internal/auth/mock/jwt.go -package mock JWTAuthentication
mockgen -source=internal/service/delete_service.go -destination=internal/service/mock/delete_service.go -package mock DeleteSerice