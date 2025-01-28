use crate::config::actix::AppState;
use crate::domain::account::CreateAccountCommand;
use crate::dto::requests::account::CreateAccountRequest;
use crate::repositories::profile::mappers::to_profile;
use crate::repositories::profile::{find_one_from_session};
use actix_session::Session;
use actix_web::web::Data;
use actix_web::{get, post, web, HttpResponse, Responder};
use log::error;
use crate::dto::views::account::AccountView;
use crate::repositories;
use crate::repositories::account::create;

// TODO Extract la logique dans un adapter
#[post("/accounts")]
pub async fn create_account(
    payload: web::Json<CreateAccountRequest>,
    session: Session,
    state: Data<AppState>,
) -> impl Responder {
    let client = state.oidc_client.as_ref().unwrap().lock().unwrap();
    match find_one_from_session(&client, &state.db_connection, &session, to_profile).await {
        Ok(profile) => {
            match profile {
                None => {
                    HttpResponse::Forbidden().finish()
                }
                Some(profile) => {
                    let command = CreateAccountCommand {
                        name: payload.name.to_owned(),
                        r#type: payload.r#type.to_owned(),
                        starting_amount: payload.starting_amount.to_owned(),
                        profile_id: profile.id,
                    };

                    create(&state.db_connection, &command).await;

                    HttpResponse::Ok().finish()
                }
            }

            HttpResponse::Ok().finish()
        }
        Err(e) => {
            HttpResponse::Forbidden().finish()
        }
    }
}

#[get("/accounts/{id}")]
pub async fn find_one(
    path: web::Path<i32>,
    session: Session,
    state: Data<AppState>
) -> impl Responder {
    let client = state.oidc_client.as_ref().unwrap().lock().unwrap();

    match find_one_from_session(&client, &state.db_connection, &session, to_profile).await {
        Ok(profile) => {
            match profile {
                None => {
                    HttpResponse::Forbidden().finish()
                }
                Some(profile) => {
                    match repositories::account::find_by_id(&state.db_connection, path.into_inner(), profile.id).await {
                        Ok(Some(account)) => {
                            HttpResponse::Ok().json(AccountView::from(account))
                        }
                        Ok(None) => HttpResponse::NotFound().body("No account found"),
                        Err(e) => {
                            error!("Error finding account: {:?}", e);
                            HttpResponse::InternalServerError().finish()
                        },
                    }
                }
            }
        }
        Err(e) => {
            HttpResponse::Forbidden().finish()
        }
    }
}

// Find all
#[get("/accounts")]
pub async fn find_all(
    session: Session,
    state: Data<AppState>
) -> impl Responder {
    let client = state.oidc_client.as_ref().unwrap().lock().unwrap();

    match find_one_from_session(&client, &state.db_connection, &session, to_profile).await {
        Ok(profile) => {
            match profile {
                None => {
                    HttpResponse::Forbidden().finish()
                }
                Some(profile) => {
                    match repositories::account::find_all(&state.db_connection, profile.id).await {
                        Ok(accounts) => {
                            let accounts_view: Vec<AccountView> = accounts.into_iter().map(|account| AccountView::from(account)).collect();
                            HttpResponse::Ok().json(accounts_view)
                        }
                        Err(e) => {
                            error!("Error finding accounts: {:?}", e);
                            HttpResponse::InternalServerError().finish()
                        },
                    }
                }
            }
        }
        Err(e) => {
            HttpResponse::Forbidden().finish()
        }
    }
}