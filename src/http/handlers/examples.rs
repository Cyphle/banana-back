use actix_web::{get, post, web, HttpResponse, Responder};

#[get("/hello")]
pub async fn hello() -> impl Responder {
    HttpResponse::Ok().body("Hello world!")
}

#[derive(serde::Deserialize)]
struct Test {
    key: String
}

#[post("/echo")]
pub async fn echo(_: web::Json<Test>) -> impl Responder {
    HttpResponse::Ok().body("Bonjour")
}

#[cfg(test)]
mod tests {
    use sea_orm::{
        entity::prelude::*, entity::*
    };

    mod actix_tests {
        use actix_web::http::header::ContentType;
        use actix_web::{test, App};
        use super::super::{hello, echo};

        #[actix_web::test]
        async fn test_hello_get() {
            let app = test::init_service(App::new().service(hello)).await;
            let req = test::TestRequest::get().uri("/hello")
                .insert_header(ContentType::plaintext())
                .to_request();
            let resp = test::call_service(&app, req).await;
            assert!(resp.status().is_success());
        }

        #[actix_web::test]
        async fn test_echo_post() {
            let app = test::init_service(App::new().service(echo)).await;
            let req = test::TestRequest::post()
                .set_payload("{ \"key\": \"Bonjour\" }")
                .insert_header(ContentType::json())
                .uri("/echo")
                .to_request();
            let resp = test::call_service(&app, req).await;
            assert!(resp.status().is_success());
        }

        #[actix_web::test]
        async fn test_echo_post_error() {
            let app = test::init_service(App::new().service(echo)).await;
            let req = test::TestRequest::post().uri("/").to_request();
            let resp = test::call_service(&app, req).await;
            assert!(resp.status().is_client_error());
        }
    }
}