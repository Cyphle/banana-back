import { MigrationInterface, QueryRunner } from "typeorm";

export class PostRefactoring1753979783733 implements MigrationInterface {

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            CREATE TABLE IF NOT EXISTS users(
                id SERIAL PRIMARY KEY,
                firstname VARCHAR(255),
                lastname VARCHAR(255),
                age INTEGER
            )
        `);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`DROP TABLE IF EXISTS users;`);
    }

}
