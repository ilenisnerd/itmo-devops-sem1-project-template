# Финальный проект 1 семестра. Простой уровень сложности.

REST API сервис для загрузки и выгрузки данных о ценах.

## Требования к системе

Написать требования к ОС или аппаратному обеспечению (если есть)

1. Linux Ubuntu 24.04
2. 2Гб Оперативной памяти и 6 Гб Объем диска
## Установка и запуск


##### 1. Скачайте проект с помощью команды
 - "git clone https://github.com/ilenisnerd/itmo-devops-sem1-project-template.git"
##### 2. Введите команду
 - "cd itmo-devops-sem1-project-template/"
##### 3. Введите команду 
- "bash setup_script.sh" (если хотите запускать без тестирования)
- "bash setup_script.sh -t" (если хотите запускать с тестированием)
## Тестирование

- Если вы запустили "bash setup_script.sh", то введите "bash ./scripts/tests.sh 1"
- Если вы не запустили "bash setup_script.sh", то введите "bash setup_script.sh -t"
- Если вы ввели "bash setup_script.sh -t", то вы молодец✔ и уже запустили тестирование!

##### Какие тесты проходит приложение?
Тесты проверяют корректность: 

POST запроса:
- curl -s -F "file=@test_data.zip" "http://localhost:8080/api/v0/prices"

GET запроса:
- curl -s "http://localhost:8080/api/v0/prices" -o "response.zip"

Работы БД PostgreSQL:
- PGPASSWORD=val1dat0r psql -h localhost -U validator -d project-sem-1 -c '\q' 2>/dev/null

## Контакт

К кому можно обращаться в случае вопросов?
 - https://t.me/dimitri_kopylov