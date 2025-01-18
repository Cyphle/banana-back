use crate::config::local::oidc_config::USER_SESSION_KEY;
use actix_session::SessionGetError;
use log::{error, info, warn};
use openid::error::Error;
use openid::{
    Bearer, Client, Discovered, DiscoveredClient, IdToken, StandardClaims, Token,
    TokenIntrospection, Userinfo,
};
use reqwest::Client as HttpClient;
use serde::Deserialize;
use crate::security::oidc::OidcAdminConfig;

// To get the username from a Bearer token
pub async fn get_username_from_bearer(
    client: &DiscoveredClient,
    bearer: &Bearer,
) -> Option<String> {
    info!("Bearer token found in session {:?}", bearer.clone());
    let token_wrapper: Token<StandardClaims> = Token::from(bearer.clone());
    match client
        .request_token_introspection::<TokenIntrospection<StandardClaims>>(&token_wrapper)
        .await
    {
        Ok(userinfo) => userinfo.username,
        Err(e) => {
            warn!("No user info found: {}", e);
            None
        }
    }
}

// To get the username from the bearer token in session
pub async fn get_username_from_session(
    client: &DiscoveredClient,
    session: &actix_session::Session,
) -> Option<String> {
    let bearer = session.get::<Bearer>(USER_SESSION_KEY);

    match bearer {
        Ok(bearer) => match bearer {
            Some(bearer) => get_username_from_bearer(client, &bearer).await,
            None => {
                error!("No bearer token found in session");
                None
            }
        },
        Err(e) => {
            error!("Error getting bearer token from session: {}", e);
            None
        }
    }
}

// Structure for the token response
#[derive(Deserialize)]
struct TokenResponse {
    access_token: String,
    token_type: String,
    expires_in: u64,
}

// To get an access token for admin account (client)
pub async fn get_admin_access_token(
    client: &Client<Discovered, StandardClaims>,
    admin: &OidcAdminConfig
) -> Result<String, Box<dyn std::error::Error>> {
    let token_endpoint = client.config().token_endpoint.clone();
    let token_request = HttpClient::new()
        .post(token_endpoint)
        .basic_auth(admin.client_id.clone(), Some(admin.client_secret.clone()))
        .form(&[("grant_type", "client_credentials")])
        .send()
        .await;

    match token_request {
        Ok(response) => {
            match response.json::<TokenResponse>().await {
                Ok(response) => Ok(response.access_token),
                Err(e) => {
                    error!("Error getting access token: {}", e);
                    Err(Box::new(e))
                }
            }
        }
        Err(e) => {
            error!("Error getting access token: {}", e);
            Err(Box::new(e))
        }
    }
}
