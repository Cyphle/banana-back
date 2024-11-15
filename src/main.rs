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

    info!("Starting the application");
    config::actix::config().await

    // TODO faudra trouver un moyen de close la connexion. Mais là on peut pas move la static_db
}
