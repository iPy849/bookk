package bookk

import (
	"time"
)

type BaseBooking struct {
	Id        string
	UserId    string
	ItemId    string
	CreatedAt time.Time
	StartsAt  time.Time
	EndsAt    time.Time
}

type Booking struct {
	BaseBooking
	Description string
	Cancelled   bool
}

type IBookingService[T any] interface {
	GetBookingById(bookingId string) (*T, error)
	GetLastBookingsByUserId(userId string, limit int) ([]*T, error)
	GetLastBookingsByGroupId(groupId string, limit int) ([]*T, error)
	GetBookingsByTimeRangeAndUserId(userId string, timeRange TimeRange) ([]*T, error)
	GetBookingsByTimeRangeAndGroupId(groupId string, timeRange TimeRange) ([]*T, error)
	GetBookingsByDateAndUserId(userId string, date time.Time) ([]*T, error)
	GetBookingsByDateAndGroupId(groupId string, date time.Time) ([]*T, error)
	CreateBooking(booking Booking) (*Booking, error)
	UpdateBooking(booking *Booking) error
	DeleteBooking(bookingId string) error
}
