## Домашнее задание 3

**Предметная область** - банковские счета и подписки зарегистрированные на них, например счет в сбербанке и подписка на Яндекс музыку привязанная к этому счету

Конфигурация для бд лежит в /configs/config.yaml

### Эндпоинты:
- **Создать аккаунт**
  ```bash
  curl.exe -X POST -H "Content-Type:application/json" -d '{"holder_name":"Dima Sudakov","balance":1000,"bank_name":"Sberbank"}' http://localhost:9000/bank-accounts
  ```
- **Получить аккаунт по ID**
  ```bash
  curl.exe -X GET http://localhost:9000/bank-accounts/3  
  ```
- **Обновить аккаунт**
  ```bash
  curl.exe -X PUT -H "Content-Type:application/json" -d '{"holder_name":"Dima Sudakov","balance":5000,"bank_name":"Sberbank"}' http://localhost:9000/bank-accounts/1
  ```
- **Удалить аккаунт**
  ```bash
  curl.exe -X DELETE http://localhost:9000/bank-accounts/1
  ```
- **Добавить подписку**
  ```bash
  curl.exe -X POST -H "Content-Type:application/json" -d '{"subscription_name": "Яндекс музыка","price": 299,"account_id": 3}' http://localhost:9000/subscriptions
  ```
- **Получить подписку по ID**
  ```bash
  curl.exe -X GET http://localhost:9000/subscriptions/6  
  ```