use crate::config::local::oidc_config::USER_SESSION_KEY;
use log::{error, info, warn};
use openid::error::Error;
use openid::{Bearer, DiscoveredClient, IdToken, StandardClaims, Token, TokenIntrospection, Userinfo};

// To get the username from a Bearer token
pub async fn get_username_from_bearer(client: &DiscoveredClient, bearer: &Bearer) -> Option<String> {
    info!("Bearer token found in session {:?}", bearer.clone());
    let token_wrapper: Token<StandardClaims> = Token::from(bearer.clone());
    match client.request_token_introspection::<TokenIntrospection<StandardClaims>>(&token_wrapper).await {
        Ok(userinfo) => {
            userinfo.username
        }
        Err(e) => {
            warn!("No user info found: {}", e);
            None
        }
    }
}

// To get the username from the bearer token in session
pub async fn get_username_from_session(client: &DiscoveredClient, session: &actix_session::Session) -> Option<String> {
    let bearer = session.get::<Bearer>(USER_SESSION_KEY);
    // TODO faut pas faire unwrap unwrap
    get_username_from_bearer(client, &bearer.unwrap().unwrap()).await
}
