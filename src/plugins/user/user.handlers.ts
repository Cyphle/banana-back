import { UserInfo } from './user.types';

export const userInfoByUsernameHandler = () => (username: string): UserInfo | undefined => {
  return {
    username: 'johndoe',
    email: 'johndoe@example.fr',
    firstName: 'John',
    lastName: 'Doe'
  }
}
