use log::{error, warn};
use openid::{Bearer, DiscoveredClient, IdToken, StandardClaims, Token, TokenIntrospection, Userinfo};
use openid::error::Error;

// To get the username from a Bearer token
pub async fn get_username_from_bearer(client: &DiscoveredClient, bearer: &Bearer) -> Option<String> {
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