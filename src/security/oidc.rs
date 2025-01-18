use chrono::Duration;
use openid::{Client, StandardClaims};

pub struct OidcConfig {
    pub issuer_url: String,
    pub client_id: String,
    pub client_secret: String,
    pub redirect_uri: String,
    pub nonce: Option<String>,
    pub max_age: Option<Duration>,
    pub logout_uri: String,
    pub admin: OidcAdminConfig,
}
impl Clone for OidcConfig {
    fn clone(&self) -> Self {
        Self {
            issuer_url: self.issuer_url.clone(),
            client_id: self.client_id.clone(),
            client_secret: self.client_secret.clone(),
            redirect_uri: self.redirect_uri.clone(),
            nonce: self.nonce.clone(),
            max_age: self.max_age.clone(),
            logout_uri: self.logout_uri.clone(),
            admin: self.admin.clone(),
        }
    }
}

pub struct OidcAdminConfig {
    pub client_id: String,
    pub client_secret: String,
    pub create_user_url: String
}

impl Clone for OidcAdminConfig {
    fn clone(&self) -> Self {
        Self {
            client_id: self.client_id.clone(),
            client_secret: self.client_secret.clone(),
        }
    }
}

impl OidcConfig {
    pub fn new(
        issuer_url: String,
        client_id: String,
        client_secret: String,
        redirect_uri: String,
        nonce: Option<String>,
        max_age: Option<Duration>,
        logout_uri: String,
        admin: OidcAdminConfig,
    ) -> Self {
        OidcConfig {
            issuer_url,
            client_id,
            client_secret,
            redirect_uri,
            nonce,
            max_age,
            logout_uri,
            admin,
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
