# Проект. Вопрос?. Ответ
Вопросом было, может ли пользоваетль не являющийся ответственным в любой организации, добавлять предложение для тендера `bids/new`, т.к явного указания небыло, то моим решением было - нет. Только пользователи являющиеся ответственными в любой организации, могут добавлять новые предложения.

# Развертывание проекта в Docker
1. Склонировать проект.
```shell
git clone https://github.com/Ki4EH/super-duper-zadanie.git
```
2. Перейти в корневую папку проекта.
3. При необходимости настроить переменные среды в файле `/service/.env`
```env
POSTGRES_USERNAME=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_HOST=db # изменить на вашу бд при запуске на локально
POSTGRES_PORT=5432
POSTGRES_DATABASE=
```
3. Выполнить команду .
```shell
docker-compose up --build
```
4. После успешного выполнения команды, и сборки сервер доступен по адресу `http://localhost:8080/`

Все необходимые таблицы и данные будут созданы автоматически при первом создании контейнера.

Необходимые комманды для создания таблиц и данных находятся в файле `service/init.sql`.

# Запуск локально
1. Склонировать проект.
```shell
git clone https://github.com/Ki4EH/super-duper-zadanie.git
```
2. Перейти в корневую папку проекта.
3. При настроить переменные среды в файле `/service/.env`
```env
POSTGRES_USERNAME=your_username
POSTGRES_PASSWORD=your_password
POSTGRES_HOST=your_host
POSTGRES_PORT=your_port
```
4. Запустить сервер.
```shell
cd service
go run main.go
```
