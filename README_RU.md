# Отправка уведомлений в Телеграм

[![Build](https://github.com/chapvic/telegram-notify/actions/workflows/build.yml/badge.svg)](https://github.com/chapvic/telegram-notify/actions/workflows/build.yml)

Отправляет уведомления о состоянии выполнения рабочих потоков в ваш Телеграм-бот с поддержкой различных Git-серверов, таких как Gitea, Gogs.

## Использование

Сначала необходимо создать Телеграм-бот с помощью [@BotFather](https://t.me/botfather). Смотрите официальное руководство: https://core.telegram.org/bots#6-botfather

Если вам необходимо получать уведомления в личный чат с помощью бота, то необходимо определить свой ID пользователя, например через бота [@jsondumpbot](https://t.me/jsondumpbot).

Также вы можете использовать каналы для уведомлений, указав название канала в формате `@channelname`.

Добавить токен вашего бота, ID пользователя (или названия Телеграм-канала) в раздел "Секреты" вашего репозитория.

Добавить следующий шаг последним в ваш рабочий поток:

```yaml
    - name: Send notification to Telegram
      uses: chapvic/telegram-notify@master
      if: always()
      with:
        token: ${{ secrets.TELEGRAM_BOT_TOKEN }} # Токен доступа к вашему боту
        chat: ${{ secrets.TELEGRAM_CHAT_ID }} # ID пользователя или название канала
        status: ${{ job.status }} # Не изменяйте это значение !!!
```

Вы можете указать необязательные аргументы в вашем рабочем потоке:

```yaml
    - name: Send notification to Telegram
      uses: chapvic/telegram-notify@master
      if: always()
      with:
        token: ${{ secrets.TELEGRAM_BOT_TOKEN }} # Токен доступа к вашему боту
        chat: ${{ secrets.TELEGRAM_CHAT_ID }} # ID пользователя или название канала
        status: ${{ job.status }} # Не изменяйте это значение !!!
        title: Заголовок вашего сообщения
        message: Основное сообщение
        footer: Подпись вашего сообщения
```

Все передваемые дополнительные аргуменеты должны содержать текст в формате Markdown
