import { decorateWithUser } from './decorators/authentication.decorator';
import { initFastify } from './config/fastify.config';
import { userPlugin } from './plugins/user/user.plugin';
import { securityPlugin } from './plugins/security/security.plugin';
import 'reflect-metadata'
import { AppDataSource } from './config/data-source.config';
import { User } from './infra/entities/User';


const fastify = initFastify(
    [decorateWithUser],
    [
      { plugin: securityPlugin, routesPrefix: '/' },
      { plugin: userPlugin, routesPrefix: '/user' },
    ]
);

// TYPEORM
AppDataSource.initialize().then(async () => {

  console.log("Inserting a new user into the database...")
  const user = new User();
  user.firstName = "Timber"
  user.lastName = "Saw"
  user.age = 25
  await AppDataSource.manager.save(user)
  console.log("Saved a new user with id: " + user.id)

  console.log("Loading users from the database...")
  const users = await AppDataSource.manager.find(User)
  console.log("Loaded users: ", users)

  console.log("Here you can setup and run express / fastify / any other framework.")

}).catch(error => console.log(error))
// END TYPEORM

const start = async () => {
  try {
    await fastify.listen({ port: 3000 })
  } catch (err) {
    fastify.log.error(err)
    process.exit(1)
  }
}
start();