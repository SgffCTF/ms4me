import grpc
from grpc_client.sso import auth_pb2_grpc
from typing import AsyncGenerator
from config import AppConfig


class GrpcClient:
    def __init__(self, host=AppConfig.SSO_HOST, port=AppConfig.SSO_PORT):
        self.addr = f'{host}:{port}'
        self.grpc_channel = None
        self.auth_client = None

    async def __aenter__(self):
        self.grpc_channel = grpc.insecure_channel(self.addr)
        self.auth_client = auth_pb2_grpc.AuthStub(self.grpc_channel)
        return self.auth_client

    async def __aexit__(self, exc_type, exc_val, exc_tb):
        pass

