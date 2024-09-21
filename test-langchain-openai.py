import os

from langchain_openai import ChatOpenAI

llm = ChatOpenAI(
    model="gpt-4",
    temperature=0,
    max_tokens=None,
    timeout=None,
    max_retries=2,
    api_key=os.getenv("OPENAI_API_KEY"),
    base_url="http://localhost:3000",
    # organization="...",
    # other params...
)

if __name__ == "__main__":
  messages = [
    (
      "system",
      "You are a helpful assistant that translates English to French. Translate the user sentence.",
    ),
    ("human", "I love programming."),
  ]
  llm.invoke(messages)
  for chunk in llm.stream(messages):
    print(chunk)

  ai_msg = llm.invoke(messages)
  print(ai_msg)
