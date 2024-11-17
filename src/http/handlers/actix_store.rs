use actix_web::{get, web, HttpResponse, Responder};
use openid::{Bearer, Client};
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
