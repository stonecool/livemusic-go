package internal

type IDataTable interface {
	SetId(id int)

	Exist() (bool, error)
}
