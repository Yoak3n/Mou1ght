// 开发环境种子数据：生成一批"AI 味"的测试内容，
// 覆盖文章/说说/留言/标签/分类/账号，方便本地调试前后端各功能。
//
// 用法：
//   go run ./cmd/seed            # 数据库为空时灌入种子数据
//   go run ./cmd/seed -force     # 清空现有数据后重新灌入
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"Mou1ght/internal/config"
	"Mou1ght/internal/domain/model/table"
	"Mou1ght/internal/pkg/database"
	"Mou1ght/internal/pkg/util"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var rng = rand.New(rand.NewSource(42))

// ---- 内容池（确定性生成，不改代码则每次结果一致）----

var articleTitles = []string{
	"深入浅出聊聊现代前端构建工具",
	"Go 里那些容易踩的坑",
	"一次数据库慢查询的排查记录",
	"我为什么从自建服务器换到云主机",
	"Vue 3 组合式 API 使用心得",
	"Nginx 反代踩坑指南",
	"SQLite 性能优化小记",
	"个人博客架构演进实录",
	"Docker Compose 部署实践",
	"Markdown 编辑器选型对比",
	"关于单元测试的一点思考",
	"缓存失效那些事",
	"HTTP 状态码漫谈",
	"Git 工作流总结",
	"浅谈 JWT 鉴权与刷新策略",
	"文件上传功能的实现细节",
	"RSS 订阅的复兴",
	"暗色模式设计心得",
}

var slopParagraphs = []string{
	"在当今快速变化的技术浪潮中，保持学习的节奏显得尤为重要。每次翻开文档，都能发现一些之前被忽略的细节，这大概就是持续迭代的意义。",
	"本文记录了我从零开始搭建这套系统时的一些思考与取舍。过程中踩了不少坑，也积累了不少经验，希望这些记录对后来者有所帮助。",
	"值得一提的是，任何方案都不是银弹。在具体场景中，我们需要根据团队规模、业务复杂度和维护成本来权衡利弊，而不是盲目追随热点。",
	"代码的可维护性往往比一时的性能提升更加重要。清晰的结构、一致的命名和恰当的注释，能让后来接手的人少走很多弯路。",
	"测试是保障质量的基石。虽然编写测试会占用一些时间，但它带来的信心和安全感，往往远超投入本身。",
	"部署环节看似简单，实则需要考虑很多细节：环境变量、日志收集、监控告警、灰度发布……每一步都值得认真对待。",
	"开源社区的力量是惊人的。当我们遇到问题时，很可能已经有人遇到过并给出了解决方案，学会高效检索和阅读源码是一项重要技能。",
	"写文章的过程本身也是一次很好的梳理。把知识讲清楚，比以为自己懂了要难得多，这也是我坚持输出的原因。",
	"性能优化的前提是测量，而不是猜测。在没有数据支撑的情况下贸然优化，往往只是自欺欺人。",
	"面向未来的设计并不等于过度设计。适度的抽象能降低成本，过度的抽象则会吞噬生产力，找到那个平衡点很重要。",
}

var sharingTexts = []string{
	"今天把博客迁移到了新服务器，速度提升明显！",
	"最近在读《代码整洁之道》，推荐给每位工程师。",
	"深夜写代码的感觉真好，世界都安静了。",
	"分享一个调试技巧：打印关键变量时带上变量名，事半功倍。",
	"刚刚发布了新版本，修复了若干已知问题。",
	"周末愉快！准备去爬山放松一下。",
	"私密说说：测试私密分享功能。",
	"又一年的总结：技术、生活与成长。",
}

var messageTexts = []string{
	"博主你好，文章写得很好，学到了！",
	"请问这个主题是怎么配置的？想参考一下。",
	"感谢分享，已收藏。",
	"第一次来，做个记号。",
	"催更！什么时候更新下一篇？",
	"留言板功能很有意思，像真的便利贴一样。",
	"路过打卡，祝博主越来越好。",
	"这个暗色模式做得很舒服，赞一个。",
	"测试一下留言定位功能。",
	"审核中的留言应该不会显示吧？",
	"从搜索引擎过来的，内容很干货。",
	"期待更多关于 Go 的文章。",
	"评论区的氛围真好。",
	"这个软木板设计太可爱了。",
	"收藏了 RSS 订阅，以后常来。",
}

var categoryDefs = []struct {
	label    string
	children []string
}{
	{"技术", []string{"前端", "后端", "数据库", "运维"}},
	{"随笔", []string{"生活", "阅读"}},
	{"教程", []string{}},
	{"笔记", []string{}},
	{"公告", []string{}},
}

