use crate::config::actix::AppState;
use crate::http::adapters::profile::get_profile_by_username;
use crate::repositories::profile::find_one_by_username;
use actix_web::{get, web, HttpResponse, Responder};
use log::{error, info};

#[get("/test")]
async fn test(state: web::Data<AppState>) -> impl Responder {
    match get_profile_by_username(&state.db_connection, "test").await {
        Some(username) => {
            info!("Found profile for username: {}", username.username);
            HttpResponse::Ok().json(username)
        }
        None => {
            error!("Nothing found");
            HttpResponse::Ok().finish()
        }
    }
}

#[cfg(test)]
mod tests {
    use super::test as test_endpoint;
    use actix_web::{test, App};
    use actix_web::http::StatusCode;
    use actix_web::web;
    use sea_orm::{DatabaseBackend, DatabaseConnection, MockDatabase};
    use chrono::{FixedOffset, NaiveDate, NaiveDateTime, NaiveTime};
    use sea_orm::prelude::DateTimeWithTimeZone;
    use std::collections::HashMap;
    use std::sync::Mutex;
    use openid::Bearer;
    use crate::config::actix::AppState;
    use crate::security::oidc::{OidcAdminConfig, OidcConfig};

    fn dummy_oidc_config() -> OidcConfig {
        OidcConfig::new(
            "http://localhost".to_string(),
            "client_id".to_string(),
            "client_secret".to_string(),
            "http://localhost/callback".to_string(),
            None,
            None,
            "http://localhost/logout".to_string(),
            OidcAdminConfig { 
                client_id: "admin_client".to_string(),
                client_secret: "admin_secret".to_string(),
                create_user_url: "http://localhost/create".to_string(),
            },
        )
    }

    fn mock_db_with_one_profile() -> &'static DatabaseConnection {
        let d = NaiveDate::from_ymd_opt(2015, 6, 3).unwrap();
        let t = NaiveTime::from_hms_milli_opt(12, 34, 56, 789).unwrap();
        Box::leak(Box::new(
            MockDatabase::new(DatabaseBackend::Postgres)
                .append_query_results([
                    vec![entity::profiles::Model {
                        id: 1,
                        username: "test".to_owned(),
                        email: "test@banana.fr".to_owned(),
                        first_name: "Test".to_owned(),
                        last_name: "User".to_owned(),
                        created_at: DateTimeWithTimeZone::from_naive_utc_and_offset(NaiveDateTime::new(d, t), FixedOffset::east_opt(0).unwrap()),
                        updated_at: DateTimeWithTimeZone::from_naive_utc_and_offset(NaiveDateTime::new(d, t), FixedOffset::east_opt(0).unwrap()),
                        deleted_at: None,
                    }],
                ])
                .into_connection(),
        ))
    }

    fn mock_db_without_profile() -> &'static DatabaseConnection {
        Box::leak(Box::new(
            MockDatabase::new(DatabaseBackend::Postgres)
                .append_query_results([Vec::<entity::profiles::Model>::new()])
                .into_connection(),
        ))
    }

    fn make_state(db: &'static DatabaseConnection) -> web::Data<AppState> {
        web::Data::new(AppState {
            db_connection: db,
            oidc_client: None,
            store: Mutex::new(HashMap::<String, Bearer>::new()),
            oidc_config: Some(dummy_oidc_config()),
        })
    }

    #[actix_web::test]
    async fn get_test_should_return_profile_when_found() {
        let state = make_state(mock_db_with_one_profile());
        let app = test::init_service(App::new().app_data(state.clone()).service(test_endpoint)).await;
        let req = test::TestRequest::get().uri("/test").to_request();
        let resp = test::call_service(&app, req).await;
        assert_eq!(resp.status(), StatusCode::OK);
        let body = test::read_body(resp).await;
        let body_str = String::from_utf8(body.to_vec()).unwrap();
        assert!(body_str.contains("\"username\":\"test\""));
    }

    #[actix_web::test]
    async fn get_test_should_return_ok_with_empty_body_when_not_found() {
        let state = make_state(mock_db_without_profile());
        let app = test::init_service(App::new().app_data(state.clone()).service(test_endpoint)).await;
        let req = test::TestRequest::get().uri("/test").to_request();
        let resp = test::call_service(&app, req).await;
        assert_eq!(resp.status(), StatusCode::OK);
        let body = test::read_body(resp).await;
        assert!(body.is_empty());
    }
}