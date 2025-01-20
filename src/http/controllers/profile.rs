use crate::config::actix::AppState;
use crate::domain::profile::CreateProfileCommand;
use crate::http::adapters::profile::get_profile_by_username;
use crate::repositories;
use crate::security::token::get_username_from_session;
use actix_session::Session;
use actix_web::{get, post, web, HttpResponse, Responder};
use log::{error, info};

#[get("/profiles")]
async fn get_profile(session: Session, state: web::Data<AppState>) -> impl Responder {
    let client = state.oidc_client.as_ref().unwrap().lock().unwrap();
    match get_username_from_session(&client, &session).await {
        Some(username) => {
            info!("Username in session: {:?}", username);
            match get_profile_by_username(&state.db_connection, &username).await {
                Some(profile) => HttpResponse::Ok().json(profile),
                None => {
                    error!("No profile found for username {}", username);
                    HttpResponse::InternalServerError().finish()
                }
            }
        }
        None => {
            error!("No username info found in session");
            HttpResponse::Ok().finish()
        }
    }
}
