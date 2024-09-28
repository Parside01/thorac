export const COOKIE_REFRESH_TOKEN_NAME = "refresh_token"
export const REQUEST_FIELD_REFRESH_TOKEN = "session"
export const REQUEST_FIELD_ACCESS_TOKEN = "authorization"
export const AUTH_COOKIE_OPTIONS = {httpOnly: true, secure: true, maxAge: 60*60*24*14, path: "/"}
export const USE_AUTH_METADATA = "UseAuthMetadata"