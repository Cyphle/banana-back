use actix_session::Session;
use actix_web::{get, web, HttpResponse, Responder};
use actix_web::web::Data;
use log::{error, info};
use openid::{Bearer, Client, Discovered, StandardClaims};
use url::Url;
use crate::AuthRequest;
use crate::config::actix::AppState;
use crate::config::local::oidc_config::USER_SESSION_KEY;

#[get("/logout")]
async fn logout(
    session: Session,
    state: Data<AppState>,
    _: web::Query<AuthRequest>,
) -> impl Responder {
    let user_session = session.get::<Bearer>(USER_SESSION_KEY);
    let client = state.oidc_client.as_ref().unwrap().lock().unwrap();
    let logout_uri: &str = state.oidc_config.logout_uri.as_ref();

    match user_session {
        Ok(bearer_opt) => match bearer_opt {
            Some(bearer) => {
                match build_logout_url(&client, &bearer.clone().id_token.unwrap(), logout_uri).await {
                    Ok(logout_url) => {
                        info!("Logout URL: {}", logout_url);

                        match reqwest::get(logout_url).await {
                            Ok(response) => {
                                // Print the response body as text
                                let body = response.text().await;

                                match body {
                                    Ok(body) => println!("Response: {}", body),
                                    Err(e) => error!("Error reading response body: {}", e),
                                }

                                // TODO and delete session
                                HttpResponse::NoContent()
                                    .append_header(("Location ", client.redirect_url()))
                                    .finish()
                            }
                            Err(e) => {
                                error!("Error generating logout URL: {}", e);
                                HttpResponse::InternalServerError().body(format!("Error generating logout URL: {}", e))
                            }
                        }
                    }
                    Err(e) => {
                        error!("Error generating logout URL: {}", e);
                        HttpResponse::InternalServerError().body(format!("Error generating logout URL: {}", e))
                    }
                }
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

pub async fn build_logout_url(
    client: &Client<Discovered, StandardClaims>,
    id_token: &str,
    logout_uri: &str,
) -> Result<Url, Box<dyn std::error::Error>> {
    // Access the discovered metadata
    let discovered_metadata = client.config();

    // Extract the end_session_endpoint
    let mut end_session_endpoint = discovered_metadata
        .end_session_endpoint
        .as_ref()
        .ok_or("End session endpoint not available in metadata")?
        .clone();

    end_session_endpoint
        .query_pairs_mut()
        .append_pair("id_token_hint", id_token)
        .append_pair("post_logout_redirect_uri", logout_uri);

    Ok(end_session_endpoint.clone())
}
