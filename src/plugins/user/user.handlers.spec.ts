import { Database } from '../../database/database';
import { userInfoByUsernameHandler } from './user.handlers';
import { Profile } from "../profile.types";

describe('userInfoByUsernameHandler', () => {
  let mockDatabase: jest.Mocked<Database>;

  beforeEach(() => {
    mockDatabase = {
      readOneByField: jest.fn(),
    } as any;
  });

  it('should return user info when user exists', () => {
    const mockProfile: Profile = {
      id: 1,
      username: 'testuser',
      email: 'test@example.com',
      firstName: 'Test',
      lastName: 'User',
    };

    mockDatabase.readOneByField.mockReturnValue(mockProfile);

    const result = userInfoByUsernameHandler()('testuser');

    expect(result).toEqual(mockProfile);
    expect(mockDatabase.readOneByField).toHaveBeenCalledWith('profiles', 'username', 'testuser');
  });

  it('should return undefined when user does not exist', () => {
    mockDatabase.readOneByField.mockReturnValue(undefined);

    const result = userInfoByUsernameHandler()('nonexistentuser');

    expect(result).toBeUndefined();
    expect(mockDatabase.readOneByField).toHaveBeenCalledWith('profiles', 'username', 'nonexistentuser');
  });
});
