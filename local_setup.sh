#!/bin/bash

echo "Обновление пакетов..."
sudo apt update

echo "Установка Golang..."
sudo apt install golang-go

echo "Установка базы данных Postgresql..."
sudo apt install postgresql

echo "Настройка базы данных Postgresql..."

sudo -u postgres psql <<EOF
-- Создаем базу данных
CREATE DATABASE "project-sem-1";

-- Создаем пользователя с указанным паролем
CREATE USER validator WITH PASSWORD 'val1dat0r';

-- Передаем владельцу базы данных все права доступа
ALTER DATABASE "project-sem-1" OWNER TO validator;
EOF

echo "База данных 'project-sem-1' и пользователь 'validator' успешно созданы."

echo "Запуск скрипта prepare.sh..."
bash scripts/prepare.sh

echo "Запуск скрипта run.sh..."
bash scripts/run.sh

if [ "$1" == "-t" ]; then
    echo "Выполнение тестов..."
    bash scripts/tests.sh 1
    exit 1
fi