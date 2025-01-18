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
- `sea-orm-cli migrate up`: run migration
- `sea-orm-cli generate entity -o entity/src`: generate entities

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

### TODO
1. Logout
    a. Créer un endpoint logout
    b. Le endpoint doit pouvoir générer une URL à partir d'un id token et appeler le logout endpoint et renvoyer un 201:redirec::logout url (à ajouter dans les param ou renvoyer vers redirect_uri)
    c. Supprimer le user de la session

   // Call the function
   match logout(&client, &user_id.clone().id_token.unwrap(), post_logout_redirect_uri).await {
       Ok(logout_url) => {
       info!("Logout URL: {}", logout_url);
       }
       Err(e) => {
       error!("Error generating logout URL: {}", e);
       }
    }

2. Avoir des endpoints de compte (create, list, get one)
3. Github Action : build test
4. Front
5. Docker + Descripteurs kube et deploy sur minikube
6. EKS + Terraform ou ECS + Terraform
7. Helm
8. Logs avec cloudwatch
9. Monitoring avec opentelemetry genre prometheus