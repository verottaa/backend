import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { ServeStaticModule } from '@nestjs/serve-static';
import { join } from 'path';
import { Constants } from './shared/constants';
import { UsersModule } from './modules/users/users.module';

@Module({
  imports: [
    UsersModule,
    ServeStaticModule.forRoot({
      rootPath: join(__dirname, '..', 'static'),
      exclude: ['/'.concat(Constants.globalPrefix, '*')],
    }),
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {
}
