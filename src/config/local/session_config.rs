use crate::config::actix::SessionConfig;

// TODO il faut construire une connexion string complète
pub fn get_session_config() -> SessionConfig {
    SessionConfig {
        store_addr: "redis://127.0.0.1:6379".to_string(),
        cookie_name: "actix_cookie".to_string(),
    }
}