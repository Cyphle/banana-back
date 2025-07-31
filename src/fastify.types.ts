import { FastifyRequest } from 'fastify';
import { Database } from './database/database';

export interface CustomFastifyRequest extends FastifyRequest {
  user?: string;
  database?: Database;
  oidcConfig?: any;
}
