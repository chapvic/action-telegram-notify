# Telegram Notify Action

Send GitHub Actions notifications to Telegram.

[🇷🇺 Русская версия](README_RU.md)

## Usage

```yaml
- name: Send notification
  uses: chapvic/telegram-notify@v1
  if: always()
  with:
    chat: ${{ secrets.TELEGRAM_CHAT_ID }}
    token: ${{ secrets.TELEGRAM_BOT_TOKEN }}
    status: ${{ job.status }}
```

> `if: always()` is **required** — without it, the notification step is skipped when a previous step fails, so you won't be notified about failures.

## Inputs

| Name | Required | Default | Description |
|------|----------|---------|-------------|
| `chat` | ✅ | — | Chat ID (numeric) or `@channel_name` |
| `token` | ✅ | — | Telegram Bot Token |
| `status` | ❌ | `info` | Job status: `success`, `failure`, `cancelled`, `info` |
| `title` | ❌ | `GITHUB_WORKFLOW` | Message title |
| `message` | ❌ | link to commit | Notification message |
| `footer` | ❌ | — | Message footer |
| `parse_mode` | ❌ | `Markdown` | `Markdown`, `HTML`, or `MarkdownV2` |
| `disable_notification` | ❌ | `false` | Send silently |
| `disable_web_page_preview` | ❌ | `false` | Disable link previews |

## Outputs

| Name | Description |
|------|-------------|
| `message_id` | ID of the sent Telegram message |

## Status → icon

| Status | Icon |
|--------|------|
| `success` | ✅ |
| `failure` | 🔴 |
| `cancelled` | ❌ |
| `info` | 🔔 |

## Example with custom title and message

```yaml
- name: Deploy
  id: deploy
  run: ./deploy.sh

- name: Notify Telegram
  uses: chapvic/telegram-notify@v1
  if: always()
  with:
    chat: ${{ secrets.TELEGRAM_CHAT_ID }}
    token: ${{ secrets.TELEGRAM_BOT_TOKEN }}
    status: ${{ job.status }}
    title: 'Deployment finished'
    message: |
      Environment: production
      Commit: ${{ github.sha }}
    footer: '_Sent by GitHub Actions_'
    disable_notification: true
```

## How it works

The action is a Docker container that:

1. Reads `INPUT_*` environment variables (provided by GitHub Actions from `with:`).
2. Normalizes and validates `status` (falls back to `info`).
3. Escapes all user input to prevent Markdown injection.
4. Builds the message with title, body, and footer.
5. Truncates to Telegram's 4096-rune limit if needed.
6. Sends via `https://api.telegram.org` and writes `message_id` to `GITHUB_OUTPUT`.

## License

[MIT](LICENSE) © 2026 Chapvic
