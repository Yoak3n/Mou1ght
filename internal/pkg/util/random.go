package util

import (
	"math/rand"
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

func GenUserID() string {
	return "user-" + string(rune(time.Now().Second()))[5:] + generateRandomString(6)
}

func GenArticleID() string {
	return "art-" + string(rune(time.Now().Second()))[5:] + generateRandomString(6)
}

func GenCategoryID() string {
	return "cat-" + string(rune(time.Now().Second()))[5:] + generateRandomString(6)
}

func GenCategoryLinkID() string {
	return "clk-" + string(rune(time.Now().Second()))[5:] + generateRandomString(6)
}

func GenTagID() string {
	return "tag-" + string(rune(time.Now().Second()))[5:] + generateRandomString(6)
}

func GenTagLinkID() string {
    return "tlk-" + string(rune(time.Now().Nanosecond()))[8:] + generateRandomString(6)
}

func GenSharingID() string {
    return "sha-" + string(rune(time.Now().Second()))[5:] + generateRandomString(6)
}

func GenMessageID() string {
    return "msg-" + string(rune(time.Now().Second()))[5:] + generateRandomString(6)
}
