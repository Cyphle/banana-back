use crate::domain::profile::{CreateProfileCommand, Profile};
use crate::dto::views::profile::ProfileView;
use crate::security::token::get_username_from_session;
use entity::profiles::Entity as ProfileEntity;
use openid::DiscoveredClient;
use sea_orm::ActiveValue::Set;
use sea_orm::{ActiveModelTrait, ColumnTrait, DatabaseConnection, DbErr, EntityTrait, QueryFilter};

pub async fn create(db_connexion: &DatabaseConnection, command: &CreateProfileCommand) -> Result<Profile, DbErr> {
    let model = entity::profiles::ActiveModel {
        username: Set(command.username.to_owned()),
        email: Set(command.email.to_owned()),
        first_name: Set(command.first_name.to_owned()),
        last_name: Set(command.last_name.to_owned()),
        ..Default::default()
    };

    model.clone().insert(db_connexion).await.map(|m| Profile::new(
        m.id,
        m.username,
        m.email,
        m.first_name,
        m.last_name
    ))
}

pub async fn find_one_by_id(db_connexion: &DatabaseConnection, id: i32) -> Result<Option<ProfileView>, DbErr> {
    ProfileEntity::find_by_id(id)
        .one(db_connexion)
        .await
        .map(|m| m.map(|m| ProfileView {
            id: m.id,
            username: m.username,
            email: m.email,
            first_name: m.first_name,
            last_name: m.last_name,
        }))
}

pub async fn find_one_by_username<F, R>(
    db_connexion: &DatabaseConnection,
    username: &str,
    mapper: F
) -> Result<Option<R>, DbErr>
where F: Fn(entity::profiles::Model) -> R
{
    ProfileEntity::find()
        .filter(entity::profiles::Column::Username.eq(format!("{}", username).as_str()))
        .one(db_connexion)
        .await
        .map(|m| m.map(|m| mapper(m)))
}

pub async fn find_one_from_session<F, R>(
    client: &DiscoveredClient,
    db_connexion: &DatabaseConnection,
    session: &actix_session::Session,
    mapper: F
) -> Result<Option<R>, DbErr>
where F: Fn(entity::profiles::Model) -> R
{
    let username = get_username_from_session(client, session).await;
    match username {
        Some(username) => find_one_by_username(db_connexion, &username, mapper).await,
        None => Ok(None)
    }
}

pub mod mappers {
    use crate::domain::profile::Profile;
    use crate::dto::views::profile::ProfileView;

    pub fn to_profile_view(model: entity::profiles::Model) -> ProfileView {
        ProfileView {
            id: model.id,
            username: model.username,
            email: model.email,
            first_name: model.first_name,
            last_name: model.last_name,
        }
    }

    pub fn to_profile(model: entity::profiles::Model) -> Profile {
        Profile::new(
            model.id,
            model.username,
            model.email,
            model.first_name,
            model.last_name
        )
    }
}

#[cfg(test)]
mod tests {
    mod read {
        use crate::domain::profile::Profile;
        use crate::dto::views::profile::ProfileView;
        use crate::repositories::profile::mappers::{to_profile, to_profile_view};
        use crate::repositories::profile::{find_one_by_id, find_one_by_username};
        use chrono::{FixedOffset, NaiveDate, NaiveDateTime, NaiveTime};
        use sea_orm::prelude::DateTimeWithTimeZone;
        use sea_orm::EntityTrait;
        use sea_orm::{DatabaseBackend, DbErr, MockDatabase, Transaction};

        #[async_std::test]
        async fn should_find_one_by_id() -> Result<(), DbErr> {
            let d = NaiveDate::from_ymd_opt(2015, 6, 3).unwrap();
            let t = NaiveTime::from_hms_milli_opt(12, 34, 56, 789).unwrap();

            let db = MockDatabase::new(DatabaseBackend::Postgres)
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
                .into_connection();

            let found = find_one_by_id(&db, 1).await?;

            assert_eq!(
                found,
                Some(ProfileView {
                    id: 1,
                    username: "johndoe".to_owned(),
                    email: "johndoe@banana.fr".to_owned(),
                    first_name: "John".to_owned(),
                    last_name: "Doe".to_owned(),
                })
            );

            // Checking transaction log
            assert_eq!(
                db.into_transaction_log(),
                [
                    Transaction::from_sql_and_values(
                        DatabaseBackend::Postgres,
                        r#"SELECT "profiles"."id", "profiles"."username", "profiles"."email", "profiles"."first_name", "profiles"."last_name", "profiles"."created_at", "profiles"."updated_at", "profiles"."deleted_at" FROM "profiles" WHERE "profiles"."id" = $1 LIMIT $2"#,
                        [1.into(), 1u64.into()],
                    ),
                ]
            );

            Ok(())
        }

