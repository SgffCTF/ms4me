import asyncio
import logging

from fastapi import APIRouter, HTTPException, Request
from fastapi.responses import PlainTextResponse
from grpc import RpcError

from models.auth import RegisterLoginBody, RegisterResponse, LoginResponse
from grpc_client.sso.auth_pb2 import RegisterRequest, LoginRequest
        

api_v1_router = APIRouter(prefix="/api/v1")


@api_v1_router.get("/health")
async def health():    
    return PlainTextResponse("OK")


@api_v1_router.post("/register")
async def register(body: RegisterLoginBody, request: Request):    
    try:
        register_response = request.app.state.auth_client.Register(RegisterRequest(username=body.username, password=body.password))
    except RpcError as e:
        return HTTPException(status_code=400, detail=e.details())
    except Exception as e:
        logging.error("/register:" + str(e))
        return HTTPException(status_code=500, detail="Internal Server Error")
    
    return RegisterResponse(id=register_response.id)


@api_v1_router.post("/login")
async def login(body: RegisterLoginBody, request: Request):    
    try:
        login_response = request.app.state.auth_client.Login(LoginRequest(username=body.username, password=body.password))
    except RpcError as e:
        return HTTPException(status_code=400, detail=e.details())
    except Exception as e:
        logging.error("/login:" + str(e))
        return HTTPException(status_code=500, detail="Internal Server Error")
    
    return LoginResponse(token=login_response.token)
