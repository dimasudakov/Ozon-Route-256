# Домашнее задание 7
___
## database

Примеры использования:
- Создать соединение:
  ```go
  db, err := database.GetConnection()
  ```
- Создать таблицу:
  ```go
  err := db.CreateTable(ctx, tableName, columns)  
  ```
- Удалить базу данных
- Удалить таблицу

Работа с записями происходит через билдеры `Insert` и `Select`
### Insert
Пример:
```go
err := Insert(db).Into(tableName).Values(record).Execute()
```
### Select
Пример:
```go
result, err := Select(db).Values(selectValues).From(tableName).Where(conditions).Execute()
```
где `selectValues` - столбцы которые должны вернуться, `conditions` - условия

___  

## cache
поддерживаемые методы:
```go
type Cache interface {  
    Set(key string, value any) error  
    SetWithExpiration(key string, value any, expiration time.Duration) error  
    Get(key string) (any, error)  
    Remove(key string) bool  
}
```

### LFU  (Least Frequently Used)
пример создания:
```go
cache := cache.New().LFU().Expiration(duration).Build()
```


### LRU (Least Recently Used)
пример создания:
```go
cache := cache.New().LRU().Expiration(duration).Build()
```
___

## Тестирование

### database
```shell
go test ./pkg/database/... -cover 
```

___
### cache
```shell
go test ./pkg/cache/... -cover 
```