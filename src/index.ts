import { decorateWithUser } from './decorators/authentication.decorator';
import { initFastify } from './config/fastify.config';
import { userPlugin } from './plugins/user/user.plugin';
import { securityPlugin } from './plugins/security/security.plugin';

const fastify = initFastify(
    [decorateWithUser],
    [
      { plugin: securityPlugin, routesPrefix: '/' },
      { plugin: userPlugin, routesPrefix: '/user' },
    ]
);

const start = async () => {
  try {
    await fastify.listen({ port: 3000 })
  } catch (err) {
    fastify.log.error(err)
    process.exit(1)
  }
}
start();