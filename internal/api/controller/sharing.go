package controller

import (
    "Mou1ght/internal/domain/entity"
    "Mou1ght/internal/domain/model/schema/request"
    "Mou1ght/internal/domain/model/table"
    "Mou1ght/internal/pkg/util"
    "Mou1ght/internal/repository/instance"
    "errors"
)

func CreateSharing(req *request.CreateSharingRequest) error {
    sid := util.GenSharingID()
    record := &table.SharingTable{
        PostBase: table.PostBase{
            ID:      sid,
            Content: req.Content,
        },
        AuthorID:   req.Author,
        Attachment: req.Attachment,
    }
    err := instance.UseDatabase().CreateSharing(record)
    if err != nil {
        return err
    }
    tagIDs := make([]string, len(req.Tags))
    for i, tag := range req.Tags {
        tagIDs[i] = tag.ID
    }
    err = CreateTagsLinkToSharing(tagIDs, sid)
    if err != nil {
        return err
    }
    return nil
}

func UpdateSharing(req *request.UpdateSharingRequest) error {
    record := &table.SharingTable{
        PostBase: table.PostBase{
            ID:      req.ID,
            Content: req.Content,
        },
        AuthorID:   req.Author,
        Attachment: req.Attachment,
    }
    err := instance.UseDatabase().UpdateSharing(record)
    if err != nil {
        return err
    }
    tagsIDs := make(map[string]bool)
    for _, tag := range req.Tags {
        tagsIDs[tag.ID] = true
    }
    err = instance.UseDatabase().UpdateTargetLinks(req.ID, 2, tagsIDs)
    if err != nil {
        return err
    }
    return nil
}

func ViewSharing(id string) error {
    return instance.UseDatabase().AddViewCountSharing(id)
}

func LikeSharing(id string) error {
    return instance.UseDatabase().AddLikeCountSharing(id)
}

func GetSharingByID(id string) (*entity.SharingEntity, error) {
    record, err := instance.UseDatabase().GetSharingByID(id)
    if err != nil {
        return nil, err
    }
    e := entity.NewSharingEntityFromTable(record)
    if e == nil {
        return nil, errors.New("sharing not exist")
    }
    return e, nil
}

func DeleteSharingByID(id string) error {
    err := instance.UseDatabase().DeleteSharingByID(id)
    if err != nil {
        return err
    }
    err = instance.UseDatabase().DeleteTagLinkFromTarget(id, 2)
    if err != nil {
        return err
    }
    return nil
}

func CreateTagsLinkToSharing(tags []string, sharingID string) error {
    for _, tag := range tags {
        lid := util.GenTagLinkID()
        record := &table.TagLinkTable{
            ID:         lid,
            TargetID:   sharingID,
            TargetType: 2,
            TagID:      tag,
        }
        err := instance.UseDatabase().CreateTagLink(record)
        if err != nil {
            return err
        }
    }
    return nil
}

