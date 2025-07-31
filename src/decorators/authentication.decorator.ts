import { FastifyInstance, FastifyReply, HookHandlerDoneFunction } from 'fastify';
import {CustomFastifyRequest} from "../fastify.types";

export const decorateWithUser = (fastify: FastifyInstance): void => {
  fastify.log.info('Decorating with user');
  fastify.decorateRequest('user', null);

  fastify.addHook('preHandler', (req: CustomFastifyRequest, reply: FastifyReply, done: HookHandlerDoneFunction) => {
    req.user = '';
    done()
  });
}