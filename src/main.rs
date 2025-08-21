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

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    config::logger::config();

    match config::app_config::AppConfig::new() {
        Ok(config) => {
            info!("Starting the application");
            config::actix::start(&config).await
        }
        Err(e) => {
            log::error!("Erreur lors du chargement de la configuration: {:?}", e);
            return Err(std::io::Error::new(std::io::ErrorKind::Other, "Configuration error"));
        }
    }
    
    // TODO faudra trouver un moyen de close la connexion. Mais là on peut pas move la static_db
    // TODO en fait il faut plutôt clone plutôt qu'un pointeur. et même sans ça, la connexion se close bien
}
