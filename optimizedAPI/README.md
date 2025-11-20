# optimizedAPI - Number Sorting API Optimization Project

## Описание проекта

REST API сервис для сортировки массива чисел.

**Основной эндпоинт:**
- `POST /sort` - принимает JSON с массивом чисел и возвращает отсортированный массив

**Пример запроса:**
```json
{
  "numbers": [5, 2, 8, 1, 9]
}
```

---

## Анализ производительности (версия 1.0 - неоптимизированная)

### Результаты нагрузочного тестирования (Vegeta)

**Параметры теста:**
- Rate: 50 req/s
- Duration: 30s
- Total requests: 1500

**Метрики:**
| Метрика | Значение |
|---------|----------|
| Success Rate | 100.00% |
| Mean Latency | 1.306 ms |
| 50th Percentile | 836.254 µs |
| 90th Percentile | 2.522 ms |
| 95th Percentile | 2.718 ms |
| 99th Percentile | 3.453 ms |
| Max Latency | 18.932 ms |
| Throughput | 50.03 req/s |

### CPU Профилирование (pprof - 30s)

**Top горячие точки:**

1. **`handler.bubbleSort`** - **45.83% CPU** (0.66s из 1.44s)
    - Основная вычислительная нагрузка
    - Неэффективный алгоритм O(n²)
    - Возможность оптимизации: замена алгоритма

2. **`encoding/json.(*decodeState).value`** - **21.53% CPU** (0.31s)
    - Десериализация JSON запроса
    - Рекурсивный парсинг структур

3. **`encoding/json.(*decodeState).object`** - **21.53% CPU** (0.31s)
    - Парсинг JSON объектов

4. **`encoding/json.(*Decoder).Decode`** - **25.00% CPU** (0.36s)
    - Общие затраты на JSON декодирование

**Ключевые наблюдения:**
- Bubble Sort - явное узкое место (46% CPU)
- JSON десериализация занимает ~25% CPU

### Benchmark результаты

```
BenchmarkSort-8   	  358591	      4712 ns/op
```

### Оптимизация

**Опираясь на результаты анализа, ключевой проблемой является неэффективный алгоритм сортировки. Для начала заменим его.**

---
## Версия 1.1 - замена алгоритма сортировки

**Изменения:**

- Самописный алгоритм пузырьковой сортировки был заменен на sort.Ints (Pattern-Defeating Quick Sort) 

## Анализ производительности

### Результаты нагрузочного тестирования (Vegeta)

**Параметры теста:**
- Rate: 50 req/s
- Duration: 30s
- Total requests: 1500

**Метрики:**
| Метрика | Значение |
|---------|----------|
| Success Rate | 100.00% |
| Mean Latency | 1.522 ms |
| 50th Percentile | 1.511 ms |
| 90th Percentile | 1.666 ms |
| 95th Percentile | 1.742 ms |
| 99th Percentile | 2.347 ms |
| Max Latency | 4.242 ms |
| Throughput | 50.03 req/s |

### CPU Профилирование (pprof - 30s)

**Top горячие точки:**

1. **`encoding/json.(*Decoder).Decode`** - **37.20% CPU** (0.61s)
    - Общие затраты на JSON декодирование

2. **`encoding/json.(*Decoder).Encode`** - **16.46% CPU** (0.27s)
    - Общие затраты на JSON кодирование

3. **`http(*response).finishRequest`** - **15.85% CPU** (0.26s)
    - Затраты на http ответ

4. **`slices.pdqsortOrdered[go shape it]`** - **11.59% CPU** (0.19s)
    - Значительно улучшена скорость сортировки

### Benchmark результаты

```
BenchmarkSort-8   	  290648	      3514 ns/op
```

### Оптимизация

**Замена алгоритма сортировки значительно улучшила производительность.**
**Теперь главную нагрузку создает обработка JSON. Следующим шагом заменю стандартный json парсер**

---

## Версия 1.2 - замена JSON парсера

**Изменения:**

- Самописный алгоритм пузырьковой сортировки был заменен на sort.Ints (Pattern-Defeating Quick Sort)

## Анализ производительности

### Результаты нагрузочного тестирования (Vegeta)

**Параметры теста:**
- Rate: 50 req/s
- Duration: 30s
- Total requests: 1500

**Метрики:**
| Метрика | Значение |
|---------|----------|
| Success Rate | 100.00% |
| Mean Latency | 1.062 ms |
| 50th Percentile | 1.169 ms |
| 90th Percentile | 1.401 ms |
| 95th Percentile | 1.55 ms |
| 99th Percentile | 2ю261 ms |
| Max Latency | 5.481 ms |
| Throughput | 50.03 req/s |

### CPU Профилирование (pprof - 30s)

**Top горячие точки:**

1. **`handler.easyjson888c126aDecodeGithubComKost0L4InternalHandler`** - **20.59% CPU** (0.21s)
    - Общие затраты на JSON декодирование (заметно лучше)

2. **`syscall.Syscall6`** - **20.59% CPU** (0.21s)
    - Системные вызовы

3. **`http(*response).finishRequest`** - **19.61% CPU** (0.20s)
    - Затраты на http ответ

4. **`handler.RequestPayload.MarshalJSON`** - **13.73% CPU** (0.14s)
    - Общие затраты на JSON кодирование (заметно лучше)

### Benchmark результаты

```
BenchmarkSort-8   	  1587265	      766 ns/op
```

### Оптимизация

**Замена стандартного парсера на easyjson заметно улучшила процессы кодирования и декодирования JSON-ов**
**Дальше попробую отказаться от стандартного http пакета в пользу более быстрого**

---