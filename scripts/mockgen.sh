#!/bin/bash
mockgen -source=internal/repository/repository.go -destination=internal/repository/mock/repository.go -package mock Repository
mockgen -source=internal/database/database.go -destination=internal/database/mock/database.go -package mock Database