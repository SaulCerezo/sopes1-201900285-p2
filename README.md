# sopes1-201900285-p2

## Rust-api creando la imagen

Hicimos un endpoint POST /input que:

* Recibe un array JSON con tweets del clima.

* Guarda los datos en memoria con Mutex Vec WeatherTweet.

* Responde con el número de tweets recibidos y el total acumulado.

Se edito en main.rs para que escuche desde 0.0.0.0
```
    .bind(("0.0.0.0", 8000))?

```
Esto permite que el contenedor acepte conexiones desde fuera, como desde Postman o curl.

### Comandos que usamos para rust-api
```
    docker build -t rust-api:local .
    docker run --rm -p 8000:8000 rust-api:local
```
## Servidro de Go
creamos un sevidor y un cliente de forma local para ver si los tweets estan llegando 

```
    go run main.go
```

## Locust
Herramienta para pruebas de carga que te permite simular muchos usuarios haciendo peticiones a tu API
```
    locust --host http://localhost:8000 
```

Docker compose
Creamos un docker compose para ejecutar los servicos y exponerlos a la misma red

```
    docker compose up
    docker compose down
```
```
    docker compose restart analyzer
```
## Subiendo las imagenes a dockerhub
Etiqueta correctamente tus imágenes con el nombre correcto:

```
    docker tag tweetsclima-rust-api saulcerezo/sopes1p2:rust-api-v1
    docker tag tweetsclima-go-entry saulcerezo/sopes1p2:go-entry-v1
    docker tag tweetsclima-analyzer saulcerezo/sopes1p2:analyzer-v1

```
hacemos push de cada una

```
    docker push saulcerezo/sopes1p2:rust-api-v1
    docker push saulcerezo/sopes1p2:go-entry-v1
    docker push saulcerezo/sopes1p2:analyzer-v1

```
-----------------------------------------

# DOCUMENTACION

### ¿Cómo funciona Kafka?
Apache Kafka es una plataforma de mensajería distribuida que permite transmitir, almacenar y procesar flujos de datos en tiempo real. Se basa en un modelo publish-subscribe, donde:

* Productores publican mensajes en un tópico.
* Consumidores se suscriben a ese tópico y reciben mensajes.
* Kafka almacena mensajes de forma persistente y permite a múltiples consumidores leer a su ritmo.
* Está diseñado para ser escalable, tolerante a fallos y de alto rendimiento.

```
    ⚠️ Para este proyecto no utilice kafka pero si 
    hice uso de RabbitMQ que cumplen funciones similares (mensajeria)
```

### ¿Cómo difiere Valkey de Redis?
Valkey es un fork (derivado) de Redis, creado después de que Redis Labs cambió su licencia a una menos permisiva.


| Característica | Redis | Valkey |
| ------------- | ------------- | ------------- |
| Licencia | Redis Source Available MIT (libre y abierta) (RSAL) | nuevo |
|Comunidad|Mantenido por Redis Inc.| Mantenido por la comunidad open source|
|Filosofía|Más control comercial|Enfocado en la comunidad y libertad
|

```
    ⚠️ Para este proyecto no se utlizo ni valkey ni Redis
```
### ¿Es mejor gRPC que HTTP?
Depende del caso de uso.

gRPC es un protocolo de comunicación basado en HTTP/2 que utiliza Protocol Buffers. Es ideal para:

* Servicios internos entre microservicios.
* Comunicación rápida y eficiente en entornos de alto rendimiento.
* Necesidad de tipado fuerte y streaming bidireccional.

Ventajas de gRPC sobre HTTP REST:
* Mejor rendimiento (más ligero que JSON).
* Soporte de streaming.
* Contratos bien definidos mediante .proto.

```
En este proyecto use gRPC entre microservicios 
(Rust → Go), lo cual es adecuado y más eficiente 
que usar REST en ese tipo de comunicación.
```

### ¿Hubo una mejora al utilizar dos réplicas en los deployments de API REST y gRPC? Justifique su respuesta.

Sí, hay mejora, sobre todo en disponibilidad y balance de carga.

* Más réplicas = mayor capacidad de atender múltiples solicitudes simultáneas.
* Tolerancia a fallos: si una instancia falla, otra sigue atendiendo.
* Mejor rendimiento bajo carga alta (como la generada con Locust).

```
usar 2 réplicas en un clúster de GCP o Kubernetes 
mejora escalabilidad y disponibilidad del sistema.
```

### Para los consumidores, ¿Qué utilizó y por qué?
Se usó RabbitMQ como sistema de mensajería y un servicio Go llamado analyzer como consumidor:

* analyzer se suscribe a la cola de RabbitMQ y procesa los mensajes recibidos.
* Se eligió Go por su eficiencia en concurrencia y su buen soporte para integrarse con RabbitMQ.

```
Se eligio RabbitMQ por su simplicidad, soporte en 
Docker, y facilidad de uso en entornos de 
desarrollo.
```
