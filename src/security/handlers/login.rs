use crate::config::actix::AppState;
use crate::AuthRequest;
use actix_session::Session;
use actix_web::web::Data;
use actix_web::{get, web, HttpResponse, Responder};
use log::{error, info};
use openid::{Bearer, Client, Options, Token};
use std::collections::HashMap;
use std::sync::{Arc, Mutex};

#[get("/login")]
async fn login(
    session: Session,
    state: Data<AppState>,
    query: web::Query<AuthRequest>,
) -> impl Responder {
    let client = state.oidc_client.as_ref().unwrap().lock().unwrap();

    match &query.code {
        Some(authorization_code) => {
            info!(
                "Requesting token with received authorization code: {}",
                authorization_code
            );

            match client.authenticate(authorization_code, None, None).await {
                Ok(token) => {
                    let access_token = token.bearer.access_token.clone();
                    let id_token = token.bearer.id_token.clone();

                    save_in_session(session, &token);

                    // TODO il faut retourner autre chose
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

            // TODO Il faut définir un nonce et max age ici pour réutiliser dans la méthode authenticate
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

fn save_in_session(session: Session, token: &Token) {
    let saved = session.insert(
        "user_id",
        token.bearer.clone(),
    );
    match saved {
        Ok(_) => info!("Token saved in session"),
        Err(e) => error!("Failed to save token in session: {}", e)
    }
}