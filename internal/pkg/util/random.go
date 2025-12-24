package util

import (
	"math/rand"
	"strconv"
	"time"
)

func generateRandomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func nowSecondToString() string {
	s := strconv.FormatInt(time.Now().Unix(), 10)
	return s
}

func nowMilliSecondToString() string {
	s := strconv.FormatInt(time.Now().UnixMilli(), 10)
	return s
}

func GenUserID() string {
	return "user-" + nowSecondToString()[5:] + generateRandomString(6)
}

func GenArticleID() string {
	return "art-" + nowSecondToString()[5:] + generateRandomString(6)
}

func GenCategoryID() string {
	return "cat-" + nowSecondToString()[5:] + generateRandomString(6)
}

func GenCategoryLinkID() string {
	return "clk-" + nowMilliSecondToString()[8:] + generateRandomString(6)
}

func GenTagID() string {
	return "tag-" + nowSecondToString()[5:] + generateRandomString(6)
}

func GenTagLinkID() string {
	return "tlk-" + nowMilliSecondToString()[8:] + generateRandomString(6)
}

func GenSharingID() string {
	return "sha-" + nowSecondToString()[5:] + generateRandomString(6)
}

func GenMessageID() string {
	return "msg-" + nowSecondToString()[5:] + generateRandomString(6)
}
