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

var bots map[string] *Bot


func RegisterBots(key string, bot *Bot) {
    if bots == nil {
        bots = make(map[string] *Bot)
    }
    bots[key] = bot
}

// Returns an array of registered bots
func GetRegisteredBots() ([]*Bot) {
    list := make([]*Bot, len(bots))
    i := 0
    for _, value := range bots {
        list[i] = value
        i++
    }
    return list
}

// Starts a bot
func StartBot(botName string) {
    go bots[botName].Run()
}
