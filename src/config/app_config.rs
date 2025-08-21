use config::{Config, Environment, File};
use serde::Deserialize;
use std::env;
use std::time::Duration;
use crate::config::database::DatabaseConfig;
use crate::config::session::SessionConfig;

// TODO à merge avec ce qui est dans le dossier security
#[derive(Debug, Deserialize, Clone)]
pub struct OidcAdminConfig {
    pub client: OidcClientConfig,
    pub create_user_url: String,
}

#[derive(Debug, Deserialize, Clone)]
pub struct OidcRealmConfig {
    pub url: String,
    pub name: String,
}

#[derive(Debug, Deserialize, Clone)]
pub struct OidcClientConfig {
    pub id: String,
    pub secret: String
}

#[derive(Debug, Deserialize, Clone)]
pub struct OidcUrlConfig {
    pub redirect: String,
    pub logout: String,
}

// TODO à merge avec ce qui est dans le dossier security
#[derive(Debug, Deserialize, Clone)]
pub struct OidcConfig {
    pub realm: OidcRealmConfig,
    pub url: OidcUrlConfig,
    pub client: OidcClientConfig,
    pub nonce: String,
    pub session_timeout_minutes: i64,
    pub admin: OidcAdminConfig,
}

#[derive(Debug, Deserialize, Clone)]
pub struct CorsConfig {
    pub allowed_origin: String,
    pub allowed_methods: Vec<String>,
    pub allowed_headers: Vec<String>,
    pub supports_credentials: bool,
    pub max_age: u64,
}

#[derive(Debug, Deserialize, Clone)]
pub struct AppConfig {
    pub app: AppServerConfig,
    pub database: DatabaseConfig,
    pub oidc: OidcConfig,
    pub session: SessionConfig,
    pub cors: CorsConfig,
    pub logging: LoggingConfig,
}

#[derive(Debug, Deserialize, Clone)]
pub struct AppServerConfig {
    pub port: u16,
    pub host: String,
}

#[derive(Debug, Deserialize, Clone)]
pub struct LoggingConfig {
    pub level: String,
}

impl AppConfig {
    pub fn new() -> Result<Self, config::ConfigError> {
        let config = Config::builder()
            .add_source(File::with_name("config/default"))
            .add_source(File::with_name("config/local").required(false))
            .add_source(Environment::with_prefix("BANANA").separator("_"))
            .build()?;

        config.try_deserialize()
    }

    pub fn database_url(&self) -> String {
        format!(
            "postgres://{}:{}@{}:{}/{}",
            self.database.username,
            self.database.password,
            self.database.host,
            self.database.port,
            self.database.schema
        )
    }

    pub fn session_ttl_duration(&self) -> Duration {
        Duration::from_secs(self.session.session_ttl_days * 24 * 60 * 60)
    }
}

// Constantes pour les clés de session
pub const USER_SESSION_KEY: &str = "user_session";
