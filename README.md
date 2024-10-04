# Онлайн библиотека песен
Проект является RestAPI сервером, предоставляющим функционал CRUD, для работы с песнями.
В нем представлены взаимодействие с БД (в нашем случае Postgres) и работа со сторонним API.

Приложение позволяет создавать, обновлять, удалять, получать песни, а так же извлекать текст конкретной песни.
Реализована фильтрация и пагинация для получения песен, а так же пагинация для извлечения текста песни.

## Как работать с приложением (запросы, которые поддерживает):
1.`http://localhost:8080/api/song`- запрос поддерживает такие HTTP методы, как: POST, PUT, DELETE.  
Но работа с этими методами будет немного отличаться.

Пример запроса (для методов POST, DELETE):  
`http://localhost:8080/api/song?group=nirvana&song=smells like teen spirit` - параметры group и song являются обязательными
Метод POST добавляет песню в БД, а DELETE удаляет.

Пример запроса (для метода PUT):  
`http://localhost:8080/api/song?group=nirvana&song=smells like teen spirit` - параметры group и song являются обязательными  
Метод PUT позволяет обновить данные песни (только поля с названиями исполнителя и песни)  
В случае использования метода следует предоставить в теле запроса json с новой информацией о песне (это тоже обязательно):

```bash
{
    "group": "Imagine Dragons",
    "song": "Radioactive"
}
```

2.`http://localhost:8080/api/song/text` - получение текста песни, запрос поддерживает только HTTP метод Get.

Пример запроса:  
`http://localhost:8080/api/song/text?group=nirvana&song=smells like teen spirit&offset=0&limit=4` - параметры group, song, offset и limit являются обязательными. Offset и limit это значения, которые будут использоваться для пагинации текста песни по куплетам

3.`http://localhost:8080/api/song/songs` - получение данных библиотеки (песен), запрос поддерживает только HTTP метод Get.
Этот запрос поддерживает только HTTP метод: GET

Пример запроса:  
`http://localhost:8080/api/songs?group=nirvana&song=smells like teen spirit&releaseDate=25.08.2009&text=here we are&link=https://www.youtube.com/watch?v=hTWKbfoikeg&offset=0&limit=4` - параметры offset и limit являются обязательными, все остальные нет. Offset и limit это значения, которые будут использоваться для пагинации кол-ва получаемых песен.  
Остальные параметры (необязательные):
* releaseDate - дата релиза песни
* text - ключевые слова для поиска по тексту песен
* link - ссылка на песню


## Информация для разработчиков:
Для запуска у себя приложния у вас должен быть конфигурационный файл .env (лежать в той же директории, что и go.mod), а так же открыто соединение с PostgreSQL (в нем уже должна быть создана БД, котоую вы укажете в конфигурационном файле .env)

## Конфигурационный файл должен иметь следующий вид:
```bash
# Данные для работы с БД:
DB_HOST=<your_host> # example: localhost
DB_PORT=<your_port> # example: 5432
DB_USER=<your_user_name>
DB_PASSWORD=<your_password>
DB_NAME=<your_db_name> # example: MusicLibrary
DRIVER_NAME=postgres
TABLE_NAME=<your_table_name> # example: music

# Данные по порту, на котором будет работать сервер
BIND_ADDR=<your_port> # example: 8080
```
