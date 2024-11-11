use actix_cors::Cors;
use actix_session::config::PersistentSession;
use actix_session::storage::RedisSessionStore;
use actix_session::Session;
use actix_session::{storage::CookieSessionStore, SessionMiddleware};
use actix_web::cookie::{time, Key};
use actix_web::{get, web, App, HttpResponse, HttpServer, Responder};
use config::actix::AppState;
use config::local::oidc_config::get_oidc_config;
use log::debug;
use log::error;
use log::info;
use openid::error::Error;
use openid::Bearer;
use openid::DiscoveredClient;
use openid::IdToken;
use openid::TokenIntrospection;
use openid::Userinfo;
use openid::{Client, Options, StandardClaims, Token};
use reqwest::Client as HttpClient;
use sea_orm::DatabaseConnection;
use security::oidc::get_client;
use serde::{Deserialize, Serialize};
use serde_json::Value;
use std::collections::HashMap;
use std::fmt::Display;
use std::sync::Arc;
use std::sync::Mutex;
use std::time::Duration;
use url::Url;

mod config;
mod domain;
mod dto;
mod http;
mod repositories;
mod security;

#[derive(Debug, Deserialize)]
struct AuthRequest {
    code: Option<String>,
    session_state: Option<String>,
    iss: Option<String>,
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    config::logger::config();

    log::info!("Starting the application");
    config::actix::config().await

    // TODO faudra trouver un moyen de close la connexion. Mais là on peut pas move la static_db
}

///////// SHARED MUTABLE STATE //////////
#[derive(Serialize)]
struct GetResponse {
    key: String,
    value: Option<Bearer>,
}

// Handler to set a key-value pair
async fn set_key(data: web::Data<AppState>) -> impl Responder {
    let mut store = data.store.lock().unwrap();
    HttpResponse::Ok().json("Key set successfully")
}

// Handler to get a value by key
async fn get_key(data: web::Data<AppState>) -> impl Responder {
    let store = data.store.lock().unwrap();
    let value = store.get(&"hello".to_string()).cloned();
    HttpResponse::Ok().json(GetResponse {
        key: "hello".to_string(),
        value: value,
    })
}

////////// SESSION //////////

// Route to retrieve data from the session
async fn get_session(
    session: Session,
    state: web::Data<AppState>,
    _: web::Query<AuthRequest>,
) -> impl Responder {
    let user_id = session.get::<Bearer>("user_id");

    match user_id {
        Ok(user_id) => match user_id {
            Some(user_id) => {
                info!("User ID found in session: {}", user_id.access_token.clone());
                HttpResponse::Ok().body(format!("Welcome back, user {}!", user_id.access_token.clone()))
            }
            None => {
                error!("No user ID found in session");
                HttpResponse::Ok().body("No user ID found in session")
            }
        },
        Err(e) => {
            error!("No session data found: {}", e);
            HttpResponse::Ok().body(format!("No session data found: {}", e))
        }
    }
}


////////// OIDC //////////
#[get("/login")]
async fn login(
    session: Session,
    // client: web::Data<Arc<Mutex<Client<openid::Discovered, StandardClaims>>>>,
    state: web::Data<AppState>,
    query: web::Query<AuthRequest>,
) -> impl Responder {
    let client = state.client.lock().unwrap();

    info!("Login with query: {:?}", query);

    match &query.code {
        Some(code) => {
            let authorization_code: &&String = &code;

            info!(
                "Requesting token with received authorization code: {}",
                authorization_code
            );

            match client.authenticate(authorization_code, None, None).await {
                Ok(token) => {
                    let access_token = token.bearer.access_token.clone();
                    let id_token = token.bearer.id_token.clone();

                    // Save in session
                    let saved = session.insert(
                        "user_id",
                        token.bearer.clone(),
                    );
                    match saved {
                        Ok(_) => info!("Token saved in session"),
                        Err(e) => error!("Error saving token in session: {}", e),
                    }

                    // Save in shared state
                    let mut store = state.store.lock().unwrap();
                    store.insert("hello".to_string(), token.bearer.clone());

                    HttpResponse::Ok().json(HashMap::from([
                        ("access_token", access_token),
                        ("id_token", id_token.unwrap_or_default()),
                    ]))
                }
                Err(err) => {
                    error!("Error exchanging code for token: {:?}", err);
                    HttpResponse::InternalServerError().body("Failed to authenticate")
                }
            }
        }
        None => {
            info!("No code provided. Starting authentication.");

            // Il faut définir un nonce et max age ici pour réutiliser à priori
            let auth_url = client.auth_url(&Options {
                scope: Some("openid email profile".into()),
                ..Default::default()
            });

            HttpResponse::Found()
                .append_header(("Location", auth_url.to_string()))
                .finish()
        }
    }
}

// GETTING INFO FROM ID TOKEN
/*
    How to use:
*/
fn get_info_from_id_token(id_token: &IdToken<StandardClaims>) -> Userinfo {
    id_token.payload().cloned().unwrap().userinfo
}

// VALIDATING/INTROSPECTING TOKEN
/*
    How to use:
    // Validating access token
    let token_from_bearer = introspect_token_from_bearer(&client, token.bearer.clone()).await;
    // Validating from token
    let introspection = introspect_token(&client, &token).await;

    match introspection {
        Ok(intro) => info!("Token introspection successful: {:?}", intro),
        Err(e) => error!("Token introspection failed: {:?}", e)
    }
    // End validating access token
*/
async fn introspect_token(
    client: &DiscoveredClient,
    token: &Token,
) -> Result<TokenIntrospection<StandardClaims>, Error> {
    client.request_token_introspection(&token).await
}

async fn introspect_token_from_bearer(
    client: &DiscoveredClient,
    bearer: Bearer,
) -> Result<TokenIntrospection<StandardClaims>, Error> {
    let token = Token::from(bearer);
    client.request_token_introspection(&token).await
}
