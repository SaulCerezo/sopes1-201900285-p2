use actix_web::{post, web, App, HttpResponse, HttpServer, Responder, middleware::Logger};
use serde::{Deserialize, Serialize};
use std::sync::Mutex;
use std::env;           // ✅ Importar para leer variables de entorno
use dotenv::dotenv;     // ✅ Importar para cargar .env
use env_logger;

pub mod weather {
    tonic::include_proto!("weather");
}
use weather::weather_service_client::WeatherServiceClient;
use weather::{Tweet, TweetBatch};
use tonic::Request;

#[derive(Debug, Deserialize, Serialize)]
struct WeatherTweet {
    description: String,
    country: String,
    weather: String,
}

struct AppState {
    tweets: Mutex<Vec<WeatherTweet>>,
}

#[post("/input")]
async fn input(
    tweets: web::Json<Vec<WeatherTweet>>,
    data: web::Data<AppState>
) -> impl Responder {
    let mut stored_tweets = data.tweets.lock().unwrap();
    let count = tweets.len();
    
    println!("Recibidos {} nuevos tweets (total: {})", count, stored_tweets.len() + count);
    
    stored_tweets.extend(tweets.into_inner());

    let grpc_tweets: Vec<Tweet> = stored_tweets
        .iter()
        .map(|t| Tweet {
            description: t.description.clone(),
            country: t.country.clone(),
            weather: t.weather.clone(),
        })
        .collect();

    // ✅ Leer la dirección del servidor gRPC desde la variable de entorno
    let grpc_addr = env::var("GRPC_SERVER").unwrap_or_else(|_| {
        println!("⚠️  GRPC_SERVER no definida, usando http://localhost:50051 por defecto");
        "http://localhost:50051".to_string()
    });

    let mut grpc_client = match WeatherServiceClient::connect(grpc_addr).await {
        Ok(client) => client,
        Err(err) => {
            println!("❌ Error al conectar con el servidor Go: {}", err);
            return HttpResponse::InternalServerError().json(serde_json::json!({
                "status": "error",
                "message": "No se pudo conectar con Go"
            }));
        }
    };

    let request = Request::new(TweetBatch {
        tweets: grpc_tweets,
    });

    match grpc_client.send_tweets(request).await {
        Ok(res) => {
            let ack = res.into_inner();
            println!("✅ Respuesta del servidor Go: {:?}", ack);
        }
        Err(err) => {
            println!("❌ Error al enviar tweets a Go: {}", err);
        }
    }

    HttpResponse::Ok().json(serde_json::json!({
        "status": "received",
        "count": count,
        "total": stored_tweets.len()
    }))
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    dotenv().ok(); // ✅ Cargar variables de entorno desde .env
    env_logger::init_from_env(env_logger::Env::new().default_filter_or("info"));

    let app_state = web::Data::new(AppState {
        tweets: Mutex::new(Vec::new()),
    });

    println!("API Rust arrancando en http://0.0.0.0:8000");

    HttpServer::new(move || {
        App::new()
            .app_data(app_state.clone())
            .wrap(Logger::default())
            .service(input)
    })
    .bind(("0.0.0.0", 8000))?
    .run()
    .await
}
