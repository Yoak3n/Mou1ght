package service

import (
	"Mou1ght/internal/domain/model/schema/request"
	"Mou1ght/internal/domain/model/table"
	"sync"

	"gorm.io/gorm"
)

var errRecordNotFound = gorm.ErrRecordNotFound

type mockArticleRepo struct {
	mu            sync.Mutex
	deleteCalls   []string
	deleteErr     error
	getByID       *table.ArticleTable
	getByIDErr    error
	articles      []*table.ArticleTable
	articleTotal  int64
	getArticlesOK bool
}

func (m *mockArticleRepo) CreateArticle(article *table.ArticleTable) error { return nil }
func (m *mockArticleRepo) UpdateArticle(article *table.ArticleTable) error { return nil }
func (m *mockArticleRepo) AddViewCountArticle(id string) error             { return nil }
func (m *mockArticleRepo) AddLikeCountArticle(id string) error             { return nil }
func (m *mockArticleRepo) GetArticleByID(id string) (*table.ArticleTable, error) {
	return m.getByID, m.getByIDErr
}
func (m *mockArticleRepo) DeleteArticleByID(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.deleteCalls = append(m.deleteCalls, id)
	return m.deleteErr
}
func (m *mockArticleRepo) GetArticlesByAuthorID(authorID string, desc bool) ([]table.ArticleTable, error) {
	return nil, nil
}
func (m *mockArticleRepo) GetArticlesByAuthorIDs(ids []string, desc bool) ([]*table.ArticleTable, error) {
	return nil, nil
}
func (m *mockArticleRepo) GetArticles(opts request.ListOptions) ([]*table.ArticleTable, int64, error) {
	return m.articles, m.articleTotal, nil
}

type mockTagRepo struct {
	deleteLinkCalls []struct {
		targetID   string
		targetType table.TagType
	}
	getTag    *table.TagTable
	getTagErr error
	updated   *table.TagTable
}

func (m *mockTagRepo) CreateTag(tag *table.TagTable) error { return nil }
func (m *mockTagRepo) DeleteTag(tagID string) error        { return nil }
func (m *mockTagRepo) GetAllTags() ([]table.TagTable, error) {
	return nil, nil
}
func (m *mockTagRepo) GetTagByID(id string) (*table.TagTable, error) {
	return m.getTag, m.getTagErr
}
func (m *mockTagRepo) GetTagsByID(ids []string) ([]table.TagTable, error) { return nil, nil }
func (m *mockTagRepo) UpdateTag(tag *table.TagTable) error {
	m.updated = tag
	return nil
}
func (m *mockTagRepo) UpdateTargetLinks(targetID string, targetType table.TagType, ids map[string]bool) error {
	return nil
}
func (m *mockTagRepo) DeleteTagWithLink(tagID string) error { return nil }
func (m *mockTagRepo) CreateTagLink(tagLink *table.TagLinkTable) error {
	return nil
}
func (m *mockTagRepo) QueryTagsByLabel(label []string) ([]table.TagTable, error) { return nil, nil }
func (m *mockTagRepo) QueryTagsByID(targetID string, targetType table.TagType) ([]table.TagTable, error) {
	return nil, nil
}
func (m *mockTagRepo) DeleteTagLink(linkID string) error    { return nil }
func (m *mockTagRepo) DeleteTagLinkByTagID(tagID string) error {
	return nil
}
func (m *mockTagRepo) DeleteTagLinkFromTarget(targetID string, targetType table.TagType) error {
	m.deleteLinkCalls = append(m.deleteLinkCalls, struct {
		targetID   string
		targetType table.TagType
	}{targetID, targetType})
	return nil
}
func (m *mockTagRepo) GetTagLinkByKeyword(keyword []string, typ string) (map[string]table.TagTable, []table.TagLinkTable, error) {
	return nil, nil, nil
}
func (m *mockTagRepo) GetArticlesFromTagLink(link *table.TagLinkTable, desc bool) ([]table.ArticleTable, error) {
	return nil, nil
}
func (m *mockTagRepo) GetSharingFromTagLink(link *table.TagLinkTable, desc bool) ([]table.SharingTable, error) {
	return nil, nil
}
func (m *mockTagRepo) CreateTagsLinkToArticle(tags []string, articleID string) error { return nil }

