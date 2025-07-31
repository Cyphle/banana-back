import fastify, {FastifyInstance, FastifyReply, HookHandlerDoneFunction} from 'fastify';
import {CustomFastifyRequest} from '../fastify.types';
import fastifyCookie from '@fastify/cookie';
import fastifySession from '@fastify/secure-session';

export const mockFastify = (opts = {}, controllers: ((fastify: FastifyInstance) => void)[]) => {
  const app: FastifyInstance = fastify(opts);
  app.register(fastifyCookie, {});
  // @ts-ignore
  app.register(fastifySession, {
    secret: 'verylonglongsecretthatisonlyfortests',
    cookie: {
      path: '/',
      httpOnly: true,
      secure: false, 
    },
  });

  app.addHook('preHandler', (req: CustomFastifyRequest, reply: FastifyReply, done: HookHandlerDoneFunction) => {
    // @ts-ignore
    req.session.set('user', { username: 'john.doe' });
    done()
  });

  controllers.forEach(controller => controller(app));

  return app
}