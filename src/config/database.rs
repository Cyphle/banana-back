use sea_orm::{Database, DatabaseConnection, DbErr};

struct DatabaseConfig {
    host: &'static str,
    port: &'static str,
    schema: &'static str,
    username: &'static str,
    password: &'static str,
    max_connections: u32,
    min_connections: u32,
    connect_timeout: u64,
    acquire_timeout: u64,
    idle_timeout: u64,
    max_lifetime: u64,
    sqlx_logging: bool,
}

fn new_database_config() -> DatabaseConfig {
    DatabaseConfig {
        host: "localhost",
        port: "5434",
        schema: "banana",
        username: "postgres",
        password: "postgres",
        max_connections: 100,
        min_connections: 5,
        connect_timeout: 8,
        acquire_timeout: 8,
        idle_timeout: 8,
        max_lifetime: 8,
        sqlx_logging: true,
    }
}

pub async fn connect() -> Result<DatabaseConnection, DbErr> {
    let config = new_database_config();

    // let mut opt = ConnectOptions::new("protocol://".to_string() + config.username + ":" + config.password + "@" + config.host + ":" + config.port + "/" + config.schema);
    // opt.max_connections(config.max_connections)
    //     .min_connections(config.min_connections)
    //     .connect_timeout(Duration::from_secs(config.connect_timeout))
    //     .acquire_timeout(Duration::from_secs(config.acquire_timeout))
    //     .idle_timeout(Duration::from_secs(config.idle_timeout))
    //     .max_lifetime(Duration::from_secs(config.max_lifetime))
    //     .sqlx_logging(config.sqlx_logging)
    //     .sqlx_logging_level(log::LevelFilter::Info)
    //     .set_schema_search_path(config.schema);
    //
    // Database::connect(opt).await

    Database::connect("postgres://".to_string() + config.username + ":" + config.password + "@" + config.host + ":" + config.port + "/" + config.schema).await
}