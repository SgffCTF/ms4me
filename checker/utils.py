import string
from random import randint, choice
from requests import Session


alphabet = string.ascii_letters + string.digits
user_agents = [f'python-requests/2.{x}.0' for x in range(15, 28)]


def generate_random_string(length: int):
    result = ''
    for _ in range(length):
        result += alphabet[randint(0, len(alphabet) - 1)]
    return result


def session_with_ua():
        sess = Session()
        sess.headers['User-Agent'] = choice(user_agents)
        return sess