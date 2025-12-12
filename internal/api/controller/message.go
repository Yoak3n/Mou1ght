package controller

import (
    "Mou1ght/internal/domain/entity"
    "Mou1ght/internal/domain/model/schema/request"
    "Mou1ght/internal/domain/model/table"
    "Mou1ght/internal/pkg/util"
    "Mou1ght/internal/repository/instance"
    "errors"
    "time"
)

func CreateMessage(req *request.CreateMessageRequest) error {
    mid := util.GenMessageID()
    record := &table.MessageTable{
        PostBase: table.PostBase{
            ID:      mid,
            Content: req.Content,
        },
        X:        req.Position.X,
        Y:        req.Position.Y,
        Z:        req.Position.Z,
        AuthorIP: req.AuthorIP,
    }
    return instance.UseDatabase().CreateMessage(record)
}

func UpdateMessage(req *request.UpdateMessageRequest) error {
    record := &table.MessageTable{
        PostBase: table.PostBase{
            ID:      req.ID,
            Content: req.Content,
        },
        X:        req.Position.X,
        Y:        req.Position.Y,
        Z:        req.Position.Z,
        AuthorIP: req.AuthorIP,
    }
    return instance.UseDatabase().UpdateMessage(record)
}

func ViewMessage(id string) error {
    return instance.UseDatabase().AddViewCountMessage(id)
}

func LikeMessage(id string) error {
    return instance.UseDatabase().AddLikeCountMessage(id)
}

func GetMessageByID(id string) (*entity.MessageEntity, error) {
    record, err := instance.UseDatabase().GetMessageByID(id)
    if err != nil {
        return nil, err
    }
    e := entity.NewMessageEntityFromTable(record)
    if e == nil {
        return nil, errors.New("message not exist")
    }
    return e, nil
}

func DeleteMessageByID(id string) error {
    return instance.UseDatabase().DeleteMessageByID(id)
}

func ListMessages(dateRange *request.PostFilterDate, sort string) ([]*table.MessageTable, error) {
    var startDate, endDate *time.Time
    if dateRange != nil {
        s, _ := time.Parse("2006-01-02 15:04:05", dateRange.StartDate)
        e, _ := time.Parse("2006-01-02 15:04:05", dateRange.EndDate)
        startDate = &s
        endDate = &e
    }
    msgs, err := instance.UseDatabase().GetMessages(startDate, endDate)
    if err != nil {
        return nil, err
    }
    if sort == "desc" {
        for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
            msgs[i], msgs[j] = msgs[j], msgs[i]
        }
    }
    return msgs, nil
}
