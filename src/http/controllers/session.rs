use actix_session::Session;
use actix_web::{get, web, HttpResponse, Responder};
use log::{error, info};
use openid::{Bearer, Client, StandardClaims, Token, TokenIntrospection};
use crate::config::local::oidc_config::USER_SESSION_KEY;
use crate::AuthRequest;
use crate::config::actix::AppState;
use crate::security::token::{get_username_from_bearer};

#[get("/get-from-session")]
async fn get_session(
    session: Session,
    state: web::Data<AppState>,
    _: web::Query<AuthRequest>,
) -> impl Responder {
    let user_id = session.get::<Bearer>(USER_SESSION_KEY);
    let client = state.oidc_client.as_ref().unwrap().lock().unwrap();

    match user_id {
        Ok(user_id) => match user_id {
            Some(user_id) => {
                println!("Bearer in session {:?}", user_id);

                // Example of retrieving username
                match get_username_from_bearer(&client, &user_id).await {
                    Some(userinfo) => {
                        info!("User info found in login: {:?}", userinfo);
                    }
                    None => {
                        error!("No user info found");
                    }
                }

                info!("Session data found: {:?}", user_id);
                // info!("user info: {:?}", res);
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