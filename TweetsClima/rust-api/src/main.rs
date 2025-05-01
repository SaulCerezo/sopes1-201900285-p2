use actix_web::{post, web, App, HttpResponse, HttpServer, Responder, middleware::Logger};
use serde::{Deserialize, Serialize};
use std::sync::Mutex;
use env_logger;

#[derive(Debug, Deserialize, Serialize)]
struct WeatherTweet {
    description: String,  // Cambiado a minúscula para seguir convención Rust
    country: String,
    weather: String,
}

// Estado compartido para almacenar tweets
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
    
    // Almacenar los tweets
    stored_tweets.extend(tweets.into_inner());
    
    HttpResponse::Ok().json(serde_json::json!({
        "status": "received",
        "count": count,
        "total": stored_tweets.len()
    }))
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    // Inicializar logger
    env_logger::init_from_env(env_logger::Env::new().default_filter_or("info"));
    
    // Estado compartido
    let app_state = web::Data::new(AppState {
        tweets: Mutex::new(Vec::new()),
    });

    println!("API Rust arrancando en http://127..0.1:8000");
    
    HttpServer::new(move || {
        App::new()
            .app_data(app_state.clone())
            .wrap(Logger::default())  // Middleware para logging
            .service(input)
    })
    .bind(("0.0.0.0", 8000))?
    .run()
    .await
}