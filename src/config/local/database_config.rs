use crate::config::database::DatabaseConfig;

pub fn get_database_config() -> DatabaseConfig {
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