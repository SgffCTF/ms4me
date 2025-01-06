from fastapi import FastAPI
import logging
from grpc_client.client import GrpcClient

from handlers.api_v1 import api_v1_router
from config import AppConfig


def init_app() -> FastAPI:
    config = AppConfig()
    
    docs_kwargs = {}
    if config.ENV == 'prod':
        logging.basicConfig(level=logging.INFO)
        docs_kwargs = {"docs_url": None, "redoc_url": None}
    elif config.ENV == 'local':
        logging.basicConfig(level=logging.DEBUG)
    
    async def lifespan(app: FastAPI):
        grpc_client = GrpcClient()
        async with grpc_client as client:
            app.state.auth_client = client
            yield
        grpc_client.grpc_channel.close()

    app = FastAPI(
        debug=config.DEBUG,
        title="SSO Gateway",
        description="API Gateway for SSO",
        version="1.0.0",
        lifespan=lifespan,
        **docs_kwargs
    )

    app.include_router(api_v1_router)

    return app


api = init_app()
