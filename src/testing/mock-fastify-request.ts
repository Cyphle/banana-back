import { FastifyRequest } from 'fastify';
import { Session } from '@fastify/secure-session';

export const mockFastifyRequest = () => ({
  session: {
    user: undefined,
    changed: false,
    deleted: false,
    get: jest.fn(),
    set: jest.fn(),
    data: jest.fn(),
    delete: jest.fn(),
    options: jest.fn(),
    touch: jest.fn(),
    regenerate: jest.fn(),
  } as Session,
} as FastifyRequest);
