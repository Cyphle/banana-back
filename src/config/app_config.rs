use config::{Config, Environment, File};
use serde::Deserialize;
use std::time::Duration;

#[derive(Debug, Deserialize, Clone)]
pub struct DatabaseConfig {
    pub host: String,
    pub port: String,
    pub schema: String,
    pub username: String,
    pub password: String,
    pub max_connections: u32,
    pub min_connections: u32,
    pub connect_timeout: u64,
    pub acquire_timeout: u64,
    pub idle_timeout: u64,
    pub max_lifetime: u64,
    pub sqlx_logging: bool,
}

#[derive(Debug, Deserialize, Clone)]
pub struct DbConfig {
    pub name: String,
}

#[derive(Debug, Deserialize, Clone)]
pub struct OidcAdminConfig {
    pub client_id: String,
    pub client_secret: String,
    pub create_user_url: String,
}

#[derive(Debug, Deserialize, Clone)]
pub struct OidcConfig {
    pub realm_url: String,
    pub realm: String,
    pub client_secret: String,
    pub login_url: String,
    pub nonce: String,
    pub session_timeout_minutes: i64,
    pub base_url: String,
    pub admin: OidcAdminConfig,
}

#[derive(Debug, Deserialize, Clone)]
pub struct SessionConfig {
    pub store_addr: String,
    pub cookie_name: String,
    pub session_ttl_days: u64,
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
    pub db: DbConfig,
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
            .add_source(Environment::with_prefix("BANANA"))
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
