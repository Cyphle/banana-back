use chrono::Duration;

use crate::security::oidc::{OidcAdminConfig, OidcConfig};

pub fn get_oidc_config() -> OidcConfig {
    return OidcConfig::new(
        "http://localhost:8181/realms/Banana".to_string(),
        "banana".to_string(),
        "banana-secret".to_string(),
        "http://localhost:8080/login".to_string(),
        Some("nonce".to_string()),
        Some(Duration::minutes(10)),
        "http://localhost:8080/logout".to_string(),
        OidcAdminConfig {
            client_id: "banana-admin".to_string(),
            client_secret: "5YMRPEgEmwq6G819T98F4dMhb1vMx7AR".to_string(),
            create_user_url: "http://localhost:8181/admin/realms/Banana/users".to_string(),
        }
    )
}

pub static USER_SESSION_KEY: &str = "user_session";