package request

type CategoryRequest struct {
	Label string `json:"label"`
	// Parent: nil 不改父级；指向空串改为根分类；指向 id 则挂到该分类下
	Parent *string `json:"parent"`
}
