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
