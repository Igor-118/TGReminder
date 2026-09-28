# 0. Колонка `tg_id` и починка миграции

## Зачем

Пользователя нужно находить по его Telegram ID: он приходит в каждом апдейте (`update.Message.From.ID`). Сейчас такой колонки нет.

Заодно текущая миграция `migrations/20260925123024_init.sql` не накатывается из-за синтаксических ошибок:
- имена колонок в одинарных кавычках (`'name'`), так в Postgres пишутся строки, а не имена;
- запятая после последней колонки;
- нет `;` после `CREATE TABLE`;
- нет запятой после `dateid integer primary key`;
- внешний ключ записан как `userid integer foreign key`, а нужно `user_id bigint references users(id)`;
- в `Down` заглушка вместо удаления таблиц.

## Что сделать

Миграцию ещё никто не накатывал, поэтому правим прямо `20260925123024_init.sql`, новую не создаём.

- `users`:
  - `id` → `bigserial primary key`;
  - добавить `tg_id bigint not null unique`;
- `dates`: поправить синтаксис, `user_id bigint not null references users(id)`;
- `Down`: `drop table` обеих таблиц (сначала `dates`, потом `users`).

Telegram ID не влезает в `integer` (int32), поэтому `bigint`.

## Готово, когда

- [ ] `goose up` проходит без ошибок
- [ ] `goose down` откатывает, повторный `goose up` снова проходит
- [ ] вставка двух пользователей с одинаковым `tg_id` падает с ошибкой уникальности
