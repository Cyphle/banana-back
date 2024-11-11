use openid::{Client, StandardClaims};

pub struct OidcConfig {
    pub issuer_url: String,
    pub client_id: String,
    pub client_secret: String,
    pub redirect_uri: String,
}

impl OidcConfig {
    pub fn new(issuer_url: String, client_id: String, client_secret: String, redirect_uri: String) -> Self {
        OidcConfig { 
            issuer_url, 
            client_id, 
            client_secret, 
            redirect_uri 
        }
    }
}

pub async fn get_client(config: &OidcConfig) -> Client<openid::Discovered, StandardClaims> {
    Client::discover(
        config.client_id.to_string(),
        config.client_secret.to_string(),
        Some(config.redirect_uri.to_string()),
        reqwest::Url::parse(&config.issuer_url).unwrap(),
    )
    .await
    .expect("Failed to discover OpenID configuration")
}