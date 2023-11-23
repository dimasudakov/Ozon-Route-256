# Домашнее задание 5

- Интеграционные тесты находятся в пакете `tests`  
- Для работы приложения и тестирования используется одна и та же БД  
- Перед каждым интеграционным тестом при необходимости в бд вставляются записи, после выполнения эти записи удаляются

Процент покрытия юнит тестов: `coverage: 52.5% of statements`
Процент покрытия интеграционных тестов: `coverage: 76.5% of statements`

### Команда для запуска юнит тестов:
```shell  
make unit-test
```  

### Команда для запуска интеграционных тестов:
```shell  
make integration-test
```  

### Команды для работы с докером:
```shell  
make docker-up 
```  
```shell  
make docker-down 
```  
```shell  
make docker-start 
```  
```shell  
make docker-stop 
```  
```shell  
make docker-ps 
```  
```shell  
make docker-restart
```  

### Команды для создания миграций:
```shell  
make migrate-up
```  
```shell  
make migrate-down
```