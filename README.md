
# Тестовое задание:
1. Общие сведения
   Название проекта: Сервис сбора статистики
   Описание проекта: Микросервис на golang для сбора статистики
   Цель проекта: Создать сервис с возможностью хранения статистики в базе данных
2. Требования к функциональности
   Должно быть 4 api ручки:
   GetOrderBook(exchange_name, pair string) ([]*DepthOrder, error)
   SaveOrderBook(exchange_name, pair string, orderBook []*DepthOrder) error
   GetOrderHistory(client *Client)  ([]*HistoryOrder, error)
   SaveOrder(client *Client, order *HistoryOrder) error
3. Технические требования
   Язык программирования: Go 1.22
   Библиотеки и фреймворки: Любые на ваш выбор
   Архитектура: REST API
   База данных:  Лучше реализовать ClickHouse, но можно и  Postgres
4. Нефункциональные требования
   Производительность: Время отклика сервера не более 200мс
   Rps: до 200 на запись, до 100 на чтение
5. Тестирование
   Тестирование: Unit-тесты для всех основных функций
6. Требования к документации
   Документация кода: Комментарии к основным модулям и функциям
   Пользовательская документация: Руководство пользователя
7. Сроки выполнения
   Ожидаемые сроки: 1 неделя
8. Детали реализации:

Таблица OrderBook
```
id          int64
exchange    string
pair        string
asks        []depthOrder
bids        []depthOrder
```

```go
type DepthOrder struct{
Price   float64
BaseQty float64
}
```

Таблица Order_History

```
client_name              	string
exchange_name   	        string
label		                string
pair  		                string
side    		            string
type                        string
base_qty                    float64
price                       float64
algorithm_name_placed       string
lowest_sell_prc             float64
highest_buy_prc             float64
commission_quote_qty        float64
time_placed                 time.time

```

```go
type HistoryOrder struct{
client_name              	string
exchange_name   	        string
label		                string
pair  		                string
side    		        string
type                            string
base_qty                        float64
price                           float64
algorithm_name_placed           string
lowest_sell_prc                 float64
highest_buy_prc                 float64
commission_quote_qty            float64
time_placed                     time.time
}
```
```go
type Client struct{
client_name              	string
exchange_name   	        string
label		                string
pair  		                string
}
```

--- 

# Пояснения к решению

Было создано две таблицы в Clickhouse: `order_book` и `order_history`.
Для каждого POST запроса проходит валидация данных,
считается что массив depthOrder хранит в себе данные о asks и bids,
из-за чего достаточно разделить его поровну для получение этих массивов. 

---

# Установка и запуск

```
git clone https://github.com/Rpqshka/vortex-stats-collector.git
```

```
cd vortex-stats-collector
```

```
docker-compose up --build
```
```
go run .\cmd\main.go
```


ТЕСТЫ:
```
go test .\tests\ -v 
```
Вы можете использовать тестовую конфигурацию, которая находится в файле ```.env```, либо настроить сервис под себя и добавить этот файл в ```.gitignore```

---

# Работа с сервисом

После загрузки БД и запуска сервиса становятся доступны эндпоинты:
- GET    /order-book          
- PUT    /order-book           
- GET    /order-history         
- PUT    /order-history   

[![postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/24093475-af31a427-e343-465a-b073-b39d5edae9ff?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D24093475-af31a427-e343-465a-b073-b39d5edae9ff%26entityType%3Dcollection%26workspaceId%3Db652d614-668b-484c-b035-be4525c69c9f)

---

# Контакты
Telegram : @rpqshka

email: rpqshka@gmail.com