use crate::config::actix::AppState;
use crate::AuthRequest;
use actix_session::Session;
use actix_web::web::Data;
use actix_web::{get, web, HttpResponse, Responder};
use openid::{Client, Discovered, StandardClaims};
use reqwest::Client as HttpClient;
use serde::{Deserialize, Serialize};
use std::error::Error;

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

// Structure for the token response
#[derive(Deserialize)]
struct TokenResponse {
    access_token: String,
    token_type: String,
    expires_in: u64,
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

    let admin_token = get_access_token(&client).await.unwrap();

    println!("Admin token: {}", admin_token);

    // Define the new user
    let new_user = NewUser {
        username: "new_user".to_string(),
        email: "new_user@example.com".to_string(),
        enabled: true,
        credentials: vec![Credential {
            r#type: "password".to_string(),
            value: "SecureP@ssw0rd!".to_string(),
            temporary: false, // Set to true if the user must change the password on first login
        }],
    };

    let api_url = format!(
        "http://localhost:8181/admin/realms/{}/users", // Replace with your IdP's API
        "Banana"
    );

    let http_client = HttpClient::new();

    let response = http_client
        .post(&api_url)
        .bearer_auth(admin_token) // Use the admin token for authorization
        .json(&new_user)          // Serialize the user data into JSON
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

// Function to fetch an access token using the client credentials flow
async fn get_access_token(client: &Client<Discovered, StandardClaims>) -> Result<String, Box<dyn Error>> {
    let token_endpoint = client
        .config()
        .token_endpoint
        .clone();

    let http_client = HttpClient::new();

    let response = http_client
        .post(token_endpoint)
        .basic_auth(client.client_id.clone(), client.client_secret.clone())
        .form(&[("grant_type", "client_credentials")])
        .send()
        .await?
        .json::<TokenResponse>()
        .await?;

    Ok(response.access_token)
}