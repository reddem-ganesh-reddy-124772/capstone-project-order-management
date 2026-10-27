
package model

// These constants define the possible states of a stock reservation.

// InventoryReservationStatusReserved means stock has been deducted
// from the product's available stock.
const InventoryReservationStatusReserved = "reserved"

// InventoryReservationStatusReleased means the reserved stock has
// been restored. Repeated release requests must not restore it again.
const InventoryReservationStatusReleased = "released"
