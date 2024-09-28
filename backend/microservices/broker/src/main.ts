import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { Logger } from "@nestjs/common";
import { FastifyAdapter, NestFastifyApplication } from "@nestjs/platform-fastify";
import { REQUEST_FIELD_ACCESS_TOKEN, REQUEST_FIELD_REFRESH_TOKEN } from "./constants";
import { ICustomRequest } from "./types";
import fastifyCookie from "@fastify/cookie";

(async () => {
  const app = await NestFactory.create<NestFastifyApplication>(AppModule, new FastifyAdapter());

  app.use((req: ICustomRequest, res, next) => {

    if (req.headers.authorization) {
      const authorizationHeader = req.headers[REQUEST_FIELD_ACCESS_TOKEN].split(" ");
      if (authorizationHeader[0] === "Bearer") {
        req.headers.authorization = authorizationHeader[1];
      }
    }
    if(req.headers.refresh_token) {
      const refreshTokenHeader = req.headers.refresh_token.split(" ");
      if (refreshTokenHeader[0] === "Bearer") {
        req.headers.refresh_token = refreshTokenHeader[1];
      }
    }

    next();
  });
  app.use(fastifyCookie)

  app.enableCors({
    origin: "http://localhost:3000",
    credentials: true,
    methods: ["GET", "POST", "PUT", "PATCH", "DELETE"],
    exposedHeaders: [REQUEST_FIELD_ACCESS_TOKEN, REQUEST_FIELD_REFRESH_TOKEN, "fingerprint", "set-cookie"]
  });

  await app.listen(8000, "0.0.0.0"); // 0.0.0.0
  Logger.log("Broker service successfully started")
})()
