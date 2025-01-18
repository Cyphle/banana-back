use crate::config::actix::AppState;
use crate::{repositories, AuthRequest};
use actix_session::Session;
use actix_web::web::Data;
use actix_web::{get, post, web, HttpResponse, Responder};
use openid::{Client, Discovered, StandardClaims};
use reqwest::Client as HttpClient;
use serde::{Deserialize, Serialize};
use std::error::Error;
use log::{error, info};
use crate::domain::profile::CreateProfileCommand;
use crate::dto::requests::profile::CreateProfileRequest;
use crate::security::token::get_admin_access_token;

#[derive(Serialize)]
struct KeycloakUser {
    username: String,
    email: String,
    enabled: bool,
    credentials: Vec<KeycloakCredential>,
}

#[derive(Serialize)]
struct KeycloakCredential {
    r#type: String,
    value: String,
    temporary: bool,
}

#[post("/register")]
pub async fn register(
    payload: web::Json<CreateProfileRequest>,
    _: Session,
    state: Data<AppState>,
    _: web::Query<AuthRequest>,
) -> impl Responder {
    let client = state.oidc_client.as_ref().unwrap().lock().unwrap();
    let admin_token = get_admin_access_token(&client, &state.oidc_config.admin).await.unwrap();

    let moved_payload = payload.into_inner();

    match repositories::profile::create(
        &state.db_connection,
        &CreateProfileCommand {
            username: moved_payload.username.to_owned(),
            email: moved_payload.email.to_owned(),
            first_name: moved_payload.first_name.to_owned(),
            last_name: moved_payload.last_name.to_owned(),
        },
    ).await {
        Ok(profile) => {
            let new_user = KeycloakUser {
                username: profile.username.clone(),
                email: profile.email.clone(),
                enabled: true,
                credentials: vec![KeycloakCredential {
                    r#type: "password".to_string(),
                    value: "Bonjour".to_string(),
                    temporary: false,
                }],
            };

            match HttpClient::new()
                .post(&state.oidc_config.admin.create_user_url)
                .bearer_auth(admin_token)
                .json(&new_user)
                .send()
                .await {
                Ok(response) => {
                    info!("User created in keycloak: {:?}", response);
                    HttpResponse::Created().finish()
                }
                Err(e) => {
                    println!("Error creating user in keycloak: {:?}", e);
                    HttpResponse::InternalServerError().finish()
                }
            }
        },
        Err(e) => {
            error!("Error creating user: {:?}", e);
            HttpResponse::InternalServerError().finish()
        },
    }
}
