use chrono::Duration;

use crate::security::oidc::OidcConfig;

pub fn get_oidc_config() -> OidcConfig {
    return OidcConfig::new(
        "http://localhost:8181/realms/Banana".to_string(),
        "banana".to_string(),
        "banana-secret".to_string(),
        "http://localhost:9000/login".to_string(),
        Some("nonce".to_string()),
        Some(Duration::minutes(10)),
        "http://localhost:9000".to_string(),
    )
}

pub static USER_SESSION_KEY: &str = "user_session";