package db

type (
	Option func(*options)

	options struct {
		limit  int
		offset int
		orders []string
	}
	sortOrder string
)

var DESC sortOrder = "desc"
var ACS sortOrder = "acs"

func NewOptions() *options {
	return &options{
		limit:  10,
		offset: -1,
	}
}

func WithLimit(limit int) Option {
	return func(opt *options) {
		opt.limit = limit
	}
}

func WithOffset(offset int) Option {
	return func(opt *options) {
		opt.offset = offset
	}
}

func WithOrder(field string, r sortOrder) Option {
	return func(opt *options) {
		opt.orders = append(opt.orders, field+" "+string(r))
	}
}
