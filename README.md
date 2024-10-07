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

#### Commands
- `sea-orm-cli migrate generate <migration name>`: create a new migration
- `sea-orm-cli migrate up`: run migration
- `sea-orm-cli generate entity -o entity/src`: generate entities

### Actix

### Logging
- See `config::logger`
- Use macros `log::info!`, `log::error!`, `log::warn!`, `log::debug!`, `log::trace!`