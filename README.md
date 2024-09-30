# Banana - Efficient bank account manager

## Tests
Launch all tests: 
`go test -v ./...`

Launch tests of one subdirectory
`go test repositoties/*`

## Commands
- To add migration file run in CLI:
`migrate create -ext sql -dir migrations -seq <migration name>`

## Docs
- https://github.com/golang-migrate/migrate/blob/master/GETTING_STARTED.md

## TODO
- Bun & test containers > OK
- Account endspoints with Echo > OK
- Bun transaction 
  - OK dans les repo.
  - Peut mieux faire et englober tout la request. Genre démarrer une transaction dans l'adapteur et la faire transiter puis l'arrêter
  - Par exemple ne plus faire que les repos ont déjà un client database mais requiert pour chaque fonction une transaction (moins pratique par contre)
- Check token against IDP genre Keycloak (cf middleware et group middlewares)
- OpenID Connect flow avec React et le back récupère le token (valide le jwt)
  - https://github.com/zitadel/oidc?tab=readme-ov-file => tout est pas clair. Genre c'est quoi le Keyfile
  - https://github.com/alexedwards/scs?tab=readme-ov-file#configuring-session-behavior
  - https://www.alexedwards.net/blog/working-with-cookies-in-go
  - Il faut lire les specs de OAuth2 et OIDC (& PKCE & JWT Profile)
- Makefile or taskfile
- Swagger

## Features
- Create a profile
- Create a bank account
- Update a bank account
- Add expenses
- Add credits
- Add budgets
- Add charges
- Get account for a given month (no indicators)
- Calculations
  - Amount at start of month => to be calculated for next month of today in advance