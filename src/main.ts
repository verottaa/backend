import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { Constants } from './shared/constants';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  app.setGlobalPrefix(Constants.globalPrefix);
  await app.listen(3000);
}

bootstrap();
