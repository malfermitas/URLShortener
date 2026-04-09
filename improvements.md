## 🔬 Профилирование и оптимизация

Проект включает инструменты для профилирования CPU и памяти с использованием `net/http/pprof`.

### Профилирование CPU

```bash
# Запустить профиль на 30 секунд
curl http://localhost:8080/debug/pprof/profile?seconds=30 > cpu.pprof

# Анализ в терминале
go tool pprof cpu.pprof

# Или через веб-интерфейс
go tool pprof -http=:8081 cpu.pprof
```

### Профилирование памяти

```bash
curl http://localhost:8080/debug/pprof/heap > heap.pprof
go tool pprof -http=:8081 heap.pprof
```

### Доступные эндпоинты pprof

| Endpoint | Описание |
|----------|----------|
| `/debug/pprof/` | Индекс всех профилей |
| `/debug/pprof/heap` | Memory (Allocations) |
| `/debug/pprof/profile?seconds=30` | CPU (30 сек) |
| `/debug/pprof/goroutine` | Текущие горутины |
| `/debug/pprof/block` | Блокировки |

## 📈 История оптимизаций

### v1.0 - Исходная версия

**Проблемы:**
- Генератор ключей использовал `math/rand` с mutex - узкое место при параллельных запросах
- При создании URL без кастомного кода выполнялось до 5 SQL-запросов для проверки уникальности
- Validator создавался на каждый запрос (ненужные аллокации)
- Множественные вызовы `GetTraceID()` в каждом логе

**Результат профилирования:**
- `urlService.Create` занимал ~45% CPU времени
- До 5 SQL запросов на один Create-запрос

### v1.1 - Оптимизация генератора ключей

**Изменения:**
- Заменён `math/rand` на Snowflake-подобный алгоритм
- Гарантированная уникальность без проверки в БД
- Убран mutex в пользу атомарных операций

**Результат:**
- Генерация ключей: 119 ns/op → **31 ns/op** (3.8x быстрее)
- SQL запросов на Create: до 5 → **1**
- Время Create в профиле: 45% → ~29%

### v1.2 - Оптимизация middleware и handlers

**Изменения:**
- Переход с `gin.Default()` на `gin.New()` + ручное добавление middleware
- Validator перенесён в конструктор handler (один раз при старте)
- Убраны дублирующиеся вызовы trace_id в логах

**Результат:**
- Уменьшен overhead от framework
- Меньше аллокаций на каждый запрос

### v1.3 - Оптимизация памяти

**Изменения:**
- OpenTelemetry sampling изменён на sdktrace.TraceIDRatioBased(0.2)
- Теперь spans создаются в 20% случаев

**Результат:**
- Общая память: 54.94 MB → **28.77 MB** (47% меньше)
- OpenTelemetry: ~4 MB → практически 0
- Основное потребление: bufio буферы (70% - норма для HTTP сервера)

### Бенчмарки

```bash
go test -bench=. -benchmem ./internal/adapter/out/generator/
```

Результаты:
```
BenchmarkURLGenerator_Generate-12           37469672    31.92 ns/op    8 B/op    1 allocs/op
BenchmarkURLGenerator_Generate_Parallel-12   16590602    69.56 ns/op    8 B/op    1 allocs/op
```

## 🧪 Тестирование

### Unit-тесты

```bash
go test ./...
```

### Benchmark-тесты

```bash
go test -bench ./internal/adapter/out/generator/
```

### Тест уникальности генератора

```bash
go test -v -run TestSnowflakeGenerator_Uniqueness ./internal/adapter/out/generator/
```
