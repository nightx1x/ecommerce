package order

import "errors"

var (
	//Order validate errors
	ErrOrderIDRequired         = errors.New("order id is required")
	ErrShippingAddressRequired = errors.New("shipping address is required")
	ErrPaymentMethodInvalid    = errors.New("invalid payment method")
	ErrStatusRequired          = errors.New("status is required")
	ErrUserIDRequired          = errors.New("user id is required")
	ErrItemIDRequired          = errors.New("item id is required")
	ErrEmptyQuantity           = errors.New("quantity must be greater than 0")
	ErrShippingAddrReq         = errors.New("shipping address is required")
	ErrInvalidPayment          = errors.New("invalid payment method")
	ErrInvalidStatus           = errors.New("invalid order status")
	ErrOrderEmpty              = errors.New("order has no items")

	//Order item errors
	ErrOrderMustContainItem   = errors.New("order must contain at least one item")
	ErrProductIDRequired      = errors.New("product id is required")
	ErrInvalidProductQuantity = errors.New("product quantity must be greater than 0")

	//Order status errors
	ErrOrderAlreadyCanceled  = errors.New("order already canceled")
	ErrCannotCancelDelivered = errors.New("cannot cancel a delivered order")

	//logic errors
	ErrOrderNotFound     = errors.New("order not found")
	ErrProductNotFound   = errors.New("product not found")
	ErrItemNotFound      = errors.New("cart item not found")
	ErrInsufficientStock = errors.New("insufficient stock for product")
	ErrCannotCancelPaid  = errors.New("cannot cancel a paid or shipped order")
)