type mockCategoryLinkRepo struct {
	deleteByArticleCalls []string
}

func (m *mockCategoryLinkRepo) CreateCategoryLink(link *table.CategoryLinkTable) error { return nil }
func (m *mockCategoryLinkRepo) CreateCategoriesLinkToArticle(categories []string, article string) error {
	return nil
}
func (m *mockCategoryLinkRepo) DeleteCategoryLink(linkID string) error { return nil }
func (m *mockCategoryLinkRepo) DeleteCategoryLinkByArticleID(articleID string) error {
	m.deleteByArticleCalls = append(m.deleteByArticleCalls, articleID)
	return nil
}
func (m *mockCategoryLinkRepo) UpdateCategoryLinks(articleID string, categoryIDs map[string]bool) error {
	return nil
}
func (m *mockCategoryLinkRepo) GetArticlesFromCategoryLink(link *table.CategoryLinkTable, desc bool) ([]table.ArticleTable, error) {
	return nil, nil
}
func (m *mockCategoryLinkRepo) GetCategoryLinkByKeyword(keyword []string) (map[string]table.CategoryTable, []table.CategoryLinkTable, error) {
	return nil, nil, nil
}

type mockMessageRepo struct {
	getMessage    *table.MessageTable
	getMessageErr error
	updated       *table.MessageTable
	getMessages   []*table.MessageTable
	messageTotal  int64
}

func (m *mockMessageRepo) CreateMessage(record *table.MessageTable) error { return nil }
func (m *mockMessageRepo) UpdateMessage(msg *table.MessageTable) error {
	m.updated = msg
	return nil
}
func (m *mockMessageRepo) UpdateMessagePosition(id string, pos request.MessagePosition, authorIP string, isAdmin bool) error {
	return nil
}
func (m *mockMessageRepo) AddViewCountMessage(id string) error { return nil }
func (m *mockMessageRepo) AddLikeCountMessage(id string) error { return nil }
func (m *mockMessageRepo) GetMessageByID(id string) (*table.MessageTable, error) {
	return m.getMessage, m.getMessageErr
}
func (m *mockMessageRepo) DeleteMessageByID(id string) error { return nil }
func (m *mockMessageRepo) DeleteOwnMessage(id string, authorIP string) error { return nil }
func (m *mockMessageRepo) GetMessages(opts request.ListOptions) ([]*table.MessageTable, int64, error) {
	return m.getMessages, m.messageTotal, nil
}
func (m *mockMessageRepo) GetOwnedMessageIDs(authorIP string) ([]string, error) { return nil, nil }
func (m *mockMessageRepo) GetMaxZ() (int, error) { return 0, nil }

type mockSharingRepo struct {
	updated      *table.SharingTable
	existing     *table.SharingTable
	existingErr  error
	sharings     []*table.SharingTable
	sharingTotal int64
}

func (m *mockSharingRepo) CreateSharing(sharing *table.SharingTable) error { return nil }
func (m *mockSharingRepo) UpdateSharing(sharing *table.SharingTable) error {
	m.updated = sharing
	return nil
}
func (m *mockSharingRepo) AddViewCountSharing(id string) error { return nil }
func (m *mockSharingRepo) AddLikeCountSharing(id string) error { return nil }
func (m *mockSharingRepo) GetSharingsByAuthorID(authorID string, desc bool) ([]table.SharingTable, error) {
	return nil, nil
}
func (m *mockSharingRepo) GetSharings(opts request.ListOptions) ([]*table.SharingTable, int64, error) {
	return m.sharings, m.sharingTotal, nil
}
func (m *mockSharingRepo) GetSharingByID(id string) (*table.SharingTable, error) {
	if m.existingErr != nil {
		return nil, m.existingErr
	}
	if m.existing != nil {
		return m.existing, nil
	}
	return &table.SharingTable{PostBase: table.PostBase{ID: id, Status: 1}}, nil
}
func (m *mockSharingRepo) DeleteSharingByID(id string) error { return nil }

