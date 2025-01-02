use crate::config::actix::AppState;
use crate::domain::profile::CreateProfileCommand;
use crate::dto::requests::profile::CreateProfileRequest;
use crate::repositories;
use actix_web::{get, post, web, HttpResponse, Responder};

#[get("/profiles/{id}")]
async fn get_profile_by_id(path: web::Path<i32>, state: web::Data<AppState>) -> impl Responder {
    // TODO utiliser la session et du coup les tokens et adapter les tests. peut être extraire pour tester facilement
    match repositories::profile::find_one_by_id(&state.db_connection, path.into_inner()).await {
        Ok(Some(profile)) => HttpResponse::Ok().json(profile),
        Ok(None) => HttpResponse::NotFound().body("No profile found"),
        Err(_) => HttpResponse::InternalServerError().finish(),
    }
}

#[post("/profiles")]
async fn create_profile(payload: web::Json<CreateProfileRequest>, state: web::Data<AppState>) -> impl Responder {
    let moved_payload = payload.into_inner();
    match repositories::profile::create(
        &state.db_connection,
        &CreateProfileCommand {
            username: moved_payload.username.to_owned(),
            email: moved_payload.email.to_owned(),
            first_name: moved_payload.first_name.to_owned(),
            last_name: moved_payload.last_name.to_owned(),
        },
    ).await {
        Ok(_) => HttpResponse::Created().finish(),
        Err(_) => HttpResponse::InternalServerError().finish(),
    }
}

#[cfg(test)]
mod tests {
    use crate::config::actix::AppState;
    use crate::config::local::oidc_config::get_oidc_config;
    use crate::http::handlers::profile::{create_profile, get_profile_by_id};
    use actix_web::http::header::ContentType;
    use actix_web::web;
    use actix_web::{test, App};
    use chrono::{FixedOffset, NaiveDate, NaiveDateTime, NaiveTime};
    use sea_orm::prelude::DateTimeWithTimeZone;
    use sea_orm::{DatabaseBackend, DatabaseConnection, MockDatabase};
    use std::sync::Mutex;

    fn get_mock_database() -> &'static DatabaseConnection {
        let d = NaiveDate::from_ymd_opt(2015, 6, 3).unwrap();
        let t = NaiveTime::from_hms_milli_opt(12, 34, 56, 789).unwrap();

        Box::leak(Box::new(MockDatabase::new(DatabaseBackend::Postgres)
            .append_query_results([
                vec![entity::profiles::Model {
                    id: 1,
                    username: "johndoe".to_owned(),
                    email: "johndoe@banana.fr".to_owned(),
                    first_name: "John".to_owned(),
                    last_name: "Doe".to_owned(),
                    created_at: DateTimeWithTimeZone::from_naive_utc_and_offset(NaiveDateTime::new(d, t), FixedOffset::east_opt(0).unwrap()),
                    updated_at: DateTimeWithTimeZone::from_naive_utc_and_offset(NaiveDateTime::new(d, t), FixedOffset::east_opt(0).unwrap()),
                    deleted_at: None,
                }],
            ])
            .into_connection()))
    }

    #[actix_web::test]
    async fn should_get_one_profile_by_id() {
        let app = test::init_service(
            App::new()
                .app_data(web::Data::new(AppState {
                    db_connection: get_mock_database(),
                    oidc_client: None,
                    store: Mutex::new(std::collections::HashMap::new()),
                    oidc_config: get_oidc_config().clone(),
                }))
                .service(get_profile_by_id)
        ).await;
        let req = test::TestRequest::get().uri("/profiles/123")
            .insert_header(ContentType::plaintext())
            .to_request();

        let resp = test::call_service(&app, req).await;

        assert!(resp.status().is_success());
    }

    #[actix_web::test]
    async fn should_create_profile() {
        let app = test::init_service(
            App::new()
                .app_data(web::Data::new(AppState {
                    db_connection: get_mock_database(),
                    oidc_client: None,
                    store: Mutex::new(std::collections::HashMap::new()),
                    oidc_config: get_oidc_config().clone(),
                }))
                .service(create_profile)
        ).await;

        let req = test::TestRequest::post()
            .set_payload("{\"username\": \"johndoe\", \"email\": \"johndoe@banana.fr\", \"first_name\": \"John\", \"last_name\": \"Doe\"}")
            .insert_header(ContentType::json())
            .uri("/profiles").to_request();

        let resp = test::call_service(&app, req).await;
        assert!(resp.status().is_success());
    }
}