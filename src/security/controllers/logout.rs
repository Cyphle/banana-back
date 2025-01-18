use openid::{Client, Discovered, StandardClaims};
use url::Url;

pub async fn logout(
    client: &Client<Discovered, StandardClaims>,
    id_token: &str,
    post_logout_redirect_uri: &str,
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
        .append_pair("post_logout_redirect_uri", post_logout_redirect_uri);

    Ok(end_session_endpoint.clone())
}
