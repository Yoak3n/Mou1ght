package service

import (
	"Mou1ght/internal/domain/entity"
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/util"
	"Mou1ght/internal/repository/interfaces"
	"fmt"
)

type TagService struct {
	tags interfaces.TagRepository
}

func NewTagService(tags interfaces.TagRepository) *TagService {
	return &TagService{tags: tags}
}

func (t *TagService) CreateTag(req *request.CreateTagRequest) error {
	record := &table.TagTable{
		ID:    util.GenTagID(),
		Label: req.Label,
	}
	return t.tags.CreateTag(record)
}

func (t *TagService) DeleteTag(id string) error {
	if id == "" {
		return fmt.Errorf("tag id is empty")
	}
	return t.tags.DeleteTagWithLink(id)
}

func (t *TagService) TagsList() []entity.PostSign {
	records, err := t.tags.GetAllTags()
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
func (t *TagService) TagListWithPost(req *request.PostListRequest, typ string) (map[string]table.TagTable, []table.TagLinkTable) {
	// 根据关键字和类型获取标签链接
	tags, links, err := t.tags.GetTagLinkByKeyword(req.Data.Keyword, typ)
	if err != nil {
		// 发生错误时返回nil
		return nil, nil
	}
	// 返回结果集
	return tags, links
}

func (t *TagService) GetSharingFromTagLink(link *table.TagLinkTable, descend bool) ([]table.SharingTable, error) {
	return t.tags.GetSharingFromTagLink(link, descend)
}

func (t *TagService) GetArticlesFromTagLink(link *table.TagLinkTable, descend bool) ([]table.ArticleTable, error) {
	return t.tags.GetArticlesFromTagLink(link, descend)
}
