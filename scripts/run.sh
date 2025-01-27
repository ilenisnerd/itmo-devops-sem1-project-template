#!/bin/bash

echo "Сборка сервера..."

go build -o server ./main.go

echo "Сервер собрался"

nohup ./server > output.log 2>&1 &

echo "Сервер запустился"