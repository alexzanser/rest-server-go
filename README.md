# rest-server-go

``` 
Сервер предоставляет следующий API:

POST   /task/              :  создаёт задачу и возвращает её ID  
GET    /task/<taskid>      :  возвращает одну задачу по её ID  
GET    /task/              :  возвращает все задачи  
DELETE /task/<taskid>      :  удаляет задачу по ID  
GET    /tag/<tagname>      :  возвращает список задач с заданным тегом  
GET   /due/<yy>/<mm>/<dd> :  возвращает список задач, запланированных на указанную дату 


Запуск сервера:
go run cmd/api 4112
``` 