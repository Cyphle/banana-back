import { FastifyInstance, FastifyReply } from 'fastify';
import { CustomFastifyRequest } from '../../fastify.types';
import * as oidcClient from 'openid-client';
import { getConfigValueOf } from '../../config/application.config';

export const loginController = (fastify: FastifyInstance): void => {
  fastify.get('/login', async (request: CustomFastifyRequest, reply: FastifyReply) => {
    try {
      // Génération des paramètres PKCE pour la sécurité
      const codeVerifier = oidcClient.randomPKCECodeVerifier();
      const codeChallenge = await oidcClient.calculatePKCECodeChallenge(codeVerifier);

      // Génération d'un state aléatoire pour prévenir les attaques CSRF
      const state = oidcClient.randomState();

      // Stockage des paramètres dans la session
      // @ts-ignore
      request.session.codeVerifier = codeVerifier;
      // @ts-ignore
      request.session.state = state;

      // Construction des paramètres d'autorisation
      const parameters = {
        redirect_uri: getConfigValueOf<string>('REDIRECT_URI', ''),
        scope: 'openid email profile',
        code_challenge: codeChallenge,
        code_challenge_method: 'S256',
        state: state,
        response_type: 'code',
        prompt: 'consent',
      };

      // Construction de l'URL d'autorisation
      const authorizationUrl = oidcClient.buildAuthorizationUrl(request.oidcConfig!!, parameters);

      fastify.log.info('Redirecting to authorization URL:', authorizationUrl.href);

      // Redirection vers le provider OIDC
      return reply.redirect(authorizationUrl.href);

    } catch (error) {
      fastify.log.error('Error during login initiation:', error);
      return reply.status(500).send({ error: 'Erreur lors de l\'initialisation de la connexion' });
    }
  });
}

// TODO refacto ce controller
export const oidcCallbackController = (fastify: FastifyInstance): void => {
  fastify.get('/callback', async (request: CustomFastifyRequest, reply: FastifyReply) => {
    try {
      // Récupération des paramètres stockés en session
      // @ts-ignore
      const { codeVerifier, state } = request.session;

      if (!codeVerifier || !state) {
        return reply.status(400).send({ error: 'Session invalide' });
      }

      // Construction de l'URL actuelle pour la validation
      const currentUrl = new URL(request.url, `${request.protocol}://${request.headers.host}`);

      // Échange du code d'autorisation contre les tokens
      const tokens = await oidcClient.authorizationCodeGrant(
          request.oidcConfig!!,
          currentUrl,
          {
            pkceCodeVerifier: codeVerifier,
            expectedState: state,
          }
      );

      fastify.log.info('Tokens received:', {
        access_token: tokens.access_token ? 'present' : 'missing',
        id_token: tokens.id_token ? 'present' : 'missing',
        token_type: tokens.token_type,
        expires_in: tokens.expires_in
      });

      let userInfo = null;
      let expectedSubject = null;
      if (tokens.id_token) {
        const base64Payload = tokens.id_token.split('.')[1];
        const decodedPayload = JSON.parse(Buffer.from(base64Payload, 'base64').toString());
        expectedSubject = decodedPayload.sub;

        userInfo = {
          sub: decodedPayload.sub,
          email: decodedPayload.email,
          name: decodedPayload.name,
          picture: decodedPayload.picture
        };
      }

      // Si on a un access token et un subject, essayer de récupérer les infos via UserInfo
      if (tokens.access_token && expectedSubject) {
        try {
          const userInfoResponse = await oidcClient.fetchUserInfo(
              request.oidcConfig!!,
              tokens.access_token,
              expectedSubject
          );

          // Combiner les infos de UserInfo avec celles de l'ID token
          userInfo = {
            ...userInfo,
            ...userInfoResponse
          };

          fastify.log.info('User info fetched from UserInfo endpoint');
        } catch (userInfoError) {
          fastify.log.warn('Failed to fetch user info from UserInfo endpoint:', userInfoError);
        }
      }

      // Stockage des informations utilisateur en session
      // @ts-ignore
      request.session.user = userInfo;
      // @ts-ignore
      request.session.tokens = {
        access_token: tokens.access_token,
        refresh_token: tokens.refresh_token,
        expires_at: tokens.expires_in ? Date.now() + (tokens.expires_in * 1000) : null
      };

      // Nettoyage des paramètres temporaires
      // @ts-ignore
      delete request.session.codeVerifier;
      // @ts-ignore
      delete request.session.state;

      // Redirection vers l'accueil
      return reply.redirect(getConfigValueOf<string>('FRONT_END_URL', ''));

    } catch (error: any) {
      fastify.log.error('Error during callback processing:', error);
      return reply.status(500).send({
        error: 'Erreur lors du traitement de la réponse d\'autorisation',
        details: error.message
      });
    }
  });
}

