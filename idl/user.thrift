namespace go user

struct UserRegisterRequest {
    1: string username ( vt.min_size = "1", vt.max_size = "32")
    2: string password ( vt.min_size = "1", vt.max_size = "32")
}

struct UserRegisterResponse {
    1: i32 status_code;
    2: optional string status_msg;
    3: i64 user_id;
    4: string token;
}

struct UserLoginRequest {
    1: string username (vt.min_size = "1", vt.max_size = "32")
    2: string password (vt.min_size = "1", vt.max_size = "32")
}

struct UserLoginResponse {
    1: i32 status_code;
    2: optional string status_msg;
    3: i64 user_id;
    4: string token;
}

struct UserRequest{
    1: i64 user_id;
    2: string token;
}

service UserService {
    UserRegisterResponse Register (1: UserRegisterRequest req)
    UserLoginResponse Login (1: UserLoginRequest req)
}