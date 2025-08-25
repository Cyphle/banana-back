use crate::http::controllers::account::{create_account, find_all, find_one};
use crate::http::controllers::example_actix_session::{delete_session, get_session};
use crate::http::controllers::example_actix_store::{add_to_store, get_from_store};
use crate::http::controllers::technical::{live, ready};
use crate::http::controllers::test::test;
use crate::security::controllers::login::login;
use crate::security::controllers::logout::logout;
use crate::security::controllers::register::register;
use crate::security::oidc::OidcConfig;
use crate::{
    config::{
        database::connect,
        local::{
            oidc_config::get_oidc_config,
            session_config::get_session_config,
        },
    },
    http::controllers::profile::get_profile,
    security::oidc::get_client,
};
use actix_cors::Cors;
use actix_session::storage::CookieSessionStore;
use actix_session::{config::PersistentSession, storage::RedisSessionStore, SessionMiddleware};
use actix_web::{
    cookie::{time, Key},
    web, App, HttpServer,
};
use log::info;
use openid::{Bearer, Client, Discovered, StandardClaims};
use sea_orm::DatabaseConnection;
use std::sync::Arc;
use std::task::ready;
use std::{collections::HashMap, sync::Mutex};
use crate::config::app_config::AppConfig;

pub struct AppState {
    pub db_connection: &'static DatabaseConnection,
    pub oidc_client: Option<Arc<Mutex<Client<Discovered, StandardClaims>>>>,
    pub store: Mutex<HashMap<String, Bearer>>,
    pub oidc_config: Option<OidcConfig>,
}

pub struct SessionConfig {
    pub store_addr: String,
    pub cookie_name: String,
}

// TODO il faut remplacer ça par l'utilisation de app_conffig.rs
pub async fn start(application_configuration: &AppConfig) -> std::io::Result<()> {
    info!("Starting Actix server");

    // TODO faut sortir pas mal de chose là et les passer en paramètres

    // Databse
    // TODO à priori, chatgpt conseille de faire un clone à envoyer à actix au lieu de faire un pointer box
    let db = connect(&application_configuration.database).await.unwrap();
    let static_db = Box::leak(Box::new(db));

    // OIDC
    // TODO Remettre l'oidc et bien régler la config
    // let oidc_config = get_oidc_config();
    // let client = Arc::new(Mutex::new(get_client(&oidc_config).await));

    // Session
    let session_config = get_session_config();
    let secret_key = Key::from(&[0; 64]);
    // let redis_store = RedisSessionStore::new(session_config.store_addr)
    let redis_store = RedisSessionStore::new(application_configuration.get_session_url())
        .await
        .unwrap();

    let state = web::Data::new(AppState {
        db_connection: static_db,
        // oidc_client: Some(client.clone()),
        oidc_client: None,
        store: Mutex::new(HashMap::new()),
        // oidc_config: oidc_config.clone(),
        oidc_config: None,
    });

    // Actix
    HttpServer::new(move || {
        App::new()
            .wrap(
                Cors::default()
                    .allowed_origin("http://localhost:9000") // Change to your frontend URL
                    .allowed_methods(vec!["GET", "POST"])
                    .allowed_headers(vec![actix_web::http::header::CONTENT_TYPE])
                    .supports_credentials() // Optional, if credentials are used
                    .max_age(3600),
            )
            .wrap(
                SessionMiddleware::builder(redis_store.clone(), secret_key.clone())
                    .session_lifecycle(
                        PersistentSession::default().session_ttl(time::Duration::days(5)),
                    )
                    .cookie_secure(false)
                    .cookie_name(session_config.cookie_name.clone())
                    .build(),
            )
            .app_data(state.clone())
            // Tests and examples endpoints
            .service(add_to_store)
            .service(get_from_store)
            .service(get_session)
            .service(delete_session)
            // End tests and examples
            .service(login)
            .service(logout)
            .service(register)
            .service(get_profile)
            .service(create_account)
            .service(find_one)
            .service(find_all)
            .service(test)
            // Technical
            .service(live)
            .service(ready)
    })
    .bind("127.0.0.1:8080")?
    .run()
    .await
}
