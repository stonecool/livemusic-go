package internal

type IDataTable interface {
	setId(id int)

	exist() (bool, error)
}
