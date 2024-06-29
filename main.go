package main

import (
    "fmt"
    "log"
    "os"
    "strings"
    "strconv"
    "context"

    "github.com/go-telegram/bot"
    "github.com/go-telegram/bot/models"
    "github.com/JohannesKaufmann/html-to-markdown/escape"
)

var (
    icons = map[string] string {
        "success":   "✅",
        "failure":   "🔴",
        "cancelled": "❌",
        "info":      "🔔",
    }

    // Init zero value for ChatID
    chat_id     = 0

    // Arguments list
    token       = os.Getenv("INPUT_TOKEN")      // required
    chat        = os.Getenv("INPUT_CHAT")       // required
    status      = os.Getenv("INPUT_STATUS")     // recommended
    title       = os.Getenv("INPUT_TITLE")      // optional
    message     = os.Getenv("INPUT_MESSAGE")    // optional
    footer      = os.Getenv("INPUT_FOOTER")     // optional

    // Github enviroment variables
    actor       = os.Getenv("GITHUB_ACTOR")             // who's made a commit
    server      = os.Getenv("GITHUB_SERVER_URL")        // git-server
    workflow    = os.Getenv("GITHUB_WORKFLOW")          // workflow name
    repo        = os.Getenv("GITHUB_REPOSITORY")        // repository name
    commit      = os.Getenv("GITHUB_SHA")               // commit message

    // Format template for title
    fmt_title = "%s %s *%s*"                    // icon, message, status
)


func main() {

    log.Printf("🚀 Starting Telegram Notify...")

    s := false
    if token == "" {
        log.Printf("    ⚡ Token is required!")
        s = true
    }
    if chat == "" {
        log.Printf("    ⚡ Chat ID is required!")
        s = true
    }
    if s == true {
        fatal("Notification was interrupted!")
    }

    // Prepare ChatID numeric value (int64)
    chat_id, err := strconv.ParseInt(chat, 10, 64)
    if err != nil {
        fatal(err)
    }

    // Make status icon text
    log.Printf("    - Check status message and create icon")
    if status == "" {
        warning("Status is not given! Set to 'info'...")
        status = "info"
    } else if status != "success" && status != "failure" && status != "cancelled" {
        warning(fmt.Sprintf("Invalid status %v! Set to 'info'...", status))
        status = "info"
    }
    icon := icons[strings.ToLower(status)]

    workflow = escape.MarkdownCharacters(workflow)
    title = escape.MarkdownCharacters(title)
    message = escape.MarkdownCharacters(message)
    footer = escape.MarkdownCharacters(footer)

    // Make title text
    m_title := ""
    log.Printf("    - Create notification title")
    if title == "" {
        m_title = fmt.Sprintf(fmt_title, icon, workflow, status)
    } else {
        m_title = fmt.Sprintf(fmt_title, icon, title, status)
    }
    msg := fmt.Sprintf("%s", m_title)

    // Append message text
    log.Printf("    - Prepare notification message")
    if message != "" {
        msg = fmt.Sprintf("%v\n%v", msg, message)
    } else {
        warning("No message given! Using default notification message...")
        m_commit := fmt.Sprintf("💬 Commit: [%v](%v/%v/commit/%v)", "open link", server, repo, commit)
        msg = fmt.Sprintf("%v\n%v", msg, m_commit)
    }

    // Append footer text
    if footer != "" {
        msg = fmt.Sprintf("%v\n%v", msg, footer)
    }

    // Sending notification message
    log.Printf("📢 Sending notification message...")
    opts := []bot.Option{}
    b, err := bot.New(token, opts...)
    if err != nil {
        fatal(err)
    }

    _, err = b.SendMessage(context.TODO(), &bot.SendMessageParams{
        ChatID: chat_id,
        Text:   msg,
        ParseMode: models.ParseModeMarkdown,
    })
    if err != nil {
        fatal(err)
    }

    log.Printf("✅ Success!")

}


// Output fatal message and terminate
func fatal(msg any) {
    log.Fatal(fmt.Sprintf(`❗ Fatal: %s`, msg))
}


// Output warning message
func warning(msg any) {
    log.Printf(fmt.Sprintf(`⚠️ Warning: %s`, msg))
}
