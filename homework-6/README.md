# Домашнее задание 6

Логирование с использованием Kafka. Логи пишутся в консоль. `Consumer` и `Producer` находятся в директории `internal/infrastructure/kafka`.
`LogSender` и `LogReceiver` находятся в директории `internal/app/logging`.

Запись логов происходит в топик `logs` 

Интеграционные тесты используют топик с названием logs_test для записи логов