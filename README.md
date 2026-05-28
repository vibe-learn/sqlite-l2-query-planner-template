        # sqlite — Планировщик запросов

        Homework-шаблон для урока **l2_query_planner** (Планировщик запросов) на платформе Vibe Learn.

        ## Что делать

        На Go с modernc.org/sqlite (без CGO) реализуй обёртку-аудитор планов. Функции:
Plan(db, query, args...) → []string — строки EXPLAIN QUERY PLAN (поле detail); UsesIndex(plan)
→ bool — есть ли SEARCH ... USING INDEX (а не только SCAN); IsCovering(plan) → bool —
встречается ли COVERING INDEX; SuggestFix(db, query) → string — эвристика: если SCAN на
таблице с >N строк, подсказать «нет подходящего индекса / возможно функция на колонке».
Тесты на :memory:-базе: наполняют таблицу, проверяют, что запрос по неиндексированной колонке
даёт SCAN (UsesIndex=false), после CREATE INDEX — SEARCH (true); что запрос только по
индексным колонкам даёт COVERING; что WHERE lower(col)=? остаётся SCAN, пока нет
expression-индекса. Где уместно — ANALYZE перед проверкой.

## Контекст (из transfer-задачи урока)

В проде таблица events на 20 млн строк: events(id INTEGER PRIMARY KEY, user_id INTEGER,
type TEXT, created_at INTEGER, payload TEXT). Есть индекс idx_created ON events(created_at).
Жалоба: дашборд «события конкретного юзера за последние 7 дней» отдаётся 9 секунд:

```sql
SELECT type, count(*) FROM events
WHERE user_id = 12345
  AND created_at >= strftime('%s','now','-7 days')
GROUP BY type;
```

## Recap из урока

- **`EXPLAIN QUERY PLAN`** (EQP) — главный инструмент в SQLite. Не выполняет запрос. `EXPLAIN ANALYZE` — это PostgreSQL, в SQLite его нет; время мерь отдельно (`.timer on`).
- Читай план: `SEARCH ... USING INDEX` — хорошо (точечно), `SCAN table` — полный проход (плохо на больших, нормально на маленьких), `COVERING INDEX` — ответ из индекса без таблицы.
- **`ANALYZE`** собирает статистику в `sqlite_stat1`; без неё планировщик гадает. `PRAGMA optimize` запускает ANALYZE по необходимости — зови регулярно.
- Индекс игнорируется при: **функции/арифметике на колонке** (нужен expression-индекс), отсутствии ведущей колонки (leftmost prefix), низкой селективности, отсутствии статистики.
- NGQP (с 3.8.0) выбирает планы по стоимости и порядок соединений; покрывающий индекс — часто самый дешёвый план.

        ## Как работать

        1. Платформа Vibe Learn создаёт копию этого репо в твоём GitHub-аккаунте по клику «Начать домашку» на странице урока (через GitHub `/generate`, codecrafters-pattern).
        2. Склонируй копию локально, реализуй TODO в `main.go`, прогони тесты, запушь.
        3. CI (`.github/workflows/ci.yml`) запускает `go vet` + `go test ./...` на каждый push. Платформа слушает результат через webhook от GitHub Actions и обновляет статус домашки на странице урока.

        ## Локальное окружение

        - Go 1.22+
        - SQLite встроена — **никакого сервера и docker-compose**. БД это один файл (`DATABASE_PATH`) или `:memory:` в тестах. Драйвер `modernc.org/sqlite` — чистый Go, без CGO, так что CI собирается без компилятора C.

        ## Запуск

        ```bash
        # Прогнать тесты (бегут на :memory:-базе, ничего поднимать не нужно)
        go test ./...

        # Запустить main (создаёт схему, печатает marker; замени stub на реализацию)
        go run .

        # На файловой БД (нужно для WAL/concurrency-уроков):
        DATABASE_PATH=./app.db go run .
        ```

        ## Заметка автора

        Это baseline-шаблон, сгенерированный платформой. Бизнес-сущность задачи (что конкретно реализовать в `main.go`, какие тесты сделать строгими) расширяется по ходу итераций — параллельно с углублением теории урока.
