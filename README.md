# Banana - Efficient bank account manager

## Tests
Launch all tests: 
`go test -v ./...`

Launch tests of one subdirectory
`go test repositoties/*`

## TODO
- Bun & test containers > OK
- Account endspoints with Echo > One OK
- Bun transaction (cf // TODO c'est pas bien cette gestion de la transaction vu qu'elle va jamais s'arrêter là. Cf https://bun.uptrace.dev/guide/transactions.html#runintx RunInTx)
- Check token against IDP genre Keycloak (cf middleware et group middlewares)
- OpenID Connect flow avec React