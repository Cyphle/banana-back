import 'reflect-metadata'
import { DataSource } from 'typeorm'
import { User } from '../infra/entities/User';

export const AppDataSource = new DataSource({
  type: 'postgres',
  host: 'localhost',
  port: 5433,
  username: 'postgres',
  password: 'postgres',
  database: 'fastifyexample',
  synchronize: false,
  logging: false,
  entities: [User],
  migrations: ['migrations/*.ts'],
  subscribers: [],
})
