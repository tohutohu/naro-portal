package model

import (
	"fmt"

	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/pborman/uuid"
)

//ChangePin Pinの構造体
type ChangePin struct {
	TweetID string `json:"tweetID,omitempty"`
}

//PostPinHandler Post /pin ピン
func PostPinHandler(c echo.Context) error {
	pin := ChangePin{}
	c.Bind(&pin)

	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}

	var userID string
	Db.Get(&userID,"SELECT user_ID FROM tweet WHERE tweet_ID=?",pin.TweetID)
	if userID!=sess.Values["UserID"]{
		return c.String(http.StatusInternalServerError, "あなたのTweetではありません")
	}

	Db.Exec("INSERT INTO pin (pin_ID, user_ID,tweet_ID) VALUES (?, ?,?)", uuid.New(), sess.Values["UserID"], pin.TweetID)
	return c.NoContent(http.StatusOK)
}

//PostDeletePinHandler Post /pinDelete Pin消去
func PostDeletePinHandler(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}

	pin := ChangePin{}
	c.Bind(&pin)

	Db.Exec("DELETE FROM pin WHERE user_ID=? AND tweet_ID=?", sess.Values["UserID"], pin.TweetID)
	return c.NoContent(http.StatusOK)
}
