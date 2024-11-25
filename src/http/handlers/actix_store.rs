use actix_web::{get, web, HttpResponse, Responder};
use actix_web::web::Data;
use openid::{Bearer, Client, Token};
use serde::Serialize;
use crate::config::actix::AppState;

#[derive(Serialize)]
struct GetResponse {
    key: String,
    value: Option<Bearer>,
}

#[get("/set-in-shared-state")]
async fn add_to_store(data: web::Data<AppState>) -> impl Responder {
    let mut store = data.store.lock().unwrap();
    // Example of storing in shared state
    // save_in_shared_state(&data, token);
    HttpResponse::Ok().json("Key set successfully")
}

#[get("/get-from-shared-state")]
async fn get_from_store(data: web::Data<AppState>) -> impl Responder {
    let store = data.store.lock().unwrap();
    let value = store.get(&"hello".to_string()).cloned();
    HttpResponse::Ok().json(GetResponse {
        key: "hello".to_string(),
        value: value,
    })
}

fn save_in_shared_state(state: &Data<AppState>, token: Token) {
    let mut store = state.store.lock().unwrap();
    store.insert("hello".to_string(), token.bearer.clone());
}
