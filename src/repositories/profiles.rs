use crate::domain;
use crate::domain::profiles::{CreateProfileCommand, Profile};
use sea_orm::ActiveValue::Set;
use sea_orm::{ActiveModelTrait, DatabaseConnection, DbErr};

pub async fn create(db_connexion: &DatabaseConnection, command: &CreateProfileCommand) -> Result<Profile, DbErr> {
    let model = entity::profiles::ActiveModel {
        username: Set(command.username.to_owned()),
        email: Set(command.email.to_owned()),
        first_name: Set(command.first_name.to_owned()),
        last_name: Set(command.last_name.to_owned()),
        ..Default::default()
    };

    model.clone().insert(db_connexion).await.map(|m| domain::profiles::new_profile(
        m.id,
        m.username,
        m.email,
        m.first_name,
        m.last_name
    ))
}

// pub async fn find_one_by_id(db_connexion: &DatabaseConnection, id: i32) -> Result<Option<TodoListView>, DbErr> {
//     TodoLists::find_by_id(id)
//         .one(db_connexion)
//         .await
//         .map(|m| m.map(|m| TodoListView {
//             id: m.id,
//             title: m.title,
//         }))
// }
//
// pub async fn find_all(db_connexion: &DatabaseConnection) -> Result<Vec<TodoListView>, DbErr> {
//     TodoLists::find()
//         .all(db_connexion)
//         .await
//         .map(|models| {
//             models
//                 .into_iter()
//                 .map(|m| TodoListView {
//                     id: m.id,
//                     title: m.title,
//                 })
//                 .collect()
//         })
// }
//
// #[cfg(test)]
// mod tests {
//     mod read {
//         use sea_orm::EntityTrait;
//         use sea_orm::{DatabaseBackend, DbErr, MockDatabase, Transaction};
//
//         use crate::dto::views::todo_list::TodoListView;
//         use crate::repositories::todo_lists::{find_all, find_one_by_id};
//
//         #[async_std::test]
//         async fn should_find_one_by_id() -> Result<(), DbErr> {
//             let db = MockDatabase::new(DatabaseBackend::Postgres)
//                 .append_query_results([
//                     vec![entity::todo_lists::Model {
//                         id: 1,
//                         title: "New York Cheese".to_owned(),
//                     }],
//                 ])
//                 .into_connection();
//
//             let found = find_one_by_id(&db, 1).await?;
//
//             assert_eq!(
//                 found,
//                 Some(TodoListView {
//                     id: 1,
//                     title: "New York Cheese".to_owned(),
//                 })
//             );
//
//             // Checking transaction log
//             assert_eq!(
//                 db.into_transaction_log(),
//                 [
//                     Transaction::from_sql_and_values(
//                         DatabaseBackend::Postgres,
//                         r#"SELECT "todo_lists"."id", "todo_lists"."title" FROM "todo_lists" WHERE "todo_lists"."id" = $1 LIMIT $2"#,
//                         [1.into(), 1u64.into()],
//                     ),
//                 ]
//             );
//
//             Ok(())
//         }
//
//         #[async_std::test]
//         async fn should_find_all() -> Result<(), DbErr> {
//             let db = MockDatabase::new(DatabaseBackend::Postgres)
//                 .append_query_results([
//                     vec![
//                         entity::todo_lists::Model {
//                             id: 1,
//                             title: "New York Cheese".to_owned(),
//                         },
//                         entity::todo_lists::Model {
//                             id: 2,
//                             title: "Apple Pie".to_owned(),
//                         },
//                     ],
//                 ])
//                 .into_connection();
//
//             let found = find_all(&db).await?;
//
//             assert_eq!(
//                 found,
//                 vec![
//                     TodoListView {
//                         id: 1,
//                         title: "New York Cheese".to_owned(),
//                     },
//                     TodoListView {
//                         id: 2,
//                         title: "Apple Pie".to_owned(),
//                     },
//                 ]
//             );
//
//             // Checking transaction log
//             assert_eq!(
//                 db.into_transaction_log(),
//                 [
//                     Transaction::from_sql_and_values(
//                         DatabaseBackend::Postgres,
//                         r#"SELECT "todo_lists"."id", "todo_lists"."title" FROM "todo_lists""#,
//                         [],
//                     ),
//                 ]
//             );
//
//             Ok(())
//         }
//     }
//
//     mod create {
//         use crate::domain::todo_list::{CreateTodoListCommand, TodoList};
//         use crate::repositories::todo_lists::create;
//         use sea_orm::{
//             entity::prelude::*, entity::*,
//             DatabaseBackend, MockDatabase, MockExecResult, Transaction,
//         };
//
//         #[async_std::test]
//         async fn should_create_todo_list() -> Result<(), DbErr> {
//             let db = MockDatabase::new(DatabaseBackend::Postgres)
//                 .append_query_results([
//                     [entity::todo_lists::Model {
//                         id: 15,
//                         title: "Apple Pie".to_owned(),
//                     }],
//                 ])
//                 .append_exec_results([
//                     MockExecResult {
//                         last_insert_id: 15,
//                         rows_affected: 1,
//                     },
//                 ])
//                 .into_connection();
//
//             // Prepare the ActiveModel
//             let created = create(&db, CreateTodoListCommand {
//                 title: "Apple Pie".to_string()
//             }).await;
//
//             assert_eq!(
//                 created.unwrap(),
//                 TodoList {
//                     id: 15,
//                     title: "Apple Pie".to_owned()
//                 }
//             );
//
//             assert_eq!(
//                 db.into_transaction_log(),
//                 [
//                     Transaction::from_sql_and_values(
//                         DatabaseBackend::Postgres,
//                         r#"INSERT INTO "todo_lists" ("title") VALUES ($1) RETURNING "id", "title""#,
//                         ["Apple Pie".into()]
//                     ),
//                 ]
//             );
//
//             Ok(())
//         }
//     }
// }