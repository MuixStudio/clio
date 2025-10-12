package errorx

type ErrorX interface {
	Error() string
	Tag() string
}
