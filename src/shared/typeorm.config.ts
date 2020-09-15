import { TypeOrmModuleOptions } from '@nestjs/typeorm/dist/interfaces/typeorm-options.interface';

const username = process.env.POSTGRES_USER || 'postgres';
const password = process.env.POSTGRES_PASSWORD || 'postgres';

export const typeormConfig: TypeOrmModuleOptions = {
  type: 'postgres',
  host: 'localhost',
  port: 5432,
  username: username,
  password: password,
  database: 'verottaa',
  synchronize: true,
  dropSchema: false,
  logging: true,
  entities: [
    __dirname + '/src/**/*.entity.ts',
    __dirname + '/dist/**/*.entity.js',
  ],
  migrations: [
    'migrations/**/*.ts',
  ],
  subscribers: [
    'subscriber/**/*.ts',
    'dist/subscriber/**/.js',
  ],
  cli: {
    entitiesDir: 'src',
    migrationsDir: 'migrations',
    subscribersDir: 'subscriber',
  },
};
