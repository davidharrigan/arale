# arale
Simple Slack bot library for Go.  Create your next Slack bot with ease, and with __style__.

## Usage

__Example__
```go
package main

import (
    "regexp"
    "github.com/davidharrigan/arale/bot"
)

// This will be triggered by our fancy bot!
func respond(m bot.Message)(bot.Message) {
    m.Text = "World!"
    return m
}


func main() {
    // Instantiate bot - needs a token. Use the below URL to create your own bot
    // https://my.slack.com/services/new/bot
    b := bot.Bot{Token:"slack-api-token"}

    // Add some commands and actions
    r, _ := regexp.Compile("Hello")
    command := bot.Command{
        Regex: r,
        Action: respond  // Must accept a Message object and return a Message object
        Mention: false  // If true, the command will only be processed if the bot is explicitly mentioned
    }
    b.RegisterCommand(command)

    // Run the bot
    b.Run()
}
```
