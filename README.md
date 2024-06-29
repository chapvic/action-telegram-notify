# Send Notifications to Telegram

[![Build](https://github.com/chapvic/telegram-notify/actions/workflows/build.yml/badge.svg)](https://github.com/chapvic/telegram-notify/actions/workflows/build.yml)

Send workflow status notifications to Telegram with support different Git-servers (i.e. Gitea, Gogs)

## Usage

First of all, you need to create a Telegram bot by talking to [@BotFather](https://t.me/botfather) bot. See official guide here: https://core.telegram.org/bots#6-botfather

If you want to get notifications to personal chat with bot, find your user id, for example by talking to [@jsondumpbot](https://t.me/jsondumpbot).

Also you can use channel for notifications, in this case just get your channel name in format `@channelname`.

Then add your bot token and user id (or channel name) to repository Secrets.

Add following minimal step to the end of your workflow:

```yaml
    - name: Send notification to Telegram
      uses: chapvic/telegram-notify@master
      if: always()
      with:
        token: ${{ secrets.TELEGRAM_BOT_TOKEN }} # Token secret
        chat: ${{ secrets.TELEGRAM_CHAT_ID }} # User ID or channel name secret
        status: ${{ job.status }} # Do not modify this line !!!
```

You can specify optional arguments in your workflow:

```yaml
    - name: Send notification to Telegram
      uses: chapvic/telegram-notify@master
      if: always()
      with:
        token: ${{ secrets.TELEGRAM_BOT_TOKEN }} # Token secret
        chat: ${{ secrets.TELEGRAM_CHAT_ID }} # User ID or channel name secret
        status: ${{ job.status }} # Do not modify this line !!!
        title: Some workflow title
        message: Your notification text message
        footer: Footer message
```

All additional arguments must be in the form of markdown text
