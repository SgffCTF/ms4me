from config import AppConfig


if __name__ == '__main__':
    import uvicorn

    uvicorn.run("app:api", host=AppConfig.HOST, port=AppConfig.PORT, workers=AppConfig.CPU_NUM * 2 + 1)
