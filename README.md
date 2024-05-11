# Eager Reader Discord bot

A Discord bot which looks for URLs in messages and posts a summary of the page.

## Usage

Running the program requires two environment variables:

- `BOT_TOKEN`: The token for the Discord bot.
- `OPENAI_TOKEN`: The API key for OpenAI.

Optionally, you can customize domains to block or allow using the samples under `/config`.

Note: Empty allowlist means all domains are allowed!

### Docker container

You can run the bot as a container using the following command:

```sh
docker run -e BOT_TOKEN=secret-token -e OPENAI_TOKEN=secret-token rutkai/eager-reader-discord-bot:latest
```

If you want to define your custom allowlist/blocklist, mount a volume to `/app/config`:

```sh
-v /path/to/config:/app/config
```
