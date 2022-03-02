package repository

import (
	"time"

	"github.com/goldalee/golangprojects/bookings/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int,error)
 	InsertRoomRestrictions(r models.RoomRestriction)error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time)([]models.Room, error)
	GetRoomByID(id int)(models.Room, error)
}
