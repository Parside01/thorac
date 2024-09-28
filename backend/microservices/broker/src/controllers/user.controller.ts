import {
    Body,
    Controller,
    ForbiddenException, Get,
    Inject,
    Put,
    Query,
    Req,
    Res,
    UnauthorizedException
} from "@nestjs/common";
import { FindUsersDTO, LoginUserDTO, User, USER_SERVICE_NAME, UserServiceClient, UserState } from "../proto/user";
import { ConfirmAccessCodeDTO, MAIL_SERVICE_NAME, MailServiceClient } from "../proto/mail";
import { SESSION_SERVICE_NAME, SessionServiceClient, VerifyTokensDTO } from "../proto";
import { AUTH_COOKIE_OPTIONS, COOKIE_REFRESH_TOKEN_NAME } from "../constants";
import { ICustomRequest, ICustomResponse } from "../types";
import { DisableAuth } from "../disable-auth.decorator";

@Controller()
export class UserController {
    constructor(
        @Inject(USER_SERVICE_NAME) private userService: UserServiceClient,
        @Inject(MAIL_SERVICE_NAME) private mailService: MailServiceClient,
        @Inject(SESSION_SERVICE_NAME) private sessionService: SessionServiceClient,
    ) {}

    @Put("user/login")
    @DisableAuth()
    async loginUser(@Body() data: LoginUserDTO): Promise<{user_id: string}> {
        const isLogin = await this.userService.loginUser({email: data.email, password: data.password}).toPromise()

        if(!isLogin.isLogined) throw new UnauthorizedException("Неверный пароль!")

        await this.mailService.sendAccessCode({user_id: isLogin.user.id, email: data.email}).toPromise()
        return {user_id: isLogin.user.id}
    }

    @Put("user/confirm")
    @DisableAuth()
    async confirmUser(@Query() data: ConfirmAccessCodeDTO, @Req() req: ICustomRequest, @Res({ passthrough: true }) res: ICustomResponse): Promise<VerifyTokensDTO> {
        if(!req.headers.fingerprint) throw new UnauthorizedException("Недостатчоно данных для создания сессии!")

        const { isConfirmed } = await this.mailService.confirmAccessCode({user_id: data.user_id, code: data.code}).toPromise()
        const ip = req.ip || req.socket.remoteAddress || req.headers["x-forwarded-for"]

        if(!isConfirmed) throw new ForbiddenException()

        await this.userService.updateUser({id: data.user_id, state: UserState.ACTIVE}).toPromise()

        const tokens = await this.sessionService.generateTokens({
            ua: req.headers["user-agent"],
            fingerprint: req.headers["fingerprint"],
            user_id: data.user_id,
            ip
        }).toPromise()

        res.setCookie(COOKIE_REFRESH_TOKEN_NAME, tokens.refresh_token, AUTH_COOKIE_OPTIONS)

        return tokens
    }

    @Get("user/me")
    async getUserBySessionTokens(@Req() req: ICustomRequest): Promise<User> {
        return await this.userService.findUser({id: req.user_id}).toPromise()
    }

    @Get("user/search")
    async searchUsers(@Req() req: ICustomRequest, @Query() query: FindUsersDTO): Promise<User[]> {
        const { users } = await this.userService.findUsers(query).toPromise()
        const withoutMe = users?.filter((user: User) => user.id !== req.user_id)

        return withoutMe ?? []
    }
}