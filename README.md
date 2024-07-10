# Banana - Efficient bank account manager

## Tests
Launch all tests: 
`go test -v ./...`

Launch tests of one subdirectory
`go test repositoties/*`

## TODO
- Bun & test containers > OK
- Account endspoints with Echo > OK manque delete et update
- Bun transaction 
  - OK dans les repo. 
  - Peut mieux faire et englober tout la request. Genre démarrer une transaction dans l'adapteur et la faire transiter puis l'arrêter
  - Par exemple ne plus faire que les repos ont déjà un client database mais requiert pour chaque fonction une transaction (moins pratique par contre)
- Check token against IDP genre Keycloak (cf middleware et group middlewares)
- OpenID Connect flow avec React
- Makefile or taskfile