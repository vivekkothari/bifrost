import os

from llama_index.core.llms import ChatMessage
from llama_index.llms.openai import OpenAI


def generate_text(prompt: str, model: str = "gpt-3.5-turbo-instruct",
    max_tokens: int = 100) -> str:
  """
  Generates text using the newer OpenAI API with chat-based completions.

  Parameters:
  - prompt: The input text to guide the model's generation.
  - model: The model to use (default: gpt-3.5-turbo).
  - max_tokens: The maximum number of tokens to generate (default: 100).

  Returns:
  - Generated text as a string.
  """
  try:
    # Use the newer chat-based API
    messages = [
      ChatMessage(
          role="system",
          content="You are a helpful assistant that translates English to French. Translate the user sentence."
      ),
      ChatMessage(role="user", content=prompt),
    ]
    llm = OpenAI(
        api_key=os.environ.get("OPENAI_API_KEY"),
        max_tokens=max_tokens,
        api_base="http://localhost:3000",
        api_version="v1",
        model=model, )
    resp = llm.chat(messages)
    print(resp)

    resp = llm.stream_chat(messages)
    for r in resp:
      print(r.delta + "\n")
  except Exception as e:
    return f"Error: {str(e)}"


if __name__ == "__main__":
  # Example prompt
  prompt = "I love programming."

  # Generate text
  result = generate_text(prompt)
  print(result)
