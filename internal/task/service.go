package task

import (
	"github.com/stonecool/livemusic-go/internal/database"
)

func CreateTask(category string, metaType string, metaId int, cronSpec string) (*Task, error) {
	repo := newRepositoryDB(database.DB)
	factory := newFactory(repo)
	return factory.createTask(category, metaType, metaId, cronSpec)
}

func GetAllCrawlTasks() ([]*Task, error) {
	repo := newRepositoryDB(database.DB)
	return repo.GetAll()
}

// FIXME
//func dataTypeIdExists(dataType string, dataId int) (bool, error) {
//	val, ok := internal.DataType2StructMap[dataType]
//	if !ok {
//		return false, fmt.Errorf("data_type:%s illegal", dataType)
//	}
//
//	originalType := reflect.TypeOf(val).Elem()
//	newVar := reflect.New(originalType).Elem()
//
//	pointer := newVar.Addr().Interface().(internal.IDataTable)
//	pointer.setId(dataId)
//
//	return pointer.exist()
//}