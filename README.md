# Отборочное задание Авито Golang

## Запуск проекта

1. Клонировать репозиторий
2. Запустить `docker-compose up`
3. Сервис будет доступен на `http://localhost:8080`

## API методы

После запуска проекта по ручке http://localhost:8080/swagger-ui можно посмотреть документацию к API.

### 1. Создание сегмента

**Метод**: `POST`
- **URL**: `/v1/segments`
- **Тело запроса**:
    - `slug` (обязательный): Название сегмента
    - `percent`: Процент пользователей для автоматического добавления этого сегмента
- **Ответы**:
    - `201 Created`
    - `400 Bad Request`
    - `500 Internal Server Error`

#### Пример тела запроса (Postman)

`POST http://localhost:8080/api/v1/segments`

```json
{
  "slug": "AVITO_VOICE_MESSAGES",
  "percent": 50
}
```

#### Пример тела ответа

```json
{
    "message": "success"
}
```

### 2. Удаление сегмента

**Метод**: `DELETE`
- **URL**: `/v1/segments`
- **Параметры запроса**:
    - `slug` (обязательный): Название сегмента для удаления
- **Ответы**:
    - `200 OK`
    - `400 Bad Request`
    - `500 Internal Server Error`

#### Пример тела запроса (Postman)

`DELETE http://localhost:8080/api/v1/segments`

```json
{
  "slug": "AVITO_VOICE_MESSAGES"
}
```

#### Пример тела ответа

```json
{
    "message": "success"
}
```

### 3. Добавление или удаление сегментов для пользователя

**Метод**: `PATCH`
- **URL**: `/v1/users`
- **Тело запроса**: JSON объект:
    - `to_create`: Массив объектов:
        - `slug` (обязательный): Название сегмента
        - `ttl`: Дата и время, до которых сегмент будет активен для пользователя, в формате ISO 8601 ("2006-01-02T15:04:
          05Z").
    - `to_delete`: Массив названий сегментов, которые нужно удалить.
    - `user_id` (обязательный)
- **Ответы**:
    - `200 OK`
    - `400 Bad Request`
    - `500 Internal Server Error`

#### Пример тела запроса (Postman)

`PATCH http://localhost:8080/api/v1/users`

```json
{
  "to_create": [
    {
      "slug": "AVITO_VOICE_MESSAGES",
      "ttl": "2023-12-31T23:59:59Z"
    }
  ],
  "to_delete": [
    "AVITO_DISCOUNT_30"
  ],
  "user_id": 1000
}
```

#### Пример тела ответа

```json
{
  "message": "success"
}
```

### 4. Получение активных сегментов пользователя

**Метод**: `GET`
- **URL**: `/v1/users`
- **Параметры запроса**:
    - `user_id` (обязательный)
- **Ответы**:
- `200 OK`:
  - `segments`: Массив названий и id сегментов.
- `400 Bad Request`
- `500 Internal Server Error`

#### Пример тела запроса (Postman)

`GET http://localhost:8080/api/v1/users?user_id=1`

#### Пример тела ответа

```json
{
  "segments": [
    {
      "id": 1,
      "slug": "AVITO_VOICE_MESSAGES"
    },
    {
      "id": 2,
      "slug": "AVITO_DISCOUNT_30"
    }
  ]
}
```

### 5. Создание логов пользователя

**Метод**: `POST`
- **URL**: `/v1/users/logs`
- **Тело запроса**:
  - `user_id` (обязательный): ID пользователя
  - `date` (обязательный): Месяц и год в формате `YYYY-MM`
- **Ответы**:
  - `201 Created`Возвращает путь к файлу лога в контейнере. Название: id пользователя + дата
  - `400 Bad Request`
  - `500 Internal Server Error`

#### Пример тела запроса (Postman)

`POST http://localhost:8080/api/v1/users/logs`

```json
{
    "user_id": 1,
    "date": "2023-08"
}
```

#### Пример тела ответа

```json
{
    "file_path": "/tmp/1_2023-08.csv"
}
```

## Дополнительные задания

1. **Логирование**

   В решении данной задачи в базе данных создана отдельная таблица `logs`, которая представлена структурой `Log`:

    ```go
      type Log struct {
	        gorm.Model
	        ID        uint64     `gorm:"primaryKey;autoIncrement"`
	        UserID    uint64     `gorm:"not null"`
	        Segment   string     `gorm:"not null"`
	        EventType string     `gorm:"not null"`
	        Time      *time.Time `gorm:"not null"`
	  }
    ```
 
   При операции добавления/удаления пользователя в сегмент (в том числе и из п.3 доп.задания и 
   автоматического удаления устаревших сегментов из п.2),
   в таблицу с логами производиться запись. 
   
   При вызове метода `/logs` данные за соответствующий месяц
   записываются в `{id пользователя}_{дата}.csv` файл в контейнере по пути `/tmp`.

2. **TTL**

    В таблицу отношения многие-ко-многим между пользователями и сегментами было добавлено поле `ExpirationDate`, которое указывает дату и время, когда данное отношение должно быть автоматически удалено.

    ```go
    type UserSegment struct {
	    gorm.Model

	    UserID         uint64
	    SegmentID      uint64
	    ExpirationDate time.Time `gorm:"default:null"`
    }
    ```

    В `main` запускается отдельная горутина `TtlWorker` и каждую минуту проверяет, есть ли отношения, которые должны быть удалены на основе поля `ExpirationDate`.

    ```go
    func TtlWorker(repo segment.IRepository) {
	    for {
		    now := time.Now().Add(time.Hour * 3)
		    next := now.Add(time.Minute).Truncate(time.Minute)
		    time.Sleep(next.Sub(now))

		    log.Println("Deleting expired segments")
		    DeleteExpiredSegments(repo, now)
	    }
    }
    ```

    `TtlWorker` вызывает функцию `DeleteExpiredSegments`, которая в свою очередь обращается к репозиторию для удаления устаревших записей. После успешного удаления, в таблицу `logs` добавляются соответствующие записи.

3. **Автоматическое добавление сегментов**
    
    В методе API для создания сегмента (`POST api/v1/segments`) был добавлен параметр `percent`, 
    позволяющий задавать процент пользователей, которые будут автоматически добавлены 
    в этот сегмент.
    Считается общее количество пользователей и вычисляется число пользователей, которые будут добавлены в сегмент на основе заданного процента
    ```go
    numUsersToAssign := int(float64(totalUsers) * float64(percent) / 100)
    ```
   Из общего списка пользователей случайным образом выбирается заданное количество, к которому добавляется сегмент с последующей записью логов.

## Вопросы и решения

1. **Формат даты**

2. **Обработка ошибок**