# messenger-server

Сервер для мобильного мессенджера [messenger-android](https://github.com/vettich/messenger-android).

Работает с СУБД RethinkDB, предоставляет API в формате GraphQL.

## Запуск

```sh
$ make
```

## Роутинги

- `/graphql` - endpoint для запросов
- `/` - GraphiQL - простой playground
- `/play` - Playground - более симпатичный playground (но не работают подписки)
