import string
from random import choice


alphabet = string.ascii_letters + string.digits


def generate_random_string(length: int):
    result = ''
    for _ in range(length):
        result += choice(alphabet)
    return result