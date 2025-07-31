import { CustomFastifyRequest } from '../fastify.types';

export const getFromSession = (request: CustomFastifyRequest, key: string): any => {
  // @ts-ignore
  return request.session[key];
}

export const setToSession = (request: CustomFastifyRequest, key: string, value: any) => {
  // @ts-ignore
  request.session[key] = value;
}