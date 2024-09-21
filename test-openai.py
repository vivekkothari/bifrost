import os

from openai import OpenAI

client = OpenAI(
    base_url="http://localhost:3000",
    api_key=os.getenv("OPEN_API_KEY")
    # this is also the default, it can be omitted
)


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
    stream = True
    response = client.chat.completions.create(
        model="gpt-4",
        messages=[
          {
            "role": "system",
            "content": "You are a helpful assistant that translates English to French. Translate the user sentence.",
          },
          {
            "role": "user",
            "content": prompt,
          },
        ],
        stream=stream,
    )
    if stream:  # Check if response is iterable (indicating streaming)
      for chunk in response:
        chunk_json = chunk.to_json()
        print(chunk_json)  # Print each chunk for debugging
    else:
      # Handle non-streaming response
      print(response.to_json())

  except Exception as e:
    return f"Error: {str(e)}"


if __name__ == "__main__":
  # Example prompt
  prompt = "I love programming."

  # Generate text
  result = generate_text(prompt)
  print(result)
