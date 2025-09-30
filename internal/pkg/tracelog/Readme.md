# tracelog - Обертка над трейсингом opentelemetry и логгером zerolog

Примеры использования приведены в `examples/`

Все примеры написаны с условием, что в качестве экспортера трейсов будет выступать Jaeger.
Если требуется какой-либо другой экспортер, то некоторые из них есть [здесь](https://github.com/open-telemetry/opentelemetry-go/tree/main/exporters)

Запустить Jaeger локально в докере:

    docker run -d --name jaeger -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
    -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778 -p 16686:16686 \
    -p 14268:14268 -p 9411:9411 --restart unless-stopped jaegertracing/all-in-one:1.6

- Web GUI: http://localhost:16686
- Слать запросы на:
  - localhost:6831 (UDP)
  - http://localhost:14268/api/traces (TCP)