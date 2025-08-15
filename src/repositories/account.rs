use crate::domain::account::{Account, CreateAccountCommand};
use entity::accounts::Entity as AccountEntity;
use sea_orm::prelude::Decimal;
use sea_orm::{ActiveModelTrait, ColumnTrait, DatabaseConnection, DbErr, EntityTrait, QueryFilter, Set};
use rust_decimal::Decimal as rDecimal;
use rust_decimal::prelude::{FromPrimitive, ToPrimitive};

pub async fn create(
    db_connexion: &DatabaseConnection,
    command: &CreateAccountCommand,
) -> Result<Account, DbErr> {
    let amount: Option<Decimal> = rDecimal::from_f64(command.starting_amount);
    let model = entity::accounts::ActiveModel {
        name: Set(command.name.to_owned()),
        r#type: Set(command.r#type.to_owned()),
        profile_id: Set(command.profile_id),
        starting_amount: Set(amount),
        ..Default::default()
    };

    model.clone().insert(db_connexion).await.map(|m| Account {
        id: m.id,
        name: m.name,
        r#type: m.r#type,
        starting_amount: m.starting_amount.unwrap().to_f64().unwrap(),
        creation_date: m.created_at.into(),
    })
}

pub async fn find_all(
    db_connexion: &DatabaseConnection,
    profile_id: i32,
) -> Result<Vec<Account>, DbErr> {
    AccountEntity::find()
        .filter(entity::accounts::Column::ProfileId.eq(profile_id))
        .all(db_connexion)
        .await
        .map(|models| {
            models
                .iter()
                .map(|m| Account {
                    id: m.id,
                    name: m.name.clone(),
                    r#type: m.r#type.clone(),
                    starting_amount: m.starting_amount.unwrap().to_f64().unwrap(),
                    creation_date: m.created_at.into(),
                })
                .collect()
        })
}

pub async fn find_by_id(db_connexion: &DatabaseConnection, account_id: i32, profile_id: i32) -> Result<Option<Account>, DbErr> {
    AccountEntity::find()
        .filter(entity::accounts::Column::Id.eq(account_id))
        .filter(entity::accounts::Column::ProfileId.eq(profile_id))
        .one(db_connexion)
        .await
        .map(|m| m.map(|m| Account {
            id: m.id,
            name: m.name,
            r#type: m.r#type,
            starting_amount: m.starting_amount.unwrap().to_f64().unwrap(),
            creation_date: m.created_at.into(),
        }))
}

#[cfg(test)]
mod tests {
    mod read {
        use crate::repositories::account::{find_all, find_by_id};
        use chrono::{FixedOffset, NaiveDate, NaiveDateTime, NaiveTime};
        use sea_orm::entity::prelude::*;
        use sea_orm::prelude::*;
        use sea_orm::{MockDatabase, Transaction};
        use sea_orm::DatabaseBackend;

