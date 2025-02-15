# Запуск проекта

```shell
make run
```

# Запуск тестов
1. 
    ```shell
    make e2e-stand
    ```
2. 
    ```shell
    make test-coverage
    ```
# Результаты нагрузочного тетсирования
Результаты напрямую зависят от количества вызовов auth ручки сервиса, поскольку хеширование и сравнение пароля занимает +-200ms.
## При убранной задержке от хеширования
```
http_req_blocked...............: avg=157.73µs min=0s     med=0s      max=1.39ms   p(90)=557.72µs p(95)=729.64µs
http_req_connecting............: avg=120.21µs min=0s     med=0s      max=998.9µs  p(90)=526.98µs p(95)=561.2µs
http_req_duration..............: avg=43.2ms   min=10.9ms med=20.42ms max=572.06ms p(90)=83.68ms  p(95)=108.63ms
  { expected_response:true }...: avg=43.2ms   min=10.9ms med=20.42ms max=572.06ms p(90)=83.68ms  p(95)=108.63ms
http_req_failed................: 0.00%   0 out of 200
http_req_receiving.............: avg=10.51ms  min=0s     med=0s      max=482.73ms p(90)=999.2µs  p(95)=1.3ms
http_req_sending...............: avg=83.37µs  min=0s     med=0s      max=5.45ms   p(90)=32.19µs  p(95)=538.62µs
http_req_tls_handshaking.......: avg=0s       min=0s     med=0s      max=0s       p(90)=0s       p(95)=0s
http_req_waiting...............: avg=32.6ms   min=10.9ms med=20.42ms max=113.18ms p(90)=79.47ms  p(95)=84.82ms
```
## Без авторизации
```
http_req_blocked...............: avg=18.62µs min=0s      med=0s     max=1.03ms   p(90)=0s      p(95)=0s
http_req_connecting............: avg=15.86µs min=0s      med=0s     max=1ms      p(90)=0s      p(95)=0s
http_req_duration..............: avg=13.97ms min=1.85ms  med=7.49ms max=402.58ms p(90)=18.52ms p(95)=62.32ms
  { expected_response:true }...: avg=13.97ms min=1.85ms  med=7.49ms max=402.58ms p(90)=18.52ms p(95)=62.32ms
http_req_failed................: 0.00%   0 out of 1371
http_req_receiving.............: avg=30.63µs min=0s      med=0s     max=1.26ms   p(90)=0s      p(95)=52.1µs
http_req_sending...............: avg=9.25µs  min=0s      med=0s     max=5ms      p(90)=0s      p(95)=0s
http_req_tls_handshaking.......: avg=0s      min=0s      med=0s     max=0s       p(90)=0s      p(95)=0s
http_req_waiting...............: avg=13.93ms min=1.85ms  med=7.45ms max=402.05ms p(90)=18.42ms p(95)=62.32ms
```

## Как есть
```
http_req_blocked...............: avg=1.25ms  min=0s   med=0s    max=28.05ms  p(90)=1.9ms    p(95)=6.66ms
http_req_connecting............: avg=1.2ms   min=0s   med=0s    max=28.05ms  p(90)=1.55ms   p(95)=6.66ms
http_req_duration..............: avg=1.31s   min=13ms med=1.22s max=3.64s    p(90)=3.33s    p(95)=3.44s
  { expected_response:true }...: avg=1.31s   min=13ms med=1.22s max=3.64s    p(90)=3.33s    p(95)=3.44s
http_req_failed................: 0.00%   0 out of 200
http_req_receiving.............: avg=5.26ms  min=0s   med=0s    max=246.47ms p(90)=999.53µs p(95)=6.59ms
http_req_sending...............: avg=57.96µs min=0s   med=0s    max=1ms      p(90)=0s       p(95)=515.75µs
http_req_tls_handshaking.......: avg=0s      min=0s   med=0s    max=0s       p(90)=0s       p(95)=0s
http_req_waiting...............: avg=1.31s   min=13ms med=1.22s max=3.64s    p(90)=3.26s    p(95)=3.44s
```

# Возникшие вопросы
1. Почему ручка /api/buy/{item} имеет запрос типа GET - решено оставить GET, для совместимости с имеющимся swagger api;
2. Может ли 1 монетка подразделяться на более мелкие "копейки" - судя по примерам ценам, было решено, что нет, монетки могут быть только целыми.