import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { ServeStaticModule } from '@nestjs/serve-static';
import { join } from 'path';
import { Constants } from './shared/constants';
import { UsersModule } from './modules/users/users.module';
import { TypeOrmModule } from '@nestjs/typeorm';
import { typeormConfig } from './shared/typeorm.config';
import { Connection } from 'typeorm';

@Module({
  imports: [
    UsersModule,
    ServeStaticModule.forRoot({
      rootPath: join(__dirname, '..', 'static'),
      exclude: ['/'.concat(Constants.globalPrefix, '*')],
    }),
    TypeOrmModule.forRoot(typeormConfig),
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {
  constructor(private connection: Connection) {
  }
}
