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

    "golang.org/x/net/websocket"
)

const rtmStart = "https://slack.com/api/rtm.start?token=%s"
const slackAPI = "https://api.slack.com/"

type rtmStartResponse struct {
    Ok bool
    Error string
    Url string
    Self responseSelf
}

type responseSelf struct {
    Id string
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
