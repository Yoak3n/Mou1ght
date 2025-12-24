package controller

import (
	"Mou1ght/internal/domain/entity"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"Mou1ght/internal/repository/instance"
	"fmt"
)

func CreateTag(req *request.CreateTagRequest) error {
	record := &table.TagTable{
		ID:    util.GenTagID(),
		Label: req.Label,
	}
	return instance.UseDatabase().DB.Create(record).Error
}

func DeleteTag(id string) error {
	if id == "" {
		return fmt.Errorf("tag id is empty")
	}
	return instance.UseDatabase().DB.Delete(&table.TagTable{ID: id}).Error
}

func TagsList() []entity.PostSign {
	records, err := instance.UseDatabase().GetAllTags()
	if err != nil {
		return nil
	}
	return entity.NewTagsInformationEntityFromTable(records)
}

// TagListWithPost 根据请求参数获取带有文章或分享的标签列表
// 参数:
//
//	req - 包含过滤条件和关键字的文章列表请求
//	isSharing - 布尔值，表示是否获取分享类型的标签
//
// 返回值:
//
//	map[string]any - 包含标签列表的映射，可能为nil(当发生错误时)
func TagListWithPost(req *request.PostListRequest, typ string) map[string]any {
	// 初始化返回结果集
	ret := make(map[string]any)
	// 获取请求中的过滤条件
	filter := req.Filter
	// 判断是否为降序排列
	descend := filter.Sort == "desc"
	// 设置默认类型为文章

	if typ == "sharing" {
		// 如果是分享类型，设置类型为分享并初始化分享标签切片
		ret["tags"] = make([]*entity.TagWithSharingEntity, 0)
	} else {
		// 否则初始化文章标签切片
		ret["tags"] = make([]*entity.TagWithArticlesEntity, 0)
	}
	// 根据关键字和类型获取标签链接
	tags, links, err := instance.UseDatabase().GetTagLinkByKeyword(req.Data.Keyword, typ)
	if err != nil {
		// 发生错误时返回nil
		return nil
	}

	// 遍历标签链接
	for i := range links {
		if typ == "sharing" {
			// 处理分享类型标签
			sharings, err := instance.UseDatabase().GetSharingFromTagLink(&links[i], descend)
			if err != nil {
				// 发生错误时跳过当前标签
				continue
			}
			// 获取标签记录并创建分享标签实体
			tagRecord := tags[links[i].TargetID]
			tag := entity.NewTagWithSharingEntityFromTable(&tagRecord, sharings)
			// 将标签添加到结果集中
			ret["tags"] = append(ret["tags"].([]*entity.TagWithSharingEntity), tag)
		} else {
			// 处理文章类型标签
			articles, err := instance.UseDatabase().GetArticlesFromTagLink(&links[i], descend)
			if err != nil {
				// 发生错误时跳过当前标签
				continue
			}
			// 获取标签记录并创建文章标签实体
			tagRecord := tags[links[i].TargetID]
			tag := entity.NewTagWithArticlesEntityFromTable(&tagRecord, articles)
			// 将标签添加到结果集中
			ret["tags"] = append(ret["tags"].([]*entity.TagWithArticlesEntity), tag)
		}
	}
	// 返回结果集
	return ret
}

func CreateTagsLinkToArticle(tags []string, articleID string) error {
	for _, tag := range tags {
		lid := util.GenTagLinkID()
		record := &table.TagLinkTable{
			ID:         lid,
			TargetID:   articleID,
			TargetType: 1,
			TagID:      tag,
		}
		err := instance.UseDatabase().CreateTagLink(record)
		if err != nil {
			return err
		}
	}
	return nil
}
