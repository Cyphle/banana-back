use actix_cors::Cors;
use actix_session::{SessionMiddleware, storage::CookieSessionStore};
use actix_session::Session;
use actix_web::cookie::Key;
use actix_web::{get, web, App, HttpResponse, HttpServer, Responder};
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
use serde::Deserialize;
use serde_json::Value;
use std::collections::HashMap;
use std::sync::Arc;
use std::time::Duration;
use tokio::sync::Mutex;
use url::Url;

mod config;
mod domain;
mod dto;
mod http;
mod repositories;

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

    // let db = config::database::connect().await.unwrap();
    // let static_db = Box::leak(Box::new(db));

    // ACTIX
    // let _ = config::actix::config(static_db).await;

    // repositories::profiles::create(static_db, &CreateProfileCommand {
    //     username: "johndoe".to_string(),
    //     email: "johndoe".to_string(),
    //     first_name: "John".to_string(),
    //     last_name: "Doe".to_string(),
    // }).await.unwrap();

    // log::info!("Application is now closed");

    /* OPENID CONNECT */
    // Set up Keycloak OIDC parameters
    let issuer_url = reqwest::Url::parse("http://localhost:8181/realms/Banana");
    let client_id = "banana";
    let client_secret = "banana-secret";
    let redirect_uri = "http://localhost:9000/login";

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
    // Generate a secure 32-byte key for cookie signing (use a random key in production)
    let secret_key = Key::generate();

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
            .wrap(SessionMiddleware::new(CookieSessionStore::default(), secret_key.clone()))
            .app_data(web::Data::new(client.clone()))
            .route("/set_session", web::get().to(set_session))
            .route("/get_session", web::get().to(get_session))
            .service(login)
    })
    .bind("127.0.0.1:8080")?
    .run()
    .await

    /* END OPENID CONNECT */

    // TODO faudra trouver un moyen de close la connexion. Mais là on peut pas move la static_db
}

////////// SESSION //////////
// Route to set data in the session
async fn set_session(session: Session) -> impl Responder {
    let res = session.insert("user_id", 42);
    match res {
        Ok(_) => HttpResponse::Ok().body("Session data set"),
        Err(e) => HttpResponse::Ok().body(format!("Error setting session data: {}", e)),
    }
}

// Route to retrieve data from the session
async fn get_session(session: Session) -> impl Responder {
    let user_id = session.get::<i32>("user_id");

    match user_id {
        Ok(user_id) => {
            match user_id {
                Some(user_id) => HttpResponse::Ok().body(format!("Welcome back, user {}!", user_id)),
                None => HttpResponse::Ok().body("No user ID found in session"),
            }
        },
        Err(e) => HttpResponse::Ok().body(format!("No session data found: {}", e)),
    }
}

////////// OIDC //////////
#[get("/login")]
async fn login(
    client: web::Data<Arc<Mutex<Client<openid::Discovered, StandardClaims>>>>,
    query: web::Query<AuthRequest>,
) -> impl Responder {
    let client = client.lock().await;

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
