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

type GetPostsQuery struct {
	Limit uint
	Since string
	Desc  bool
	Sort  string
}

func CreateGetPostsQuery() GetPostsQuery {
	return GetPostsQuery{
		Limit: 100,
		Since: "",
		Desc:  false,
		Sort:  "flat",
	}
}

type FullPostQuery struct {
	User   bool
	Forum  bool
	Thread bool
}

func CreateFullPostQuery() FullPostQuery {
	return FullPostQuery{
		User:   false,
		Forum:  false,
		Thread: false,
	}
}
