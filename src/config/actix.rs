use crate::http::handlers::examples::hello;
use crate::http::state::HandlerState;
use actix_web::{web, App, HttpServer};
use log::info;
use sea_orm::DatabaseConnection;

pub async fn config(db_connection: &'static DatabaseConnection) -> std::io::Result<()> {
    info!("Starting Actix server");

    HttpServer::new(|| {
        App::new()
            .app_data(web::Data::new(HandlerState {
                db_connection: db_connection
            }))
            .service(hello)
    })
        .bind(("127.0.0.1", 8080))?
        .run()
        .await
}