package utils

import (
    "time"
	"fmt"
    "math/rand"
)

var profileImgsNameList = []string{
    "Garfield", "Tinkerbell", "Annie", "Loki", "Cleo", "Angel", "Bob", "Mia", "Coco",
    "Gracie", "Bear", "Bella", "Abby", "Harley", "Cali", "Leo", "Luna", "Jack", "Felix", "Kiki",
}

var profileImgsCollectionsList = []string{
    "notionists-neutral", "adventurer-neutral", "fun-emoji",
}

func GetDefaultProfileImg() string {
    rand.Seed(time.Now().UnixNano())
    collection := profileImgsCollectionsList[rand.Intn(len(profileImgsCollectionsList))]
    name := profileImgsNameList[rand.Intn(len(profileImgsNameList))]
    return fmt.Sprintf("https://api.dicebear.com/6.x/%s/svg?seed=%s", collection, name)
}