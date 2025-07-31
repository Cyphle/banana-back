import { Database } from './database';
import { Profile } from "../plugins/profile.types";

describe('In memory database', () => {
  test('should get mock data and store them in memory', () => {
    const database = new Database();

    expect(Object.keys(database.dump())).toEqual(['accounts', 'profiles']);
  });

  test('should get elements from given table', () => {
    const profiles = new Database().read<Profile>('profiles');

    expect(profiles).toEqual([{
      'id': 1,
      'username': 'john.doe',
      'firstName': 'John',
      'lastName': 'Doe',
      'email': 'john.doe@banana.fr'
    }]);
  });

  test('should get element from given table for given identifier key and value', () => {
    const profiles = new Database().readOneById<Profile>('profiles', 1);

    expect(profiles).toEqual({
      'id': 1,
      'username': 'john.doe',
      'firstName': 'John',
      'lastName': 'Doe',
      'email': 'john.doe@banana.fr'
    });
  });

  test('should get element from given table for given search field and criteria key and value', () => {
    const profiles = new Database().readOneByField<Profile>('profiles', 'email', 'john.doe@banana.fr');

    expect(profiles).toEqual({
      'id': 1,
      'username': 'john.doe',
      'firstName': 'John',
      'lastName': 'Doe',
      'email': 'john.doe@banana.fr'
    });
  });

  test('should create element for given table', () => {
    const database = new Database();

    database.create('profiles', {
      'username': 'john.smith',
      'firstName': 'John',
      'lastName': 'Smith',
      'email': 'john.smith@banana.fr'
    });

    expect(database.read('profiles')).toEqual([
      {
        'id': 1,
        'username': 'john.doe',
        'firstName': 'John',
        'lastName': 'Doe',
        'email': 'john.doe@banana.fr'
      },
      {
        'id': 2,
        'username': 'john.smith',
        'firstName': 'John',
        'lastName': 'Smith',
        'email': 'john.smith@banana.fr'
      }
    ]);
  });

  test('should update an element for given table and id', () => {
    const database = new Database();

    database.update('profiles', 1, {
      'username': 'john.smith',
      'firstName': 'John',
      'lastName': 'Smith',
      'email': 'john.smith@banana.fr'
    });

    expect(database.read('profiles')).toEqual([
      {
        'id': 1,
        'username': 'john.smith',
        'firstName': 'John',
        'lastName': 'Smith',
        'email': 'john.smith@banana.fr'
      }
    ])
  });

  test('should delete an element for given table and id', () => {
    const database = new Database();

    database.delete('profiles', 1);

    expect(database.read('profiles')).toEqual([]);
  });
});