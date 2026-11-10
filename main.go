package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/JohannesKaufmann/html-to-markdown/escape"
    "github.com/go-telegram/bot"
    "github.com/go-telegram/bot/models"
)

const (
    maxMessageRunes = 4096
)

var (
    icons = map[string]string{
        "success":   "✅",
        "failure":   "🔴",
        "cancelled": "❌",
        "info":      "🔔",
    }

    // Arguments list
    token                    = os.Getenv("INPUT_TOKEN")                    // required
    chat                     = os.Getenv("INPUT_CHAT")                     // required
    status                   = os.Getenv("INPUT_STATUS")                   // recommended
    title                    = os.Getenv("INPUT_TITLE")                    // optional
    message                  = os.Getenv("INPUT_MESSAGE")                  // optional
    footer                   = os.Getenv("INPUT_FOOTER")                   // optional
    parseModeRaw             = os.Getenv("INPUT_PARSE_MODE")               // optional
    disableNotificationRaw   = os.Getenv("INPUT_DISABLE_NOTIFICATION")     // optional
    disableWebPagePreviewRaw = os.Getenv("INPUT_DISABLE_WEB_PAGE_PREVIEW") // optional

    // GitHub environment variables
    actor    = os.Getenv("GITHUB_ACTOR")      // who's made a commit
    server   = os.Getenv("GITHUB_SERVER_URL") // git-server
    workflow = os.Getenv("GITHUB_WORKFLOW")   // workflow name
    repo     = os.Getenv("GITHUB_REPOSITORY") // repository name
    commit   = os.Getenv("GITHUB_SHA")        // commit sha

    // Format template for title
    fmtTitle = "%s %s *%s*" // icon, message, status
)

// prepareChatID returns int64 for numeric chat IDs, or string for @channel_name
func prepareChatID(chat string) any {
    if id, err := strconv.ParseInt(chat, 10, 64); err == nil {
        return id
    }
    return chat
}

// normalizeStatus lowercases and trims status, falling back to "info".
func normalizeStatus(s string) (string, string) {
    s = strings.ToLower(strings.TrimSpace(s))
    if s == "" {
        return "info", "Status is not given! Set to 'info'..."
    }
    if _, ok := icons[s]; !ok {
        return "info", fmt.Sprintf("Invalid status %q! Set to 'info'...", s)
    }
    return s, ""
}

// truncateMessage cuts s to at most limit runes, appending "…" on truncation
func truncateMessage(s string, limit int) string {
    r := []rune(s)
    if len(r) <= limit {
        return s
    }
    return string(r[:limit-1]) + "…"
}

// buildTitle composes the notification title from icon, title/workflow and status
func buildTitle(icon, title, status, workflow string) string {
    titleSource := workflow
    if title != "" {
        titleSource = title
    }
    return fmt.Sprintf(fmtTitle, icon, titleSource, status)
}

// buildCommitLine composes the commit link, or "" if GitHub env is incomplete
func buildCommitLine(server, repo, commit string) string {
    if server == "" || repo == "" || commit == "" {
        return ""
    }
    return fmt.Sprintf("💬 Commit: [open link](%v/%v/commit/%v)", server, repo, commit)
}

// resolveParseMode maps INPUT_PARSE_MODE to models.ParseMode, default Markdown
func resolveParseMode(raw string) models.ParseMode {
    switch strings.ToLower(strings.TrimSpace(raw)) {
    case "html":
        return models.ParseMode("HTML")
    case "markdownv2":
        return models.ParseMode("MarkdownV2")
    default:
        return models.ParseMode("Markdown")
    }
}

// resolveBool parses "true"/"1"/"yes" as true (case-insensitive), anything else as false
func resolveBool(raw string) bool {
    switch strings.ToLower(strings.TrimSpace(raw)) {
    case "true", "1", "yes":
        return true
    default:
        return false
    }
}

// writeOutput appends KEY=VALUE to GITHUB_OUTPUT (or logs if not set)
func writeOutput(key, value string) {
    path := os.Getenv("GITHUB_OUTPUT")
    if path == "" {
        log.Printf("    - GITHUB_OUTPUT not set, skipping output %q=%q", key, value)
        return
    }
    f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        warning(fmt.Sprintf("failed to open GITHUB_OUTPUT: %v", err))
        return
    }
    defer f.Close()
    if _, err := fmt.Fprintf(f, "%s=%s\n", key, value); err != nil {
        warning(fmt.Sprintf("failed to write GITHUB_OUTPUT: %v", err))
    }
}

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
    if s {
        fatal("Notification was interrupted!")
    }

    chatID := prepareChatID(chat)

    // Validate status and build icon
    log.Printf("    - Check status message and create icon")
    status, statusWarn := normalizeStatus(status)
    if statusWarn != "" {
        warning(statusWarn)
    }
    icon := icons[status]

    // Resolve parse mode and flags
    parseMode := resolveParseMode(parseModeRaw)
    disableNotification := resolveBool(disableNotificationRaw)
    disableWebPagePreview := resolveBool(disableWebPagePreviewRaw)

    // Escape everything that goes into Markdown/HTML
    safeStatus   := escape.MarkdownCharacters(status)
    safeWorkflow := escape.MarkdownCharacters(workflow)
    safeTitle    := escape.MarkdownCharacters(title)
    safeMessage  := escape.MarkdownCharacters(message)
    safeFooter   := escape.MarkdownCharacters(footer)
    safeServer   := escape.MarkdownCharacters(server)
    safeRepo     := escape.MarkdownCharacters(repo)
    safeCommit   := escape.MarkdownCharacters(commit)
    safeActor    := escape.MarkdownCharacters(actor)

    // Build title
    log.Printf("    - Create notification title")
    msg := buildTitle(icon, safeTitle, safeStatus, safeWorkflow)

    // Append message
    log.Printf("    - Prepare notification message")
    if safeMessage != "" {
        msg = fmt.Sprintf("%v\n%v", msg, safeMessage)
    } else if line := buildCommitLine(safeServer, safeRepo, safeCommit); line != "" {
        warning("No message given! Using default notification message...")
        msg = fmt.Sprintf("%v\n%v", msg, line)
    } else {
        warning("No message given and GitHub env is incomplete — message body skipped")
    }

    // Append footer (or default "by @actor" when footer is empty)
    if safeFooter != "" {
        msg = fmt.Sprintf("%v\n%v", msg, safeFooter)
    } else if safeActor != "" {
        msg = fmt.Sprintf("%v\n_by @%v_", msg, safeActor)
    }

    // Enforce Telegram length limit
    if truncated := truncateMessage(msg, maxMessageRunes); truncated != msg {
        warning(fmt.Sprintf("Message exceeds %d runes — truncated", maxMessageRunes))
        msg = truncated
    }

    // Send
    log.Printf("📢 Sending notification message...")
    b, err := bot.New(token)
    if err != nil {
        fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    sent, err := b.SendMessage(ctx, &bot.SendMessageParams{
        ChatID:              chatID,
        Text:                msg,
        ParseMode:           parseMode,
        DisableNotification: disableNotification,
        LinkPreviewOptions:  &models.LinkPreviewOptions{IsDisabled: &disableWebPagePreview},
    })
    if err != nil {
        fatal(err)
    }

    log.Printf("✅ Success! message_id=%d", sent.ID)
    writeOutput("message_id", strconv.Itoa(sent.ID))
}

// fatal logs a fatal message and terminates
func fatal(msg any) {
    log.Fatalf("❗ Fatal: %v", msg)
}

// warning logs a warning message
func warning(msg any) {
    log.Printf("⚠️ Warning: %v", msg)
}
