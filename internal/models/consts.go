package models

const (
	ContentTypePost    = "posts"
	ContentTypeArticle = "articles"
	ContentTypeNews    = "news"
)

const (
	LevelAll    = ""
	LevelEasy   = "easy"
	LevelMedium = "medium"
	LevelHard   = "hard"
)

const (
	ViewsAll     = ""
	ViewsNone    = "rated0"
	ViewsFew     = "rated10"
	ViewsSome    = "rated25"
	ViewsMany    = "rated50"
	ViewsPopular = "rated100"
)

const (
	PeriodDaily   = "daily"
	PeriodWeekly  = "weekly"
	PeriodMonthly = "monthly"
	PeriodYearly  = "yearly"
	PeriodAlltime = "alltime"
)

const (
	SortNew = ""
	SortTop = "top"
)

const (
	SearchSortRelevance = "relevance"
	SearchSortDate      = "date"
	SearchSortRate      = "rating"
)
