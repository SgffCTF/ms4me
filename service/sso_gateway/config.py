from os import getenv


class AppConfig:
    ENV = getenv('ENV') or 'local'
    DEBUG = bool(getenv('DEBUG') or False)
    SSO_HOST = getenv('SSO_HOST') or '127.0.0.1'
    SSO_PORT = int(getenv('SSO_PORT') or 15001)
    HOST = getenv('HOST') or '127.0.0.1'
    PORT = int(getenv('PORT') or 15002)
    CPU_NUM = int(getenv('CPU_NUM') or 4)
