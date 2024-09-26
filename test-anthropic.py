import anthropic

if __name__ == "__main__":
  # Example prompt
  client = anthropic.Anthropic(
      # defaults to os.environ.get("ANTHROPIC_API_KEY")
      base_url="http://localhost:3000"
  )
  message = client.messages.create(
      model="claude-3-5-sonnet-20240620",
      max_tokens=1024,
      messages=[
        {"role": "user", "content": "Hello, Claude"}
      ]
  )
  print(message.content)
