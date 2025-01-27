#!/bin/bash

echo "Установка зависимостей..."
go mod tidy

echo "Подготовка базы данных PostgreSQL..."

PGPASSWORD=val1dat0r psql -U validator -h localhost -p 5432 -d project-sem-1 -c "
CREATE TABLE IF NOT EXISTS prices (
    id SERIAL PRIMARY KEY,          
    created_at DATE NOT NULL,
    name VARCHAR(255) NOT NULL,     
    category VARCHAR(255) NOT NULL,  
    price DECIMAL(10, 2) NOT NULL
);"

echo "База данных PostgreSQL подготовлена"