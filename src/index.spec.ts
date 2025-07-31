'use strict';

import fastify from 'fastify';
import { describe, expect, test } from '@jest/globals';

function build(opts = {}) {
  const app = fastify(opts)
  app.get('/', async function (request, reply) {
    return { hello: 'world' }
  })

  return app
}

describe('Fastify integration test example', () => {
  const fakeApp = async (): Promise<{ status: number, body: string }> => {
    const app = build()

    const response = await app.inject({
      method: 'GET',
      url: '/'
    })

    return Promise.resolve({ status: response.statusCode, body: response.body });
  }

  test('should return 200 status code and "hello world" body', async () => {
    const { status, body } = await fakeApp()
    expect(status).toBe(200)
    expect(body).toBe('{"hello":"world"}')
  });

  test('another example', (done) => {
    build()
      .inject()
      .get('/')
      .headers({ foo: 'bar' })
      .query({ foo: 'bar' })
      .end((err, res) => {
        expect(res?.statusCode).toEqual(200);
        done();
      });
  });
});
