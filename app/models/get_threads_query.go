package models

type GetThreadsQuery struct {
	Limit uint
	Since string
	Desc  bool
}

func CreateGetThreadsQuery() GetThreadsQuery {
	return GetThreadsQuery{
		Limit: 100,
		Since: "",
		Desc:  false,
	}
}
