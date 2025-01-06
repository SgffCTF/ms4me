from pydantic import BaseModel, field_validator

class RegisterLoginBody(BaseModel):
    username: str
    password: str
    
    @field_validator("username")
    def validate_username(cls, v):
        if len(v) < 3:
            raise ValueError("username must be at least 3 characters long")
        if len(v) > 64:
            raise ValueError("username cannot be longer than 64 characters")
        return v
    
    @field_validator("password")
    def validate_password(cls, v):
        if len(v) < 8:
            raise ValueError("password must be at least 8 characters long")
        if len(v) > 64:
            raise ValueError("password cannot be longer than 64 characters")
        return v

class RegisterResponse(BaseModel):
    id: int

class LoginResponse(BaseModel):
    token: str
