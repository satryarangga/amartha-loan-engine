package models

type SortDirection string

const (
	SortDirectAscending  SortDirection = "asc"
	SortDirectDescending SortDirection = "desc"
)

type FindAllParam struct {
	Limit          int
	Offset         int
	SearchKeyword  string
	FieldsToSearch []string
	PreloadTables  []string
	JoinTables     []string
	SortBy         SortBy
}

type SortBy struct {
	FieldName string
	Direction SortDirection
}