var tagLabels = []string{
	"Go", "Vue", "前端", "后端", "数据库", "Nginx",
	"Docker", "Git", "随笔", "经验", "教程", "性能",
}

// ---- 小工具 ----

func randInt(min, max int) int {
	if max <= min {
		return min
	}
	return min + rng.Intn(max-min+1)
}

func pick[T any](items []T) T {
	return items[rng.Intn(len(items))]
}

func randomSubset[T any](items []T, n int) []T {
	perm := rng.Perm(len(items))
	count := n
	if count > len(items) {
		count = len(items)
	}
	out := make([]T, 0, count)
	for _, i := range perm[:count] {
		out = append(out, items[i])
	}
	return out
}

// spreadTime 把第 i 项均匀分布在最近 days 天内
func spreadTime(i, total, days int) time.Time {
	if total <= 1 {
		return time.Now().Add(-time.Hour)
	}
	day := time.Duration((total-1-i)*days/(total-1)) * 24 * time.Hour
	return time.Now().Add(-day).Add(-time.Duration(rng.Intn(3600)) * time.Second)
}

func buildArticleContent(title string) string {
	paras := randomSubset(slopParagraphs, randInt(3, 6))
	content := "# " + title + "\n\n"
	for _, p := range paras {
		content += p + "\n\n"
	}
	content += "## 小结\n\n" + pick(slopParagraphs) + "\n"
	return content
}

func wipe(db *gorm.DB) {
	tables := []any{
		&table.TagLinkTable{},
		&table.CategoryLinkTable{},
		&table.AttachmentLinkTable{},
		&table.MessageTable{},
		&table.SharingTable{},
		&table.ArticleTable{},
		&table.TagTable{},
		&table.CategoryTable{},
		&table.AttachmentTable{},
		&table.UserTable{},
	}
	for _, m := range tables {
		if err := db.Unscoped().Where("1 = 1").Delete(m).Error; err != nil {
			log.Fatalf("清空 %T 失败: %v", m, err)
		}
	}
	log.Println("已清空旧数据")
}

// ---- 数据写入 ----

func createUser(db *gorm.DB) *table.UserTable {
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	u := &table.UserTable{
		ID:        util.GenUserID(),
		UserName:  "admin",
		Password:  string(hash),
		Email:     "admin@mou1ght.local",
		Avatar:    "",
		Role:      1,
		CreatedAt: time.Now(),
	}
	if err := db.Create(u).Error; err != nil {
		log.Fatalf("创建管理员失败: %v", err)
	}
	return u
}

func createCategories(db *gorm.DB) []*table.CategoryTable {
	var out []*table.CategoryTable
	for _, def := range categoryDefs {
		root := &table.CategoryTable{ID: util.GenCategoryID(), Label: def.label, CreatedAt: time.Now()}
		if err := db.Create(root).Error; err != nil {
			log.Fatalf("创建分类 %s 失败: %v", def.label, err)
		}
		out = append(out, root)
		for _, child := range def.children {
			c := &table.CategoryTable{ID: util.GenCategoryID(), Label: child, ParentID: root.ID, CreatedAt: time.Now()}
			if err := db.Create(c).Error; err != nil {
				log.Fatalf("创建子分类 %s 失败: %v", child, err)
			}
			out = append(out, c)
		}
	}
	return out
}

func createTags(db *gorm.DB) []*table.TagTable {
	var out []*table.TagTable
	for _, label := range tagLabels {
		t := &table.TagTable{ID: util.GenTagID(), Label: label, CreatedAt: time.Now()}
		if err := db.Create(t).Error; err != nil {
			log.Fatalf("创建标签 %s 失败: %v", label, err)
		}
		out = append(out, t)
	}
	return out
}

