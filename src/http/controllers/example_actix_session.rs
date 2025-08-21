use actix_session::Session;
use actix_web::{get, web, HttpResponse, Responder};
use log::{error, info};
use openid::{Bearer, Client, StandardClaims, Token, TokenIntrospection};
use crate::config::local::oidc_config::USER_SESSION_KEY;
use crate::config::actix::AppState;
use crate::security::controllers::auth_request::AuthRequest;
use crate::security::controllers::logout::logout;
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


                // LOGOUT TEST
                // URL to redirect the user after logout
                // let post_logout_redirect_uri = "http://localhost:9000";
                //
                // // Generate the logout URL
                // match logout(&client, &user_id.clone().id_token.unwrap(), post_logout_redirect_uri).await {
                //     Ok(logout_url) => {
                //         info!("Logout URL: {}", logout_url);
                //     }
                //     Err(e) => {
                //         error!("Error generating logout URL: {}", e);
                //     }
                // }

                // END LOGOUT TEST

                // Example of retrieving username
                // match get_username_from_bearer(&client, &user_id).await {
                //     Some(userinfo) => {
                //         info!("User info found in login: {:?}", userinfo);
                //     }
                //     None => {
                //         error!("No user info found");
                //     }
                // }

                // info!("Session data found: {:?}", user_id);
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
#[get("/delete-from-session")]
async fn delete_session(
    session: Session,
    state: web::Data<AppState>,
    _: web::Query<AuthRequest>,
) -> impl Responder {
    let before = session.get::<Bearer>(USER_SESSION_KEY);
    info!("session before delete: {:?}", before);
    session.remove(USER_SESSION_KEY);
    let user_id = session.get::<Bearer>(USER_SESSION_KEY);
    info!("session after delete: {:?}", user_id);
    HttpResponse::Ok().body("Session deleted")
}