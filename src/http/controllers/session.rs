use actix_session::Session;
use actix_web::{get, web, HttpResponse, Responder};
use log::{error, info};
use openid::{Bearer, Client};
use crate::config::local::oidc_config::USER_SESSION_KEY;
use crate::AuthRequest;
use crate::config::actix::AppState;

#[get("/get-from-session")]
async fn get_session(
    session: Session,
    _: web::Data<AppState>,
    _: web::Query<AuthRequest>,
) -> impl Responder {
    let user_id = session.get::<Bearer>(USER_SESSION_KEY);

    match user_id {
        Ok(user_id) => match user_id {
            Some(user_id) => {
                info!("Session data found: {:?}", user_id);
                HttpResponse::Ok().json(user_id)
            }
            None => {
                error!("No session data found");
                HttpResponse::Ok().body("No session data found")
            }
        },
        Err(e) => {
            error!("No session data found: {}", e);
            HttpResponse::Ok().body(format!("No session data found: {}", e))
        }
    }
}