use std::ops::{Deref, DerefMut};
use sea_orm::DatabaseConnection;
use crate::domain::profile::CreateProfileCommand;

mod config;
mod repositories;
mod domain;
mod http;
mod dto;

#[actix_web::main]
async fn main() {
    config::logger::config();

    log::info!("Starting the application");

    let db = config::database::connect().await.unwrap();
    let static_db = Box::leak(Box::new(db));

    config::actix::config(static_db).await;

    // repositories::profiles::create(static_db, &CreateProfileCommand {
    //     username: "johndoe".to_string(),
    //     email: "johndoe".to_string(),
    //     first_name: "John".to_string(),
    //     last_name: "Doe".to_string(),
    // }).await.unwrap();

    log::info!("Application is now closed");

    // TODO faudra trouver un moyen de close la connexion. Mais là on peut pas move la static_db
}
