// Copyright © 2017 David Harrigan
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bot

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "regexp"
    "strings"
    "sync/atomic"

    "golang.org/x/net/websocket"
)

const rtmStart = "https://slack.com/api/rtm.start?token=%s"
const slackAPI = "https://api.slack.com/"
var counter uint64


type Bot struct {
    Token string
    BotId string
    Commands []Command
    socket *websocket.Conn
}

type rtmStartResponse struct {
    Ok bool
    Error string
    Url string
    Self responseSelf
}

type responseSelf struct {
    Id string
}

type Message struct {
    Id uint64 `json:"id"`
    Type string `json:"type"`
    Channel string `json:"channel"`
    Text string `json:"text"`
}

type Command struct {
    Regex *regexp.Regexp
    Mention bool
    Action func(Message) Message
}

/**
 *  Starts a connection with Slack
 */
func connectToSlack(token string) (ws *websocket.Conn, id string, err error) {
    // Start a RTM connection
    url := fmt.Sprintf(rtmStart, token)
    response, err := http.Get(url)
    if err != nil {
        return
    }
    if response.StatusCode != 200 {
        err = fmt.Errorf("API request failed with code %d", response.StatusCode)
        return
    }

    // Grab the response body
    body, err := ioutil.ReadAll(response.Body)
    response.Body.Close()
    if err != nil {
        return
    }

    // JSON-fy
    var responseObj rtmStartResponse
    err = json.Unmarshal(body, &responseObj)
    if err != nil {
        return
    }

    if !responseObj.Ok {
        err = fmt.Errorf("Slack error: %s", responseObj.Error)
        return
    }

    websocketUrl := responseObj.Url
    id = responseObj.Self.Id

    ws, err = websocket.Dial(websocketUrl, "", slackAPI)
    if err != nil {
        return
    }

    return
}

/**
 * Run the bot
 */
func (b *Bot) Run() {
    ws, id, err := connectToSlack(b.Token)
    if err != nil {
        log.Fatal(err)
    }
    b.socket = ws
    b.BotId = id

    c := make(chan Message)
    go b.Listen(c)

    // Main loop
    for {
        message := <-c
        for _, cmd := range b.Commands {
            // This command explicitly expects mention, but the actual message didn't contain any
            if cmd.Mention && !strings.Contains(message.Text, "<@"+b.BotId+">") {
                continue
            }
            // Match!
            if cmd.Regex.MatchString(message.Text) {
                log.Println("Processing command:", message.Text)
                response := cmd.Action(message)
                log.Println("Responding with:", response)
                err = b.SendMessage(response)
                if err != nil {
                    log.Fatal(err)
                }
            }
        }
    }
}

/**
 * Continuously listen for new incoming messages
 */
func (b *Bot) Listen(c chan Message) {
    for {
        message, err := b.ReceiveMessage()
        if err != nil {
            log.Fatal(err)
        } else {
            log.Println("Received: ", message)
            c<-message
        }
    }
}

/**
 * Receive message from the websocket.
 */
func (b *Bot) ReceiveMessage() (message Message, err error) {
    err = websocket.JSON.Receive(b.socket, &message)
    return
}

/**
 * Send message via the websocket.
 */
func (b *Bot) SendMessage(message Message) (err error) {
    message.Id = atomic.AddUint64(&counter, 1)
    err = websocket.JSON.Send(b.socket, message)
    return
}

/**
 * Register a new command
 */
func (b *Bot) RegisterCommand(command Command) {
    b.Commands = append(b.Commands, command)
}
