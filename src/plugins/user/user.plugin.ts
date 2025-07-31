import { FastifyError, FastifyInstance, FastifyPluginOptions, FastifyReply, FastifyRequest } from 'fastify';
import { logoutController, userInfoController } from './user.controller';
import { userInfoByUsernameHandler } from './user.handlers';

export const userPlugin = async (fastify: FastifyInstance, options: FastifyPluginOptions) => {
  fastify.log.info('Initiating user plugin');

  fastify.setErrorHandler((error: FastifyError, request: FastifyRequest, reply: FastifyReply) => {
    throw error
  })

  userInfoController(userInfoByUsernameHandler)(fastify);
  logoutController()(fastify);
}