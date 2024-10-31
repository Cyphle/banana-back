use std::{ops::{Deref, DerefMut}, sync::Arc};
use openid::{Client, DiscoveredClient, Options};
use reqwest::Url;
use sea_orm::DatabaseConnection;
use serde::Deserialize;
use serde_json::Value;
use tokio::sync::Mutex;
use crate::domain::profile::CreateProfileCommand;

mod config;
mod repositories;
mod domain;
mod http;
mod dto;

#[derive(Debug, Deserialize)]
struct AuthRequest {
    code: String,
}

#[actix_web::main]
async fn main() {
    config::logger::config();

    log::info!("Starting the application");

    let db = config::database::connect().await.unwrap();
    let static_db = Box::leak(Box::new(db));


    /* OPENID CONNECT */
    let issuer_url = reqwest::Url::parse("http://localhost:8181/realms/Banana").unwrap();
    let client_id = "banana";
    let client_secret = "banana-secret";
    let redirect_uri = Some("http://localhost:9000/callback");
    
    let client: Client<openid::Discovered, Value> = Client::discover(
        client_id.to_string(),
        client_secret.to_string(),
        Some(redirect_uri.to_string()),
        issuer_url,
    )
    .await
    .expect("Failed to discover OpenID configuration");

    // Wrap client in Arc and Mutex for sharing across Actix handlers
    let client = Arc::new(Mutex::new(client));

    /* END OPENID CONNECT */

    // let _ = config::actix::config(static_db).await;

    // repositories::profiles::create(static_db, &CreateProfileCommand {
    //     username: "johndoe".to_string(),
    //     email: "johndoe".to_string(),
    //     first_name: "John".to_string(),
    //     last_name: "Doe".to_string(),
    // }).await.unwrap();

    log::info!("Application is now closed");

    // TODO faudra trouver un moyen de close la connexion. Mais là on peut pas move la static_db
}

/* 
WORKING OIDC EXAMPLE

use actix_cors::Cors;
use actix_web::{web, App, HttpResponse, HttpServer, Responder};
use openid::{Client, Options, StandardClaims, Token};
use serde::Deserialize;
use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::Mutex;
use url::Url;

#[derive(Debug, Deserialize)]
struct AuthRequest {
    code: String,
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    // Set up Keycloak OIDC parameters
    let issuer_url = reqwest::Url::parse("http://localhost:8181/realms/Banana");
    let client_id = "banana";
    let client_secret = "banana-secret";
    let redirect_uri = "http://localhost:8080/auth/callback";

    // Initialize OpenID Client with Keycloak discovery
    let client: Client<openid::Discovered, StandardClaims> = Client::discover(
        client_id.to_string(),
        client_secret.to_string(),
        Some(redirect_uri.to_string()),
        issuer_url.unwrap(),
    )
    .await
    .expect("Failed to discover OpenID configuration");

    // Wrap client in Arc and Mutex for sharing across Actix handlers
    let client = Arc::new(Mutex::new(client));

    // Start Actix server
    HttpServer::new(move || {
        App::new()
            .wrap(
                Cors::default()
                    .allowed_origin("http://localhost:9000") // Change to your frontend URL
                    .allowed_methods(vec!["GET", "POST"])
                    .allowed_headers(vec![actix_web::http::header::CONTENT_TYPE])
                    .supports_credentials() // Optional, if credentials are used
                    .max_age(3600),
            )
            .app_data(web::Data::new(client.clone()))
            .route("/", web::get().to(authenticate))
            .route("/auth/callback", web::get().to(auth_callback))
    })
    .bind("127.0.0.1:8080")?
    .run()
    .await
}

// Redirect user to the authorization URL
async fn authenticate(
    client: web::Data<Arc<Mutex<Client<openid::Discovered, StandardClaims>>>>,
) -> impl Responder {
    let client = client.lock().await;
    let auth_url = client.auth_url(&Options {
        scope: Some("openid email profile".into()),
        ..Default::default()
    });
    HttpResponse::Found()
        .append_header(("Location", auth_url.to_string()))
        .finish()
}

// Handle the authorization callback
async fn auth_callback(
    client: web::Data<Arc<Mutex<Client<openid::Discovered, StandardClaims>>>>,
    query: web::Query<AuthRequest>,
) -> impl Responder {
    let client = client.lock().await;
    let code = &query.code;

    println!("Received authorization code: {}", code);

    // Exchange the authorization code for tokens
    match client.request_token(code).await {
        Ok(token) => {
            println!("Token: {:?}", token);
            let access_token = token.access_token.clone();
            let id_token = token.id_token.clone();

            let token_wrapper = Token::from(token);

            // Optional: Fetch user info
            if let Ok(userinfo) = client.request_userinfo(&token_wrapper).await {
                return HttpResponse::Ok().json(userinfo);
            }

            // If user info is not fetched, return tokens only
            HttpResponse::Ok().json(HashMap::from([
                ("access_token", access_token),
                ("id_token", id_token.unwrap_or_default()),
            ]))
        }
        Err(err) => {
            eprintln!("Error exchanging code for token: {:?}", err);
            HttpResponse::InternalServerError().body("Failed to authenticate")
        }
    }
}

[dependencies]
actix-web = "4"
openid = "0.15"
reqwest = { version = "0.11", features = ["json"] }
tokio = { version = "1", features = ["full"] }
serde = { version = "1.0", features = ["derive"] }
serde_json = "1.0"
url = "2.2"
actix-cors = "0.6"
*/