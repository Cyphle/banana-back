import { FastifyInstance, FastifyReply, HookHandlerDoneFunction } from 'fastify';
import { CustomFastifyRequest } from '../fastify.types';
import { database } from "../config/database.config";

export const decorateWithDatabase = (fastify: FastifyInstance): void => {
  fastify.log.info('Decorating with database');
  fastify.decorateRequest('database', null);
  fastify.addHook('preHandler', (req: CustomFastifyRequest, reply: FastifyReply, done: HookHandlerDoneFunction) => {
    req.database = database
    done()
  });
}