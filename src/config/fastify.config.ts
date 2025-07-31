import fastify, { FastifyInstance, FastifyPluginOptions } from 'fastify';
import cors from '@fastify/cors'
import fastifyCookie from '@fastify/cookie';
import fastifySession from '@fastify/secure-session';

export const initFastify = (
  decorators: ((fastify: FastifyInstance) => void)[] = [],
  plugins: { plugin: (fastify: FastifyInstance, options: FastifyPluginOptions) => void, routesPrefix: string }[],
): FastifyInstance => {
  const server = fastify({
    logger: true
  });

  server.register(cors, {});
  server.register(fastifyCookie, {});
  // @ts-ignore
  server.register(fastifySession, {
    secret: 'a-very-secure-secret-keythatisverylonglonglong', // TODO Use a proper secret key in production
    cookie: {
      path: '/',
      httpOnly: true,
      secure: false,  // TODO In production, set to true for HTTPS
    },
  });

  decorators.forEach(decorator => decorator(server));

  plugins.forEach(plugin => {
    server.register(
      plugin.plugin,
      {
        prefix: plugin.routesPrefix
      }
    );
  });

  return server;
}