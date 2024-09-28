import { Module } from "@nestjs/common";
import { ClientsModule, Transport } from "@nestjs/microservices";
import { MAIL_PACKAGE_NAME, MAIL_SERVICE_NAME } from "./proto/mail";
import { join } from "path"
import { SESSION_PACKAGE_NAME, SESSION_SERVICE_NAME } from "./proto";
import { STORE_PACKAGE_NAME, STORE_SERVICE_NAME } from "./proto/store";
import { TASKS_PACKAGE_NAME, TASKS_SERVICE_NAME } from "./proto/tasks";
import { USER_PACKAGE_NAME, USER_SERVICE_NAME } from "./proto/user";
import { TaskController, UserController } from "./controllers";

@Module({
    imports: [
        ClientsModule.register([
            {
                name: MAIL_SERVICE_NAME,
                transport: Transport.GRPC,
                options: {
                    package: MAIL_PACKAGE_NAME,
                    protoPath: join(process.cwd(), `src/proto/mail.proto`)
                }
            },
            {
                name: SESSION_SERVICE_NAME,
                transport: Transport.GRPC,
                options: {
                    package: SESSION_PACKAGE_NAME,
                    protoPath: join(process.cwd(), `src/proto/session.proto`)
                }
            },
            {
                name: STORE_SERVICE_NAME,
                transport: Transport.GRPC,
                options: {
                    package: STORE_PACKAGE_NAME,
                    protoPath: join(process.cwd(), `src/proto/store.proto`)
                }
            },
            {
                name: TASKS_SERVICE_NAME,
                transport: Transport.GRPC,
                options: {
                    package: TASKS_PACKAGE_NAME,
                    protoPath: join(process.cwd(), `src/proto/tasks.proto`)
                }
            },
            {
                name: USER_SERVICE_NAME,
                transport: Transport.GRPC,
                options: {
                    package: USER_PACKAGE_NAME,
                    protoPath: join(process.cwd(), `src/proto/user.proto`)
                }
            },
        ])
    ],
  controllers: [UserController, TaskController],
})
export class AppModule {}
