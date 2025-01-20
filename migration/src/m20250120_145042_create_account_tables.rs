use sea_orm_migration::{prelude::*, schema::*};

#[derive(DeriveMigrationName)]
pub struct Migration;

#[async_trait::async_trait]
impl MigrationTrait for Migration {
    async fn up(&self, manager: &SchemaManager) -> Result<(), DbErr> {
        let db = manager.get_connection();

        // type: PERSONAL, PERSONAL_SHARED, PERSONAL_SAVINGS, SHARED_SAVINGS
        db.execute_unprepared(
            "CREATE TABLE IF NOT EXISTS accounts(
                id              SERIAL,
                name            VARCHAR(255) NOT NULL,
                type            VARCHAR(255) NOT NULL,
                profile_id      INTEGER NOT NULL,
                starting_amount NUMERIC(20, 2),
                created_at      timestamptz NOT NULL DEFAULT now(),
                updated_at      timestamptz NOT NULL DEFAULT now(),
                deleted_at      timestamptz,
                PRIMARY KEY (id),
                CONSTRAINT fk_accounts_profiles FOREIGN KEY (profile_id) REFERENCES profiles (id)
                ON DELETE CASCADE
                ON UPDATE CASCADE
            )"
        )
            .await?;

        db.execute_unprepared(
            ""
        )
            .await?;

        Ok(())
    }

    async fn down(&self, manager: &SchemaManager) -> Result<(), DbErr> {
        manager
            .get_connection()
            .execute_unprepared("DROP TABLE IF EXISTS accounts;")
            .await?;

        Ok(())
    }
}

