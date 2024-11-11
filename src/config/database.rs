use std::time::Duration;
use sea_orm::{ConnectOptions, Database, DatabaseConnection, DbErr};

pub struct DatabaseConfig {
    pub host: &'static str,
    pub port: &'static str,
    pub schema: &'static str,
    pub username: &'static str,
    pub password: &'static str,
    pub max_connections: u32,
    pub min_connections: u32,
    pub connect_timeout: u64,
    pub acquire_timeout: u64,
    pub idle_timeout: u64,
    pub max_lifetime: u64,
    pub sqlx_logging: bool,
}

pub async fn connect(config: &DatabaseConfig) -> Result<DatabaseConnection, DbErr> {

    let mut opt = ConnectOptions::new("postgres://".to_string() + config.username + ":" + config.password + "@" + config.host + ":" + config.port + "/" + config.schema);
    opt.max_connections(config.max_connections)
        .min_connections(config.min_connections)
        .connect_timeout(Duration::from_secs(config.connect_timeout))
        .acquire_timeout(Duration::from_secs(config.acquire_timeout))
        .idle_timeout(Duration::from_secs(config.idle_timeout))
        .max_lifetime(Duration::from_secs(config.max_lifetime))
        .sqlx_logging(config.sqlx_logging)
        .sqlx_logging_level(log::LevelFilter::Info);

    Database::connect(opt).await
}