package account

import (
	"fmt"

	"github.com/stonecool/livemusic-go/internal/config"
	"gorm.io/gorm"
)

type factoryImpl struct {
	repo IRepository
}

func NewFactory(repo IRepository) IFactory {
	return &factoryImpl{repo: repo}
}

func (f *factoryImpl) CreateAccount(category string) (*Account, error) {
	v := NewValidator()
	if err := v.validateCategory(category); err != nil {
		return nil, fmt.Errorf("invalid account category: %w", err)
	}

	if _, ok := config.AccountMap[category]; !ok {
		return nil, fmt.Errorf("unsupported account category: %s", category)
	}

	model := &model{Category: category}
	account := model.toEntity()

	if err := v.ValidateAccount(account); err != nil {
		return nil, fmt.Errorf("invalid account: %w", err)
	}

	// return f.repo.Transaction(func(r IRepository) error {
	// return r.Create(account)
	// })

	if err := f.repo.Create(account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	account.Init()
	return account, nil
}

// 在 internal/account/factory.go 中添加一个便捷的创建方法
func CreateAccount(db *gorm.DB, category string) (*Account, error) {
	repo := NewRepositoryDB(db)
	factory := NewFactory(repo)
	return factory.CreateAccount(category)
}

// 例如在 API handler 中
// func CreateAccountHandler(c *gin.Context) {
//     category := c.PostForm("category")

//     account, err := account.CreateAccount(db, category)
//     if err != nil {
//         c.JSON(400, gin.H{"error": err.Error()})
//         return
//     }

//     // 初始化账号
//     account.Init()

//     c.JSON(200, account)
// }

//func getCrawl(id int) (interface{}, error) {
//	account := &Account{ID: id}
//	err := account.Get()
//	if err != nil {
//		return nil, err
//	}
//
//	factory, err := GetFactory(account.Category)
//	if err != nil {
//		return nil, err
//	}
//
//	cfg, ok := config.AccountMap[account.Category]
//	if !ok {
//		return nil, fmt.Errorf("config not found for category: %s", account.Category)
//	}
//
//	crawl := factory.CreateCrawl(&cfg)
//	go startCrawl(crawl)
//	return crawl, nil
//}