// TODO refacto ce controller et le faire marcher
export const logoutController = (fastify: FastifyInstance): void => {
  fastify.get('/logout', async (request: CustomFastifyRequest, reply: FastifyReply) => {
    try {
      // Révocation du refresh token si disponible
      // @ts-ignore
      if (request.session.tokens?.refresh_token) {
        try {
          await request.oidcConfig!!.revoke(
              request.oidcConfig!!,
              // @ts-ignore
              request.session.tokens.refresh_token,
              'refresh_token'
          );
        } catch (revokeError) {
          fastify.log.warn('Failed to revoke refresh token:', revokeError);
        }
      }

      // Destruction de la session
      // @ts-ignore
      request.session.destroy();

      return reply.send({ message: 'Déconnexion réussie' });

    } catch (error) {
      fastify.log.error('Error during logout:', error);
      return reply.status(500).send({ error: 'Erreur lors de la déconnexion' });
    }
  });
}

// TODO refacto
export const userMeController = (fastify: FastifyInstance): void => {
  fastify.get('/users/me', async (request: CustomFastifyRequest, reply: FastifyReply) => {
    // @ts-ignore
    if (!request.session.user) {
      return reply.status(401).send({
        error: 'Non autorisé',
        login_url: '/login'
      });
    }

    // @ts-ignore
    console.log(request.session.user);

    // Vérification de l'expiration du token
    // @ts-ignore
    if (request.session.tokens?.expires_at && Date.now() > request.session.tokens.expires_at) {

      // Tentative de refresh si un refresh token est disponible
      // @ts-ignore
      if (request.session.tokens.refresh_token) {
        try {
          const newTokens = await oidcClient.refreshTokenGrant(
              request.oidcConfig,
              // @ts-ignore
              request.session.tokens.refresh_token
          );

          // @ts-ignore
          request.session.tokens = {
            access_token: newTokens.access_token,
            // @ts-ignore
            refresh_token: newTokens.refresh_token || request.session.tokens.refresh_token,
            expires_at: newTokens.expires_in ? Date.now() + (newTokens.expires_in * 1000) : null
          };

          fastify.log.info('Token refreshed successfully');

        } catch (refreshError) {
          fastify.log.error('Failed to refresh token:', refreshError);
          // @ts-ignore
          request.session.destroy();
          return reply.status(401).send({
            error: 'Session expirée, veuillez vous reconnecter',
            login_url: '/login'
          });
        }
      } else {
        // @ts-ignore
        request.session.destroy();
        return reply.status(401).send({
          error: 'Session expirée, veuillez vous reconnecter',
          login_url: '/login'
        });
      }
    }

    fastify.log.info({
      message: 'Profil utilisateur',
      // @ts-ignore
      user: request.session.user,
      token_info: {
        // @ts-ignore
        expires_at: request.session.tokens?.expires_at,
        // @ts-ignore
        has_refresh_token: !!request.session.tokens?.refresh_token
      }
    })
    return reply.send({
      message: 'Profil utilisateur',
      // @ts-ignore
      user: request.session.user,
      token_info: {
        // @ts-ignore
        expires_at: request.session.tokens?.expires_at,
        // @ts-ignore
        has_refresh_token: !!request.session.tokens?.refresh_token
      }
    });
  });
}

