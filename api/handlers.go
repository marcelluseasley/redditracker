package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	tmpl *template.Template

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

const (
	maxRows = 100
)

type WebSocketMessage struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func init() {
	var err error
	tmpl, err = template.ParseGlob("../../web/templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) health(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func (s *Server) homeHandler(c *gin.Context) {

	data := struct {
		Title string
		Port int
	}{
		Title: fmt.Sprintf("%s - Subreddit Tracker", s.redditClient.Config.SubReddit),
		Port: s.redditClient.Config.Port,
	}
	err := tmpl.ExecuteTemplate(c.Writer, "main-template.html", data)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing template"})
	}
}

func (s *Server) wsHandler(c *gin.Context) {

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	postRows := []string{}
	usersRows := []string{}
	rowEntry := `<tr>
	<td>%s</td>
	<td>%d</td>
	</tr>`

	go func() {
		i := 0
		j := 0

		for {

			select {
			case postCase := <-s.postDataChannel:
				for _, post := range postCase {
					postRows = append(postRows, fmt.Sprintf(rowEntry, post.Title, post.Ups))
					i++

					if i == maxRows {
						rowsAsString := strings.Join(postRows, "")

						msg := WebSocketMessage{
							Type: "post",
							Data: rowsAsString,
						}
						msgJson, err := json.Marshal(msg)
						if err != nil {
							log.Println(err)
							break
						}

						if err := ws.WriteMessage(websocket.TextMessage, msgJson); err != nil {
							log.Println(err)
							break
						}
						i = 0
						postRows = []string{}
					}
				}
			case userCase := <-s.userCountChannel:
				for _, user := range userCase {
					usersRows = append(usersRows, fmt.Sprintf(rowEntry, user.Username, user.PostCount))
					j++

					if j == maxRows {
						rowsAsString := strings.Join(usersRows, "")

						msg := WebSocketMessage{
							Type: "user",
							Data: rowsAsString,
						}
						msgJson, err := json.Marshal(msg)
						if err != nil {
							log.Println(err)
							break
						}

						if err := ws.WriteMessage(websocket.TextMessage, msgJson); err != nil {
							log.Println(err)
							break
						}
						j = 0
						usersRows = []string{}
					}

				}
			}
		}

	}()
}
