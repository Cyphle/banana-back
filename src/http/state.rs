use sea_orm::DatabaseConnection;

pub struct HandlerState {
    pub db_connection: &'static DatabaseConnection
}