# Установка бинарей
```bash
make bin-deps
```

# Генерация кода

```bash
make generate
```

# Установка зависимостей

```bash
go mod download
```
# Запуск

### Server
```bash
go run cmd/server/server.go
```

### Client
```bash
go run cmd/client/client.go
```

# Отправка сообщений через WS

```bash
websocat ws://localhost:8081/api/v1/stream/bid
```

```json
{"id": "test"}
```

# JS code для CORS

```js
fetch('http://localhost:8081/api/v1/notes/ID')
  .then(response => {
    // Check if the request was successful (status code 200-299)
    if (!response.ok) {
      throw new Error('Network response was not ok: ' + response.statusText);
    }
    // Parse the response body as JSON
    return response.json();
  })
  .then(data => {
    // Work with the retrieved JSON data
    console.log('Data received:', data);
  })
  .catch(error => {
    // Handle any errors that occurred during the fetch operation
    console.error('There has been a problem with your fetch operation:', error);
  });
```