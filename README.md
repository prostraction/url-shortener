# urlshort
# Docker контейнер для создания и хранения коротких URL

### Порты, используемые в приложении

Данное решение использует по умолчанию следующие порты (задаются в docker-compose.yml):
- `:50001` порт для сервера gRPC
- `:50002` порт для сервера gRPC Gateway (HTTP)

### Методы и вызовы:

Для gRPC:
- `call ToShortLink (Value)`: принимает полный URL, сохраняет полный и короткий URL, возвращает короткий URL
- `call ToFullLink (Value)`: принимает короткий URL, возвращает полный URL

Для gRPC Gateway:
- `POST /post`: принимает полный URL в теле запроса (`{"Value"="test-link"}`), возвращает короткий URL
- `GET /get/{Value}`: где {Value} = короткий URL. Запрос возвращает полный URL

### Аргументы, принимаемые программой:

Сервер (cmd/server/main.go) принимает аргумент:
- `method=db` (по умолчанию, хранение URL в БД) или `method=memory` (хранение URL в мапе)

# Компиляция и запуск контейнера:
```
git clone https://github.com/prostraction/urlshort
cd urlshort
docker-compose run urlshort --method=memory
```

# Примеры использования:

### Для получения коротких URL: 

- curl (gRPC Gateway):

```
curl --request POST "http://127.0.0.1:50002/post" -d '{"Value": "test-full-link"}'
```

- client.go (gRPC):

```
go run cmd/client/main.go toShort test-full-link
```

### Для получения полных URL:

- curl (gRPC Gateway):

```
curl --request GET "http://127.0.0.1:50002/get/test-short-link"
```

- client.go (gRPC):

```
go run cmd/client/main.go toFull test-short-link
```
