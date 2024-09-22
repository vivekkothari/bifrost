# bifrost

## Build docker

```
docker build -t bifrost .
docker run -p 3000:3000 bifrost
```

## Run locally

You can run the python scripts locally. You need to have python3 installed.

```
export OPENAI_API_KEY=<your_openai_api_key>
```

OpenAI SDK

```
python3 test-openai.py
```

Llama Index SDK

```
python3 test-llama-index.py
```

Langchain OpenAI SDK

```
python3 test-langchain-openai.py
```

