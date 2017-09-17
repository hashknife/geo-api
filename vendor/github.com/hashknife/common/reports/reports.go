package reports

// PDFer
type PDFer interface {
	PDF() error
}

// PackagesDeliveredByCourier
type PackagesDeliveredByCourier struct{}

// PackagesDeliveredByAccount
type PackagesDeliveredByAccount struct{}

// DeliveryRequestsByAccount
type DeliveryRequestsByAccount struct{}

// DeliveryRequestsOverall
type DeliveryRequestsOverall struct{}
