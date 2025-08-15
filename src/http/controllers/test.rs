use crate::config::actix::AppState;
use crate::http::adapters::profile::get_profile_by_username;
use crate::repositories::profile::find_one_by_username;
use actix_web::{get, web, HttpResponse, Responder};
use log::{error, info};

#[get("/test")]
async fn test(state: web::Data<AppState>) -> impl Responder {
    match get_profile_by_username(&state.db_connection, "test").await {
        Some(username) => {
            info!("Found profile for username: {}", username.username);
            HttpResponse::Ok().json(username)
        }
        None => {
            error!("Nothing found");
            HttpResponse::Ok().finish()
        }
    }
}