type mockAttachmentRepo struct {
	getAttachment    *table.AttachmentTable
	getAttachmentErr error
	deleted          []string
}

func (m *mockAttachmentRepo) CreateAttachment(attachment *table.AttachmentTable) error { return nil }
func (m *mockAttachmentRepo) GetAttachmentByID(id string) (*table.AttachmentTable, error) {
	return m.getAttachment, m.getAttachmentErr
}
func (m *mockAttachmentRepo) GetAttachmentsByIDs(ids []string) ([]table.AttachmentTable, error) {
	return nil, nil
}
func (m *mockAttachmentRepo) GetAttachmentBySha256(sha256 string, size int64) (*table.AttachmentTable, error) {
	return nil, nil
}
func (m *mockAttachmentRepo) ListAttachments() ([]table.AttachmentTable, error) { return nil, nil }
func (m *mockAttachmentRepo) DeleteAttachment(id string) error {
	m.deleted = append(m.deleted, id)
	return nil
}

type mockAttachmentLinkRepo struct {
	refCount int64
}

func (m *mockAttachmentLinkRepo) ReplaceSharingAttachments(sharingID string, attachmentIDs []string) error {
	return nil
}
func (m *mockAttachmentLinkRepo) DeleteBySharingID(sharingID string) error { return nil }
func (m *mockAttachmentLinkRepo) GetAttachmentIDsBySharingID(sharingID string) ([]string, error) {
	return nil, nil
}
func (m *mockAttachmentLinkRepo) CountByAttachmentID(attachmentID string) (int64, error) {
	return m.refCount, nil
}

type mockUserRepo struct {
	getUser    *table.UserTable
	getUserErr error
	byName     map[string]*table.UserTable
	updatedPwd string
	profile    map[string]any
}

func (m *mockUserRepo) CreateUser(user *table.UserTable) error { return nil }
func (m *mockUserRepo) GetUser(id string) (*table.UserTable, error) {
	return m.getUser, m.getUserErr
}
func (m *mockUserRepo) GetUserByName(name string) (*table.UserTable, error) {
	if u, ok := m.byName[name]; ok {
		return u, nil
	}
	return nil, errRecordNotFound
}
func (m *mockUserRepo) QueryUsers(username []string) ([]table.UserTable, error) { return nil, nil }
func (m *mockUserRepo) CountUsers() (int64, error)                              { return 1, nil }
func (m *mockUserRepo) UpdateUser(user *table.UserTable) error                  { return nil }
func (m *mockUserRepo) UpdateUserProfile(id string, fields map[string]any) error {
	m.profile = fields
	return nil
}
func (m *mockUserRepo) UpdateUserPassword(id, hashedPassword string) error {
	m.updatedPwd = hashedPassword
	return nil
}
func (m *mockUserRepo) DeleteUser(user *table.UserTable) error { return nil }

type mockPostRepo struct {
	lastType   string
	lastID     string
	lastStatus int8
	err        error
}

func (m *mockPostRepo) UpdatePostStatus(postType string, id string, status int8) error {
	m.lastType = postType
	m.lastID = id
	m.lastStatus = status
	return m.err
}

type mockCategoryRepo struct{}

func (m *mockCategoryRepo) CreateCategory(category *table.CategoryTable) error { return nil }
func (m *mockCategoryRepo) UpdateCategory(category *table.CategoryTable) error { return nil }
func (m *mockCategoryRepo) DeleteCategory(categoryID string) error             { return nil }
func (m *mockCategoryRepo) GetAllCategories() ([]table.CategoryTable, error)   { return nil, nil }
func (m *mockCategoryRepo) GetCategoriesByID(ids []string) ([]table.CategoryTable, error) {
	return nil, nil
}
func (m *mockCategoryRepo) QueryCategoriesByArticleID(articleID string) ([]table.CategoryTable, error) {
	return nil, nil
}
