import { FastifyRequest } from 'fastify';

export const getNumberParam = (request: FastifyRequest, paramName: string): number => {
  return parseInt((request.params as { [key: string]: string })[paramName]);
}

export const getStringQuery = (request: FastifyRequest, queryName: string): string => {
  return (request.query as { [key: string]: string })[queryName];
}

export const getStringBodyElement = <T>(request: FastifyRequest, elementName: string): T => {
  return (request.body as { [key: string]: T })[elementName];
}