        #[async_std::test]
        async fn should_find_one_by_username_with_profile_mapper() -> Result<(), DbErr> {
            let d = NaiveDate::from_ymd_opt(2015, 6, 3).unwrap();
            let t = NaiveTime::from_hms_milli_opt(12, 34, 56, 789).unwrap();

            let db = MockDatabase::new(DatabaseBackend::Postgres)
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
                .into_connection();

            let found = find_one_by_username(&db, "johndoe", to_profile).await?;


            assert_eq!(
                found,
                Some(Profile {
                    id: 1,
                    username: "johndoe".to_owned(),
                    email: "johndoe@banana.fr".to_owned(),
                    first_name: "John".to_owned(),
                    last_name: "Doe".to_owned(),
                })
            );

            // Checking transaction log
            assert_eq!(
                db.into_transaction_log(),
                [
                    Transaction::from_sql_and_values(
                        DatabaseBackend::Postgres,
                        r#"SELECT "profiles"."id", "profiles"."username", "profiles"."email", "profiles"."first_name", "profiles"."last_name", "profiles"."created_at", "profiles"."updated_at", "profiles"."deleted_at" FROM "profiles" WHERE "profiles"."username" = $1 LIMIT $2"#,
                        ["johndoe".into(), 1u64.into()],
                    ),
                ]
            );

            Ok(())
        }

        #[async_std::test]
        async fn should_find_one_by_username_with_profile_view_mapper() -> Result<(), DbErr> {
            let d = NaiveDate::from_ymd_opt(2015, 6, 3).unwrap();
            let t = NaiveTime::from_hms_milli_opt(12, 34, 56, 789).unwrap();

            let db = MockDatabase::new(DatabaseBackend::Postgres)
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
                .into_connection();

            let found = find_one_by_username(&db, "johndoe", to_profile_view).await?;


            assert_eq!(
                found,
                Some(ProfileView {
                    id: 1,
                    username: "johndoe".to_owned(),
                    email: "johndoe@banana.fr".to_owned(),
                    first_name: "John".to_owned(),
                    last_name: "Doe".to_owned(),
                })
            );

            // Checking transaction log
            assert_eq!(
                db.into_transaction_log(),
                [
                    Transaction::from_sql_and_values(
                        DatabaseBackend::Postgres,
                        r#"SELECT "profiles"."id", "profiles"."username", "profiles"."email", "profiles"."first_name", "profiles"."last_name", "profiles"."created_at", "profiles"."updated_at", "profiles"."deleted_at" FROM "profiles" WHERE "profiles"."username" = $1 LIMIT $2"#,
                        ["johndoe".into(), 1u64.into()],
                    ),
                ]
            );

            Ok(())
        }
    }

    mod create {
        use crate::domain::profile::{CreateProfileCommand, Profile};
        use crate::repositories::profile::create;
        use chrono::{FixedOffset, NaiveDate, NaiveDateTime, NaiveTime};
        use sea_orm::{
            entity::prelude::*, entity::*,
            DatabaseBackend, MockDatabase, MockExecResult, Transaction,
        };

        #[async_std::test]
        async fn should_create_todo_list() -> Result<(), DbErr> {
            let d = NaiveDate::from_ymd_opt(2015, 6, 3).unwrap();
            let t = NaiveTime::from_hms_milli_opt(12, 34, 56, 789).unwrap();

            let db = MockDatabase::new(DatabaseBackend::Postgres)
                .append_query_results([
                    [entity::profiles::Model {
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
                .append_exec_results([
                    MockExecResult {
                        last_insert_id: 15,
                        rows_affected: 1,
                    },
                ])
                .into_connection();

            // Prepare the ActiveModel
            let created = create(&db, &CreateProfileCommand {
                username: "johndoe".to_owned(),
                email: "johndoe@banana.fr".to_owned(),
                first_name: "John".to_owned(),
                last_name: "Doe".to_owned(),
            }).await;

            assert_eq!(
                created.unwrap(),
                Profile::new(
                    1,
                    "johndoe".to_owned(),
                    "johndoe@banana.fr".to_owned(),
                    "John".to_owned(),
                    "Doe".to_owned(),
                )
            );

            assert_eq!(
                db.into_transaction_log(),
                [
                    Transaction::from_sql_and_values(
                        DatabaseBackend::Postgres,
                        r#"INSERT INTO "profiles" ("username", "email", "first_name", "last_name") VALUES ($1, $2, $3, $4) RETURNING "id", "username", "email", "first_name", "last_name", "created_at", "updated_at", "deleted_at""#,
                        ["johndoe".into(), "johndoe@banana.fr".into(), "John".into(), "Doe".into()]
                    ),
                ]
            );

            Ok(())
        }
    }
}