        #[tokio::test]
        async fn test_find_by_id() {
            let d = NaiveDate::from_ymd_opt(2015, 6, 3).unwrap();
            let t = NaiveTime::from_hms_milli_opt(12, 34, 56, 789).unwrap();

            let db = MockDatabase::new(DatabaseBackend::Postgres)
                .append_query_results(vec![vec![entity::accounts::Model {
                    id: 1,
                    name: "Test Account".to_owned(),
                    r#type: "Savings".to_owned(),
                    profile_id: 1,
                    starting_amount: Some(Decimal::new(1000, 2)),
                    created_at: DateTimeWithTimeZone::from_naive_utc_and_offset(NaiveDateTime::new(d, t), FixedOffset::east_opt(0).unwrap()),
                    updated_at: DateTimeWithTimeZone::from_naive_utc_and_offset(NaiveDateTime::new(d, t), FixedOffset::east_opt(0).unwrap()),
                    deleted_at: None,
                }]])
                .into_connection();

            let result = find_by_id(&db, 1, 1).await.unwrap().unwrap();

            assert_eq!(result.id, 1);
            assert_eq!(result.name, "Test Account");
            assert_eq!(result.r#type, "Savings");
            assert_eq!(result.starting_amount, 10.0);

            assert_eq!(
                db.into_transaction_log(),
                [
                    Transaction::from_sql_and_values(
                        DatabaseBackend::Postgres,
                        r#"SELECT "accounts"."id", "accounts"."name", "accounts"."type", "accounts"."profile_id", "accounts"."starting_amount", "accounts"."created_at", "accounts"."updated_at", "accounts"."deleted_at" FROM "accounts" WHERE "accounts"."id" = $1 AND "accounts"."profile_id" = $2 LIMIT $3"#,
                        [1.into(), 1.into(), 1u64.into()],
                    ),
                ]
            );
        }

        #[tokio::test]
        async fn test_find_all() {
            let d = NaiveDate::from_ymd_opt(2015, 6, 3).unwrap();
            let t = NaiveTime::from_hms_milli_opt(12, 34, 56, 789).unwrap();

            let db = MockDatabase::new(DatabaseBackend::Postgres)
                .append_query_results(vec![vec![
                    entity::accounts::Model {
                        id: 1,
                        name: "Test Account 1".to_owned(),
                        r#type: "Savings".to_owned(),
                        profile_id: 1,
                        starting_amount: Some(Decimal::new(1000, 2)),
                        created_at: DateTimeWithTimeZone::from_naive_utc_and_offset(NaiveDateTime::new(d, t), FixedOffset::east_opt(0).unwrap()),
                        updated_at: DateTimeWithTimeZone::from_naive_utc_and_offset(NaiveDateTime::new(d, t), FixedOffset::east_opt(0).unwrap()),
                        deleted_at: None,
                    },
                    entity::accounts::Model {
                        id: 2,
                        name: "Test Account 2".to_owned(),
                        r#type: "Checking".to_owned(),
                        profile_id: 1,
                        starting_amount: Some(Decimal::new(2000, 2)),
                        created_at: DateTimeWithTimeZone::from_naive_utc_and_offset(NaiveDateTime::new(d, t), FixedOffset::east_opt(0).unwrap()),
                        updated_at: DateTimeWithTimeZone::from_naive_utc_and_offset(NaiveDateTime::new(d, t), FixedOffset::east_opt(0).unwrap()),
                        deleted_at: None,
                    },
                ]])
                .into_connection();

            let results = find_all(&db, 1).await.unwrap();

            assert_eq!(results.len(), 2);
            assert_eq!(results[0].id, 1);
            assert_eq!(results[0].name, "Test Account 1");
            assert_eq!(results[0].r#type, "Savings");
            assert_eq!(results[0].starting_amount, 10.0);
            assert_eq!(results[1].id, 2);
            assert_eq!(results[1].name, "Test Account 2");
            assert_eq!(results[1].r#type, "Checking");
            assert_eq!(results[1].starting_amount, 20.0);

            assert_eq!(
                db.into_transaction_log(),
                [
                    Transaction::from_sql_and_values(
                        DatabaseBackend::Postgres,
                        r#"SELECT "accounts"."id", "accounts"."name", "accounts"."type", "accounts"."profile_id", "accounts"."starting_amount", "accounts"."created_at", "accounts"."updated_at", "accounts"."deleted_at" FROM "accounts" WHERE "accounts"."profile_id" = $1"#,
                        [1.into()],
                    ),
                ]
            );
        }
    }

    mod create {
        use crate::repositories::account::create;
        use crate::domain::account::{CreateAccountCommand, Account};
        use sea_orm::{DatabaseBackend, MockDatabase, MockExecResult, Transaction};
        use sea_orm::prelude::*;
        use rust_decimal::Decimal as rDecimal;
        use chrono::{FixedOffset, NaiveDate, NaiveDateTime, NaiveTime};
        use rust_decimal::prelude::FromPrimitive;
        use crate::domain::profile::Profile;

        #[tokio::test]
        async fn should_create_account() -> Result<(), DbErr> {
            let d = NaiveDate::from_ymd_opt(2015, 6, 3).unwrap();
            let t = NaiveTime::from_hms_milli_opt(12, 34, 56, 789).unwrap();
            let profile = Profile::new(
                1,
                "johndoe".to_owned(),
                "johndoe@banana.fr".to_owned(),
                "John".to_owned(),
                "Doe".to_owned(),
            );

            let db = MockDatabase::new(DatabaseBackend::Postgres)
                .append_query_results([
                    [entity::accounts::Model {
                        id: 1,
                        name: "Test Account".to_owned(),
                        r#type: "Savings".to_owned(),
                        profile_id: 1,
                        starting_amount: Some(Decimal::new(1000, 2)),
                        created_at: DateTimeWithTimeZone::from_naive_utc_and_offset(NaiveDateTime::new(d, t), FixedOffset::east_opt(0).unwrap()),
                        updated_at: DateTimeWithTimeZone::from_naive_utc_and_offset(NaiveDateTime::new(d, t), FixedOffset::east_opt(0).unwrap()),
                        deleted_at: None,
                    }],
                ])
                .append_exec_results([
                    MockExecResult {
                        last_insert_id: 1,
                        rows_affected: 1,
                    },
                ])
                .into_connection();

            let command = CreateAccountCommand {
                name: "Test Account".to_owned(),
                r#type: "Savings".to_owned(),
                profile_id: profile.id,
                starting_amount: 10.0,
            };

            let created = create(&db, &command).await?;

            assert_eq!(
                created,
                Account {
                    id: 1,
                    name: "Test Account".to_owned(),
                    r#type: "Savings".to_owned(),
                    starting_amount: 10.0,
                    creation_date: DateTimeWithTimeZone::from_naive_utc_and_offset(NaiveDateTime::new(d, t), FixedOffset::east_opt(0).unwrap()).into(),
                }
            );

            assert_eq!(
                db.into_transaction_log(),
                [
                    Transaction::from_sql_and_values(
                        DatabaseBackend::Postgres,
                        r#"INSERT INTO "accounts" ("name", "type", "profile_id", "starting_amount") VALUES ($1, $2, $3, $4) RETURNING "id", "name", "type", "profile_id", "starting_amount", "created_at", "updated_at", "deleted_at""#,
                        ["Test Account".into(), "Savings".into(), 1.into(), rDecimal::from_f64(10.0).unwrap().into()]
                    ),
                ]
            );

            Ok(())
        }
    }
}
