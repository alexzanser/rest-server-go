# rest-server-go

``` 
Сервер предоставляет следующий API:

POST   /task/              :  создаёт задачу и возвращает её ID  
GET    /task/<taskid>      :  возвращает одну задачу по её ID  
GET    /task/              :  возвращает все задачи  
DELETE /task/<taskid>      :  удаляет задачу по ID  
GET    /tag/<tagname>      :  возвращает список задач с заданным тегом  
GET   /due/<yy>/<mm>/<dd> :  возвращает список задач, запланированных на указанную дату 


Запуск сервера на локальной машине:
"make build_containers".

Отправить запрос можно с помощью curl, например:
"curl "localhost:4112"

```
![alt text](my_life.pdf)​ 