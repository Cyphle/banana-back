use crate::config::actix::AppState;
use crate::AuthRequest;
use actix_session::Session;
use actix_web::web::Data;
use actix_web::{get, web, HttpResponse, Responder};
use openid::{Client, Discovered, StandardClaims};
use reqwest::Client as HttpClient;
use serde::{Deserialize, Serialize};
use std::error::Error;
use crate::security::token::get_admin_access_token;

#[derive(Serialize)]
struct NewUser {
    username: String,
    email: String,
    enabled: bool,
    credentials: Vec<Credential>,
}

// Structure for the credentials (password)
#[derive(Serialize)]
struct Credential {
    r#type: String,
    value: String,
    temporary: bool,
}



// TODO
/*
    J'utilise les mêmes client id et client secret que l'app cliente. C'est pas bon il faut un autre dédié admin
    A noter, qu'il faut set:
    - "service accounts roles" on dans settings du client
    - que ce soit un client confidential
    - client authentication on dans settings
    - dans les rôles, il faut ajouter "manage-users" pour le client
 */
#[get("/register")]
pub async fn register(
    session: Session,
    state: Data<AppState>,
    query: web::Query<AuthRequest>,
) -> impl Responder {
    let client = state.oidc_client.as_ref().unwrap().lock().unwrap();

    let admin_token = get_admin_access_token(&client, &state.oidc_config.admin).await.unwrap();

    println!("Admin token: {}", admin_token);

    // Define the new user
    let new_user = NewUser {
        username: "new_user".to_string(),
        email: "new_user@example.com".to_string(),
        enabled: true,
        credentials: vec![Credential {
            r#type: "password".to_string(),
            value: "SecureP@ssw0rd!".to_string(),
            temporary: false,
        }],
    };

    let response = HttpClient::new()
        .post(&state.oidc_config.admin.create_user_url)
        .bearer_auth(admin_token)
        .json(&new_user)
        .send()
        .await;

    match response {
        Ok(response) => {
            println!("User created successfully {:?}", response);
        }
        Err(e) => {
            println!("Error creating user: {:?}", e);
        }
    }

    HttpResponse::Created().finish()
}
