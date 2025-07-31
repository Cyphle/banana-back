import { mockFastify } from '../../testing/mock-fastify';
import { logoutController, userInfoController } from './user.controller';
import { userInfoByUsernameHandler } from './user.handlers';

describe('User controller', () => {
  test('should get user info for given username', (done) => {
    mockFastify({}, [userInfoController(userInfoByUsernameHandler)])
      .inject()
      .get('/info')
      .query({ username: 'john.doe' })
      .end((err, res) => {
        expect(res?.statusCode).toEqual(200);
        expect(JSON.parse(res?.body ?? '{}')).toEqual({ id: 1, username: 'john.doe', firstName: 'John', lastName: 'Doe', email: 'john.doe@banana.fr' });
        done();
      });
  });

  test('should logout user', (done) => {
    mockFastify({}, [logoutController()])
      .inject()
      .get('/logout')
      .end((err, res) => {
        expect(res?.statusCode).toEqual(200);
        done();
      });
  });
});
