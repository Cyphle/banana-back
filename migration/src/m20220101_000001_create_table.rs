use sea_orm_migration::{prelude::*, schema::*};

#[derive(DeriveMigrationName)]
pub struct Migration;

#[async_trait::async_trait]
impl MigrationTrait for Migration {
    async fn up(&self, manager: &SchemaManager) -> Result<(), DbErr> {
        let db = manager.get_connection();

        // Use `execute_unprepared` if the SQL statement doesn't have value bindings
        db.execute_unprepared(
            "CREATE TABLE IF NOT EXISTS profiles(
                id          SERIAL,
                username    VARCHAR(255) NOT NULL,
                email       VARCHAR(255) NOT NULL,
                first_name  VARCHAR(255) NOT NULL,
                last_name   VARCHAR(255) NOT NULL,
                created_at  timestamptz NOT NULL DEFAULT now(),
                updated_at  timestamptz NOT NULL DEFAULT now(),
                deleted_at  timestamptz,
                PRIMARY KEY (id)
            )"
        )
            .await?;

        Ok(())
    }

    async fn down(&self, manager: &SchemaManager) -> Result<(), DbErr> {
        manager
            .get_connection()
            .execute_unprepared("DROP TABLE IF EXISTS profiles;")
            .await?;

        Ok(())
    }
}
