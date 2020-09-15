import { Connection, Repository } from 'typeorm';
import { UserEntity } from '../../../entities/user.entity';
import { Constants } from '../../../shared/constants';

export const userProviders = [
  {
    provide: Constants.USER_REPOSITORY,
    useFactory: (connection: Connection) => connection.getRepository(UserEntity),
    inject: [Constants.DATABASE_CONNECTION],
  },
];
