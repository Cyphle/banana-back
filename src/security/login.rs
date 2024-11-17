use std::collections::HashMap;
use actix_session::Session;
use actix_web::{get, web, HttpResponse, Responder};
use log::{error, info};
use openid::{Client, Options};
use crate::AuthRequest;
use crate::config::actix::AppState;

// TODO clean all this mess and test it
#[get("/login")]
async fn login(
    session: Session,
    state: web::Data<AppState>,
    query: web::Query<AuthRequest>,
) -> impl Responder {
    let client = state.oidc_client.as_ref().unwrap().lock().unwrap();

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
                        Err(e) => error!("Failed to save token in session: {}", e)
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