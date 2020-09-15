import { createConnection } from 'typeorm';
import { Constants } from '../../shared/constants';

const username = process.env.POSTGRES_USER || 'postgres';
const password = process.env.POSTGRES_PASSWORD || 'postgres';

export const databaseProviders = [
  {
    provide: Constants.DATABASE_CONNECTION,
    useFactory: async () => await createConnection({
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
    }),
  },
];
