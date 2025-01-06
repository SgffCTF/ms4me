FROM python:3.12-slim

RUN adduser --disabled-password -u 1000 user

USER user

WORKDIR /sso_gateway

ENV PIP_DISABLE_PIP_VERSION_CHECK=1

COPY requirements.txt .

RUN pip3 install --no-cache-dir -r requirements.txt
RUN pip3 install --no-cache-dir gunicorn

COPY --chown=user:user grpc_client/ ./grpc_client/
COPY --chown=user:user handlers/ ./handlers/
COPY --chown=user:user models/ ./models/
COPY --chown=user:user app.py config.py wsgi.py /sso_gateway/

CMD ["python3", "wsgi.py"]