Система работы с письмами представляет собой набор сервисов, подключенных к одной БД postgresql.

Сервисы следующие:
* mail server - почтовый сервис который можно развернуть, и отсылать или получать через него письма.
* currency demon - небольшой сервис который по крону ходит во внешний API для получения курса чежской кроны, сохраняя его в БД.
* currency server - веб-сервер который содержит API для получения информации о курсах валют за промежуток времени, а так же API для отправки и получения писем.

Все сервисы написаны на языке програмирования golang для того что бы все запустить нужно сделать следующее:
1) склонировать репозиторий
2) создать базу данных currency в постгресе
3) Изменить строку подключения в currency server/internal/app/app.go:23
4) Изменить строку подключения в currency demon/repositories/repository.go:25
5) Запустить веб - сервер cd /currency server/ && go run main.go
6) Запустить крон сервер cd /currency demon/ && go run main.go
7) Запустить почтовый сервер cd /mail server/ && go run main.go mail.ru

Для работы с API есть openAPI спецификация которую можно экспортировать в постман которая находится в currency server openapi.yaml

