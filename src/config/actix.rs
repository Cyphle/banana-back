use crate::http::handlers::profile::{create_profile, get_profile_by_id};
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
            .service(get_profile_by_id)
            .service(create_profile)
    })
        .bind(("127.0.0.1", 8080))?
        .run()
        .await
}