use crate::dto::views::profile::ProfileView;
use crate::repositories::profile::find_one_by_username;
use sea_orm::DatabaseConnection;

pub async fn get_profile_by_username(db_connexion: &DatabaseConnection, username: &str) -> Option<ProfileView> {
    find_one_by_username(db_connexion, username).await.unwrap_or_else(|_| None)
}

#[cfg(test)]
mod tests {
    use crate::http::adapters::profile::get_profile_by_username;
    use chrono::{FixedOffset, NaiveDate, NaiveDateTime, NaiveTime};
    use sea_orm::prelude::DateTimeWithTimeZone;
    use sea_orm::{DatabaseBackend, DatabaseConnection, MockDatabase};

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

    #[async_std::test]
    async fn should_get_one_profile_by_username() {
        let db = get_mock_database();

        let profile = get_profile_by_username(db, "johndoe").await;

        assert_eq!(profile.unwrap().username, "johndoe");
    }
}