import { Module } from '@nestjs/common';
import { UsersController } from './users.controller';
import { UsersService } from './users.service';
import { UserRepositoryService } from './repository/user.repository.service';
import { userProviders } from './repository/user.providers';
import { UserEntity } from '../../entities/user.entity';
import { TypeOrmModule } from '@nestjs/typeorm';

@Module({
  imports: [
    TypeOrmModule.forFeature([UserEntity]),
  ],
  controllers: [UsersController],
  providers: [
    UsersService,
    UserRepositoryService,
    ...userProviders,
  ],
  exports: [],
})
export class UsersModule {
}
