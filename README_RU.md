# Telegram Notify Action

Отправка уведомлений GitHub Actions в Telegram.

[🇬🇧 English version](README.md)

## Использование

```yaml
- name: Отправить уведомление
  uses: chapvic/telegram-notify@v1
  if: always()
  with:
    chat: ${{ secrets.TELEGRAM_CHAT_ID }}
    token: ${{ secrets.TELEGRAM_BOT_TOKEN }}
    status: ${{ job.status }}
```

> `if: always()` **обязателен** — без него шаг уведомления пропускается, если предыдущий шаг упал, и вы не получите сообщение о падении.

## Входные параметры

| Имя | Обязателен | По умолчанию | Описание |
|-----|-----------|--------------|----------|
| `chat` | ✅ | — | ID чата (число) или `@channel_name` |
| `token` | ✅ | — | Токен Telegram-бота |
| `status` | ❌ | `info` | Статус job: `success`, `failure`, `cancelled`, `info` |
| `title` | ❌ | `GITHUB_WORKFLOW` | Заголовок сообщения |
| `message` | ❌ | ссылка на коммит | Текст уведомления |
| `footer` | ❌ | — | Подвал сообщения |
| `parse_mode` | ❌ | `Markdown` | `Markdown`, `HTML` или `MarkdownV2` |
| `disable_notification` | ❌ | `false` | Отправить без звука |
| `disable_web_page_preview` | ❌ | `false` | Отключить превью ссылок |

## Выходные параметры

| Имя | Описание |
|-----|----------|
| `message_id` | ID отправленного сообщения в Telegram |

## Статус → иконка

| Статус | Иконка |
|--------|--------|
| `success` | ✅ |
| `failure` | 🔴 |
| `cancelled` | ❌ |
| `info` | 🔔 |

## Пример с кастомным заголовком и сообщением

```yaml
- name: Деплой
  id: deploy
  run: ./deploy.sh

- name: Уведомить Telegram
  uses: chapvic/telegram-notify@v1
  if: always()
  with:
    chat: ${{ secrets.TELEGRAM_CHAT_ID }}
    token: ${{ secrets.TELEGRAM_BOT_TOKEN }}
    status: ${{ job.status }}
    title: 'Деплой завершён'
    message: |
      Окружение: production
      Коммит: ${{ github.sha }}
    footer: '_Отправлено из GitHub Actions_'
    disable_notification: true
```

## Как это работает

Action — это Docker-контейнер, который:

1. Читает переменные окружения `INPUT_*` (GitHub Actions передаёт их из `with:`).
2. Нормализует и валидирует `status` (fallback на `info`).
3. Экранирует пользовательский ввод для защиты от Markdown-инъекций.
4. Собирает сообщение из заголовка, тела и подвала.
5. Обрезает сообщение до лимита Telegram в 4096 символов (рун), если нужно.
6. Отправляет через `https://api.telegram.org` и пишет `message_id` в `GITHUB_OUTPUT`.

## Лицензия

[MIT](LICENSE) © 2026 Chapvic
