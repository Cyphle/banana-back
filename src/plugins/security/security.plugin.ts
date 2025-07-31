import { FastifyInstance, FastifyPluginOptions } from 'fastify';
import * as oidcClient from 'openid-client';
import { getConfigValueOf } from '../../config/application.config';
import { loginController, logoutController, oidcCallbackController, userMeController } from './security.controller';

export const securityPlugin = async (fastify: FastifyInstance, options: FastifyPluginOptions) => {
  fastify.log.info('Initiating security plugin');
  await oidcConfiguration(fastify);

  loginController(fastify);
  oidcCallbackController(fastify);
  logoutController(fastify);
  userMeController(fastify);
}

const oidcConfiguration = async (fastify: FastifyInstance) => {
  try {
    const discoveryOptions: any = getConfigValueOf<string>('ENV', 'dev') === 'dev' ? {
      execute: [oidcClient.allowInsecureRequests]
    } : {};

    const oidcConfig = await oidcClient.discovery(
        new URL(getConfigValueOf<string>('IDP_SERVER_URL', '')),
        getConfigValueOf<string>('CLIENT_ID', ''),
        { client_secret: getConfigValueOf<string>('CLIENT_SECRET', '') },
        undefined, // clientMetadata
        discoveryOptions
    );

    fastify.decorateRequest('oidcConfig', {
      getter: () => oidcConfig
    });

    fastify.log.info('OIDC Configuration loaded successfully');
  } catch (error) {
    console.log(error);
    fastify.log.error('Failed to load OIDC configuration:', error);
    process.exit(1);
  }
}