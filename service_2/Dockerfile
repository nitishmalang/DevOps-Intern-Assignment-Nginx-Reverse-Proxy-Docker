FROM python:3.11-slim
WORKDIR /app

RUN pip install --no-cache-dir --upgrade pip

RUN pip install --no-cache-dir uv==0.7.14

COPY pyproject.toml uv.lock app.py .

RUN uv sync --frozen

EXPOSE 8002

CMD ["uv", "run", "python", "app.py"]