func createArticles(db *gorm.DB, admin *table.UserTable, categories []*table.CategoryTable, tags []*table.TagTable) []*table.ArticleTable {
	const total = 18
	var out []*table.ArticleTable
	for i := 0; i < total; i++ {
		title := articleTitles[i%len(articleTitles)]
		status := int8(1)
		switch {
		case i%6 == 4:
			status = 0 // draft
		case i%6 == 5:
			status = 2 // archive
		}
		createdAt := spreadTime(i, total, 180)
		a := &table.ArticleTable{
			PostBase: table.PostBase{
				ID:        util.GenArticleID(),
				Content:   buildArticleContent(title),
				Like:      int64(randInt(5, 60)),
				View:      int64(randInt(120, 800)),
				Status:    status,
				CreatedAt: createdAt,
				UpdatedAt: createdAt,
			},
			Title:    title,
			AuthorID: admin.ID,
		}
		if err := db.Create(a).Error; err != nil {
			log.Fatalf("创建文章失败: %v", err)
		}
		// 标签链接（1-2 个 / 文章）
		for _, t := range randomSubset(tags, randInt(1, 2)) {
			link := &table.TagLinkTable{
				ID:         util.GenTagLinkID(),
				TagID:      t.ID,
				TargetID:   a.ID,
				TargetType: table.ArticleTag,
				CreatedAt:  createdAt,
			}
			if err := db.Create(link).Error; err != nil {
				log.Fatalf("创建文章标签链接失败: %v", err)
			}
		}
		// 分类链接（已发布文章 1 个分类）
		if status == 1 {
			c := pick(categories)
			link := &table.CategoryLinkTable{
				ID:         util.GenCategoryLinkID(),
				ArticleID:  a.ID,
				CategoryID: c.ID,
				CreatedAt:  createdAt,
			}
			if err := db.Create(link).Error; err != nil {
				log.Fatalf("创建文章分类链接失败: %v", err)
			}
		}
		out = append(out, a)
	}
	return out
}

func createSharings(db *gorm.DB, admin *table.UserTable, tags []*table.TagTable) []*table.SharingTable {
	const total = 8
	var out []*table.SharingTable
	for i := 0; i < total; i++ {
		status := int8(1)
		if i == 6 {
			status = 0 // 私密
		}
		createdAt := spreadTime(i, total, 60)
		s := &table.SharingTable{
			PostBase: table.PostBase{
				ID:        util.GenSharingID(),
				Content:   sharingTexts[i%len(sharingTexts)],
				Like:      int64(randInt(0, 30)),
				View:      int64(randInt(20, 300)),
				Status:    status,
				CreatedAt: createdAt,
				UpdatedAt: createdAt,
			},
			AuthorID: admin.ID,
		}
		if err := db.Create(s).Error; err != nil {
			log.Fatalf("创建说说失败: %v", err)
		}
		if status == 1 {
			for _, t := range randomSubset(tags, randInt(0, 2)) {
				link := &table.TagLinkTable{
					ID:         util.GenTagLinkID(),
					TagID:      t.ID,
					TargetID:   s.ID,
					TargetType: table.SharingTag,
					CreatedAt:  createdAt,
				}
				if err := db.Create(link).Error; err != nil {
					log.Fatalf("创建说说标签链接失败: %v", err)
				}
			}
		}
		out = append(out, s)
	}
	return out
}

func createMessages(db *gorm.DB) []*table.MessageTable {
	const total = 15
	var out []*table.MessageTable
	for i := 0; i < total; i++ {
		status := int8(1)
		switch {
		case i >= 10 && i < 13:
			status = 3 // pending review
		case i >= 13:
			status = 2 // archive
		}
		createdAt := spreadTime(i, total, 30)
		m := &table.MessageTable{
			PostBase: table.PostBase{
				ID:        util.GenMessageID(),
				Content:   messageTexts[i%len(messageTexts)],
				Like:      int64(randInt(0, 15)),
				View:      int64(randInt(5, 120)),
				Status:    status,
				CreatedAt: createdAt,
				UpdatedAt: createdAt,
			},
			X:        randInt(3, 75),
			Y:        randInt(5, 70),
			Z:        i,
			AuthorIP: fmt.Sprintf("seed-visitor-%d", i+1),
		}
		if err := db.Create(m).Error; err != nil {
			log.Fatalf("创建留言失败: %v", err)
		}
		out = append(out, m)
	}
	return out
}

func main() {
	force := flag.Bool("force", false, "清空现有数据后重新灌入种子数据")
	flag.Parse()

	config.GetConfig()
	db := database.InitDatabase()

	var userCount int64
	db.Model(&table.UserTable{}).Count(&userCount)
	if userCount > 0 && !*force {
		log.Println("数据库已有数据，跳过种子灌入。如需清空重灌请使用 -force")
		return
	}
	if *force {
		wipe(db)
	}

	admin := createUser(db)
	categories := createCategories(db)
	tags := createTags(db)
	articles := createArticles(db, admin, categories, tags)
	sharings := createSharings(db, admin, tags)
	messages := createMessages(db)

	fmt.Printf(`
种子数据写入完成！
  管理员账号: admin / admin123
  分类: %d 个（含子分类）
  标签: %d 个
  文章: %d 篇（含草稿与归档）
  说说: %d 条（含 1 条私密）
  留言: %d 条（含待审核与归档）
`, len(categories), len(tags), len(articles), len(sharings), len(messages))
}
