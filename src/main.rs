use actix_web::cookie::{time, Key};
use actix_web::{get, web, App, HttpResponse, HttpServer, Responder};
use log::info;
use openid::{Client, Options, StandardClaims, Token};
use serde::{Deserialize, Serialize};

mod config;
mod domain;
mod dto;
mod http;
mod repositories;
mod security;

#[derive(Debug, Deserialize)]
struct AuthRequest {
    code: Option<String>,
    session_state: Option<String>,
    iss: Option<String>,
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    config::logger::config();

    // Test de la nouvelle configuration
    match config::app_config::AppConfig::new() {
        Ok(config) => {
            info!("Configuration chargée avec succès!");
            info!("Port: {}", config.app.port);
            info!("Host: {}", config.app.host);
            info!("Base de données: {}:{}", config.database.host, config.database.port);
            info!("OIDC Realm: {}", config.oidc.realm);
            info!("Session cookie: {}", config.session.cookie_name);
        }
        Err(e) => {
            log::error!("Erreur lors du chargement de la configuration: {:?}", e);
            return Err(std::io::Error::new(std::io::ErrorKind::Other, "Configuration error"));
        }
    }

    // info!("Starting the application");
    // config::actix::config().await

    Result::Ok(())
    // TODO faudra trouver un moyen de close la connexion. Mais là on peut pas move la static_db
}
