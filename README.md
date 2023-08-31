# Test Avito Service - тестовое задание на стажера Backend Go
# Сервис динамического сегментирования пользователей

Это сервис, хранящий пользователя и сегменты, в которых он состоит (создание, изменение, удаление сегментов, а также добавление и удаление пользователей в сегмент)

**Использовались следующие технологии:**
- Golang 1.20, gin
- PostgreSQL 15
- Docker и Docker Compose

Инструкции по запуску:
есть .env файл, там указываем все параметры для того, чтобы наш сервис запустился с БД.

Пример файла:

DB_USER_NAME=postgres

DB_PASSWORD=postgres

DB_HOST=database

DB_PORT=5432

DB_NAME=avito

DB_HOST_PORT=5432

DB_DOCKER_PORT=5432

GO_HOST_PORT=8080

GO_DOCKER_PORT=8080

Опционально можно менять под определенный сервер/машину

Локально у себя я поднимал с таким .env следующей командой:

```sh
$ docker compose up --build
```

Дальше локально запросы кидать через http://localhost:$GO_HOST_PORT

в данном случае .env - http://localhost:8083

### Он имеет 5 методов API:
(далее будут даваться примеры JSON, т.е. параметры, которые нужны запросу)

1. **/create_segment** - POST метод

```sh
http://localhost:8083/create_segment
```
Запрос:
```json
{
  "segment_name": "A",
  "percents": 25
}
 ```
Ответ:
```json
{
"message": "segment is successfully created"
}
```
Метод добавляет сегмент в список текущих сегментов.
На выходе либо JSON с сообщением об успешном добавлении либо текст ошибки.

2. **/delete_segment** - DELETE метод
```sh
http://localhost:8083/delete_segment
```
Запрос:
```json
{
  "segment_name": "A"
}
 ```
Ответ:
```json
{
  "message": "Segment deleted successfully"
}
```
Если сегмент не добавляли автоматически пользователям(то есть был процент >0),
то метод удаляет метод из списка текущих сегментов, также удаляет этот сегмент у всех пользователей и в истории время удаления сегмента меняется на текущее, если он еще не был удален.
На выходе либо JSON с сообщением об успешном удалении либо текст ошибки.
3. **/get_user_segments** - GET метод

```sh
http://localhost:8083//get_user_segments
```
Запрос:
```json
{
  "user_id": 1
}
 ```
Ответ:
```json
[
  "A",
  "B"
]
```
Метод возвращает список сегментов пользователя, в которых он состоит.

4) **/change_segments** - PUT метод
```json
{
  "adding_segments": [
    {
      "segment_name": "A",
      "delete_time": "2023-09-01T19:25:00Z"
    },
		{
      "segment_name": "B",
      "delete_time": "2023-09-01T19:25:00Z"
		}
  ],
  "removing_segments": [],
  "user_id": 1
}
```
```json
{
	"message": "segments updated successfully"
}
```

5) **/get_report_csv** - GET метод

Запрос:
```json
{
  "user_id": [1],
	"year": 2023,
	"month": 9
}
```
Ответ:
придет csv в ответе 
## Наша база данных состоит из 4 таблиц:

- ### id_name_segments

| Столбец    | Тип               |
|------------|-------------------|
| segment_id | int (primary key) |
|segment_name| text              |


- ### users
| Столбец | Тип    |
|---|--------|
| user_id | bigint |

- ### users_segments
| Столбец | Тип       |
|---|-----------|
|user_id| bigint    |
|segment_id| int       |
|endtime| timestamp |

- ### user_segment_history
| Столбец    | Тип       |
|------------|-----------|
| id         | bigint (primary key)    |
| segment_id | int       |
| timestamp  | timestamp |
|operation| varchar(10)|
|user_id|bigint|
Какие были вопросы/проблемы и как решались:
1) Изначально планировалось сделать 2 таблицы users_segments и user_segment_history.
Но позже было решено добавить еще 2 таблицы - id_name_segments и users.
Таблица id_name_segments нужна нам для проверки существования сегмента, так как без нее мы бы не смогли верно контролировать сегменты. 
В целом, концепция связей между таблицами очень помогла избавиться от избыточности и путаницы в данных.
2) 
3) В ТЗ в доп задании 3 написано, что сегмент, добавленный автоматически, должен отдаваться всегда. Поэтому если нам поступит команда об удалении сегмента, мы его удалять не будем.
Что касается автоматического назначения, если мы хотим это сделать, то обязательно в JSON указываем поле percents со значением 1<=percents<=100.
Если поля percents нет или его значение 0, то сегмент не добавится автоматически никому, и его можно будет впоследствии удалить.
4) Изначально я начал хранить в истории id сегмента, но потом понял, что если сегмент удалить, то мы его название не получим, поэтому в историю добавляем сразу имя сегмента
5) В задании о реализации TTL появилась задача об удалении сегментов в определенное время.
В начале можно было реализовать все так
6) Также в целях того, чтобы наши данные целостно и взаимосвязанно добавлялись, без логических нарушений, были использованы транзакции, чтобы в случае падения какого-то запроса другой тоже откатился, и запрос вернулся с ошибкой.