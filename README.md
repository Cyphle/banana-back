# Banana

## Technical instructions

### Sea ORM

#### Installation
1. In Cargo.toml, add the following dependencies:
``` toml
sea-orm = { version = "1.0.0-rc.5", features = [ <DATABASE_DRIVER>, <ASYNC_RUNTIME>, "macros" ] }
```
2. See `config::database` for database configuration
3. Add migration crate `sea-orm-cli migrate init`
4. Configure migration crate `.toml` file and add the crate in the main crate
5. Add entity crate `cargo new entity`
6. Run the command to generate the entities and add the crate in dependencies of the main one

#### Commands
- `sea-orm-cli migrate generate <migration name>`: create a new migration
- `sea-orm-cli migrate up --databade-url "postgres://postgres:postgres@localhost:5434/banana"`: run migration
- `sea-orm-cli generate entity -o entity/src`: generate entities
- `cargo build --release`: build the project for production

### Actix
1. Add the dependency
2. Define the server and launch the main in async mode
```rust
pub async fn config(db_connection: &'static DatabaseConnection) -> std::io::Result<()> {
    info!("Starting Actix server");

    HttpServer::new(|| {
        App::new()
            .app_data(web::Data::new(HandlerState {
                db_connection: db_connection
            }))
            .service(get_profile_by_id)
            .service(create_profile)
    })
        .bind(("127.0.0.1", 8080))?
        .run()
        .await
}

config::actix::config(static_db).await;
```

### Logging
- See `config::logger`
- Use macros `log::info!`, `log::error!`, `log::warn!`, `log::debug!`, `log::trace!`

### OIDC
- See `config::oidc`
- See crate `security`

### TODO
2. Avoir des endpoints de profile et de compte (create, list, get one, delete)
- Il faut aussi refiare un peu de front avec un parcours un peu mieux
3. Front
4. Github Action : build test
5. Docker + Descripteurs kube et deploy sur minikube
6. EKS + Terraform ou ECS + Terraform
7. Deploy sur AWS
8. Helm
9. Logs avec cloudwatch
10. Monitoring avec opentelemetry genre prometheus
11. Pour tester plus facilement, il faudrait passer en objet et injecter les dépendances genre initialiser les repositories et les mettre dans le state actix et apres un peu mock avec trait tout ça 