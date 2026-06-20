package orders

import (
	"testing"

	"github.com/DavoodHakimi/warehouse-app/internal/audit"
	"github.com/DavoodHakimi/warehouse-app/internal/company"
	"github.com/DavoodHakimi/warehouse-app/internal/partners"
	"github.com/DavoodHakimi/warehouse-app/internal/products"
	"github.com/DavoodHakimi/warehouse-app/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// setupService wires a Service over a fresh in-memory DB that contains all
// tables referenced by Order (including audit.Log, since UpdateOrder records
// change entries). Returns the service, the db handle (for seeding), and a
// helper that creates an order row in a given state.
func setupService(t *testing.T) (*Service, *gorm.DB) {
	t.Helper()
	db := testutil.NewTestDB(t,
		&Order{}, &OrderItem{}, &Currency{},
		&company.Company{},
		&partners.BusinessPartner{}, &partners.BusinessPartnerType{},
		&products.Product{}, &products.Stock{},
		&audit.Log{},
	)
	repo := NewRepository(db)
	return NewService(repo), db
}

// seedOrder creates and persists an Order with the given type+status.
// (seedRefs, makeOrder and itoa are defined in repository_test.go.)
func seedOrder(t *testing.T, db *gorm.DB, orderType, status string) Order {
	t.Helper()
	compID, bpID, curID := seedRefs(t, db)
	o := makeOrder(compID, bpID, curID, orderType, status)
	o.OrderNumber = "ORD-" + itoa(uint(statusSeq()))
	require.NoError(t, db.Create(&o).Error)
	return o
}

var statusCounter int

func statusSeq() int { statusCounter++; return statusCounter }

func TestService_CreateOrder(t *testing.T) {
	svc, db := setupService(t)
	compID, bpID, curID := seedRefs(t, db)
	prod := products.Product{Name: "Widget", ProductNumber: "P1", CompanyID: compID}
	require.NoError(t, db.Create(&prod).Error)

	req := &CreateOrderRequest{
		OrderType:         "sale",
		BusinessPartnerID: bpID,
		CurrencyID:        curID,
		ExchangeRate:      1.0,
		OrderItems: []OrderItemReq{
			{ProductID: prod.ID, Quantity: 2, PerItemPrice: 10},
		},
	}
	require.NoError(t, svc.CreateOrder(req, int(compID)))

	var stored Order
	require.NoError(t, db.Preload("OrderItems").First(&stored).Error)
	assert.Equal(t, "sale", stored.OrderType)
	assert.Equal(t, "Pending", stored.Status)
	assert.Equal(t, compID, stored.CompanyID)
	assert.Contains(t, stored.OrderNumber, "ORD-")
	require.Len(t, stored.OrderItems, 1)
	assert.Equal(t, 2, stored.OrderItems[0].Quantity)
}

func TestService_ReadOrder(t *testing.T) {
	svc, db := setupService(t)
	o := seedOrder(t, db, "sale", "Pending")

	got, err := svc.ReadOrder(itoa(o.ID), int(o.CompanyID))
	require.NoError(t, err)
	assert.Equal(t, "sale", got.OrderType)
	assert.Equal(t, "Pending", got.Status)
	assert.Contains(t, got.OrderNumber, "ORD-")
	assert.Equal(t, "USD", got.Currency)
	assert.Equal(t, "Acme", got.BusinessPartnerName)
}

func TestService_ReadOrder_InvalidID(t *testing.T) {
	svc, _ := setupService(t)
	_, err := svc.ReadOrder("not-a-number", 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid Order ID")
}

func TestService_CheckOrderExist_NotFound(t *testing.T) {
	svc, _ := setupService(t)
	_, _, err := svc.CheckOrderExist("99999", 1)
	require.Error(t, err)
}

func TestService_Approve_FromPending(t *testing.T) {
	svc, db := setupService(t)
	o := seedOrder(t, db, "sale", "Pending")

	require.NoError(t, svc.Approve(itoa(o.ID), int(o.CompanyID)))
	got, err := svc.ReadOrder(itoa(o.ID), int(o.CompanyID))
	require.NoError(t, err)
	assert.Equal(t, "Approved", got.Status)
}

func TestService_Approve_RejectsNonPending(t *testing.T) {
	svc, db := setupService(t)
	o := seedOrder(t, db, "sale", "Approved")

	err := svc.Approve(itoa(o.ID), int(o.CompanyID))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "can not changed to approve")
}

func TestService_Pack_FromApprovedSale(t *testing.T) {
	svc, db := setupService(t)
	o := seedOrder(t, db, "sale", "Approved")

	require.NoError(t, svc.Pack(itoa(o.ID), int(o.CompanyID)))
	got, err := svc.ReadOrder(itoa(o.ID), int(o.CompanyID))
	require.NoError(t, err)
	assert.Equal(t, "Packed", got.Status)
}

func TestService_Pack_RejectsWrongType(t *testing.T) {
	svc, db := setupService(t)
	o := seedOrder(t, db, "purchase", "Approved") // Pack requires a sale order

	err := svc.Pack(itoa(o.ID), int(o.CompanyID))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "can not changed to Packed")
}

func TestService_Ship_FromPackedSale(t *testing.T) {
	svc, db := setupService(t)
	o := seedOrder(t, db, "sale", "Packed")

	require.NoError(t, svc.Ship(itoa(o.ID), int(o.CompanyID)))
	got, err := svc.ReadOrder(itoa(o.ID), int(o.CompanyID))
	require.NoError(t, err)
	assert.Equal(t, "Shipped", got.Status)
}

func TestService_Ship_RejectsNonPacked(t *testing.T) {
	svc, db := setupService(t)
	o := seedOrder(t, db, "sale", "Approved") // not yet Packed

	err := svc.Ship(itoa(o.ID), int(o.CompanyID))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "can not changed to Shipped")
}

func TestService_Receive_FromWaitingPurchase(t *testing.T) {
	svc, db := setupService(t)
	o := seedOrder(t, db, "purchase", "Waiting")

	require.NoError(t, svc.Receive(itoa(o.ID), int(o.CompanyID)))
	got, err := svc.ReadOrder(itoa(o.ID), int(o.CompanyID))
	require.NoError(t, err)
	assert.Equal(t, "Received", got.Status)
}

func TestService_Receive_RejectsWrongType(t *testing.T) {
	svc, db := setupService(t)
	o := seedOrder(t, db, "sale", "Waiting") // Receive requires "purchase"

	err := svc.Receive(itoa(o.ID), int(o.CompanyID))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "can not changed to Received")
}

// seedProductWithStock creates a product and a stock row with the given
// available/reserved amounts, returning the product ID.
func seedProductWithStock(t *testing.T, db *gorm.DB, compID uint, avail, reserved int) uint {
	t.Helper()
	prod := products.Product{Name: "Widget", ProductNumber: "P" + itoa(uint(avail*100+reserved)), CompanyID: compID}
	require.NoError(t, db.Create(&prod).Error)
	require.NoError(t, db.Create(&products.Stock{
		ProductID: prod.ID, AvailableStock: avail, ReservedStock: reserved,
	}).Error)
	return prod.ID
}

// stockOf reads a product's current stock row.
func stockOf(t *testing.T, db *gorm.DB, productID uint) products.Stock {
	t.Helper()
	var st products.Stock
	require.NoError(t, db.Where("product_id = ?", productID).First(&st).Error)
	return st
}

// seedOrderWithItem creates and persists an order of the given type/status with
// a single line item of qty units of prodID. Refs (company/partner/currency)
// are passed in so callers can reuse them for product+stock seeding.
func seedOrderWithItem(t *testing.T, db *gorm.DB, compID, bpID, curID uint, orderType, status string, prodID uint, qty int) Order {
	t.Helper()
	o := makeOrder(compID, bpID, curID, orderType, status)
	o.OrderNumber = "ORD-" + itoa(uint(statusSeq()))
	o.OrderItems = []OrderItem{{ProductID: prodID, Quantity: qty, PerItemPrice: 10}}
	require.NoError(t, db.Create(&o).Error)
	return o
}

// --- Sale stock transitions ---

func TestService_Approve_SaleReservesStock(t *testing.T) {
	svc, db := setupService(t)
	compID, bpID, curID := seedRefs(t, db)
	prodID := seedProductWithStock(t, db, compID, 10, 0)
	o := seedOrderWithItem(t, db, compID, bpID, curID, "sale", "Pending", prodID, 4)

	require.NoError(t, svc.Approve(itoa(o.ID), int(o.CompanyID)))

	got, err := svc.ReadOrder(itoa(o.ID), int(o.CompanyID))
	require.NoError(t, err)
	assert.Equal(t, "Approved", got.Status)

	st := stockOf(t, db, prodID)
	assert.Equal(t, 6, st.AvailableStock) // 10 - 4
	assert.Equal(t, 4, st.ReservedStock)  // 0 + 4
}

func TestService_Approve_SaleInsufficientStockRollsBack(t *testing.T) {
	svc, db := setupService(t)
	compID, bpID, curID := seedRefs(t, db)
	prodID := seedProductWithStock(t, db, compID, 2, 0) // only 2 available
	o := seedOrderWithItem(t, db, compID, bpID, curID, "sale", "Pending", prodID, 5)

	err := svc.Approve(itoa(o.ID), int(o.CompanyID))
	require.Error(t, err)

	// Status stayed Pending and stock is untouched.
	got, err := svc.ReadOrder(itoa(o.ID), int(o.CompanyID))
	require.NoError(t, err)
	assert.Equal(t, "Pending", got.Status)
	st := stockOf(t, db, prodID)
	assert.Equal(t, 2, st.AvailableStock)
	assert.Equal(t, 0, st.ReservedStock)
}

func TestService_Ship_FulfillsReserved(t *testing.T) {
	svc, db := setupService(t)
	compID, bpID, curID := seedRefs(t, db)
	prodID := seedProductWithStock(t, db, compID, 0, 4) // already reserved
	o := seedOrderWithItem(t, db, compID, bpID, curID, "sale", "Packed", prodID, 4)

	require.NoError(t, svc.Ship(itoa(o.ID), int(o.CompanyID)))

	got, err := svc.ReadOrder(itoa(o.ID), int(o.CompanyID))
	require.NoError(t, err)
	assert.Equal(t, "Shipped", got.Status)
	st := stockOf(t, db, prodID)
	assert.Equal(t, 0, st.AvailableStock)
	assert.Equal(t, 0, st.ReservedStock) // 4 - 4
}

func TestService_Cancel_SaleApprovedReleasesReservation(t *testing.T) {
	svc, db := setupService(t)
	compID, bpID, curID := seedRefs(t, db)
	prodID := seedProductWithStock(t, db, compID, 6, 4) // 4 reserved
	o := seedOrderWithItem(t, db, compID, bpID, curID, "sale", "Approved", prodID, 4)

	require.NoError(t, svc.Cancel(itoa(o.ID), int(o.CompanyID)))

	got, err := svc.ReadOrder(itoa(o.ID), int(o.CompanyID))
	require.NoError(t, err)
	assert.Equal(t, "Canceled", got.Status)
	st := stockOf(t, db, prodID)
	assert.Equal(t, 10, st.AvailableStock) // 6 + 4 released
	assert.Equal(t, 0, st.ReservedStock)   // 4 - 4
}

func TestService_Cancel_SalePendingNoStockChange(t *testing.T) {
	svc, db := setupService(t)
	compID, bpID, curID := seedRefs(t, db)
	prodID := seedProductWithStock(t, db, compID, 10, 0) // nothing reserved yet
	o := seedOrderWithItem(t, db, compID, bpID, curID, "sale", "Pending", prodID, 4)

	require.NoError(t, svc.Cancel(itoa(o.ID), int(o.CompanyID)))

	got, err := svc.ReadOrder(itoa(o.ID), int(o.CompanyID))
	require.NoError(t, err)
	assert.Equal(t, "Canceled", got.Status)
	st := stockOf(t, db, prodID)
	assert.Equal(t, 10, st.AvailableStock) // unchanged
	assert.Equal(t, 0, st.ReservedStock)
}

func TestService_Cancel_AlreadyCanceled(t *testing.T) {
	svc, db := setupService(t)
	o := seedOrder(t, db, "sale", "Canceled")

	err := svc.Cancel(itoa(o.ID), int(o.CompanyID))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already canceled")
}

// --- Purchase stock transitions ---

func TestService_MarkWaiting_PurchaseFromApproved(t *testing.T) {
	svc, db := setupService(t)
	compID, bpID, curID := seedRefs(t, db)
	prodID := seedProductWithStock(t, db, compID, 0, 0)
	o := seedOrderWithItem(t, db, compID, bpID, curID, "purchase", "Approved", prodID, 5)

	require.NoError(t, svc.MarkWaiting(itoa(o.ID), int(o.CompanyID)))

	got, err := svc.ReadOrder(itoa(o.ID), int(o.CompanyID))
	require.NoError(t, err)
	assert.Equal(t, "Waiting", got.Status)
	// no stock change on waiting
	st := stockOf(t, db, prodID)
	assert.Equal(t, 0, st.AvailableStock)
	assert.Equal(t, 0, st.ReservedStock)
}

func TestService_Receive_PurchaseAddsStock(t *testing.T) {
	svc, db := setupService(t)
	compID, bpID, curID := seedRefs(t, db)
	prodID := seedProductWithStock(t, db, compID, 3, 0)
	o := seedOrderWithItem(t, db, compID, bpID, curID, "purchase", "Waiting", prodID, 5)

	require.NoError(t, svc.Receive(itoa(o.ID), int(o.CompanyID)))

	got, err := svc.ReadOrder(itoa(o.ID), int(o.CompanyID))
	require.NoError(t, err)
	assert.Equal(t, "Received", got.Status)
	st := stockOf(t, db, prodID)
	assert.Equal(t, 8, st.AvailableStock) // 3 + 5
	assert.Equal(t, 0, st.ReservedStock)
}

func TestService_Cancel_PurchaseNoStockChange(t *testing.T) {
	svc, db := setupService(t)
	compID, bpID, curID := seedRefs(t, db)
	prodID := seedProductWithStock(t, db, compID, 3, 0)
	o := seedOrderWithItem(t, db, compID, bpID, curID, "purchase", "Waiting", prodID, 5)

	require.NoError(t, svc.Cancel(itoa(o.ID), int(o.CompanyID)))

	got, err := svc.ReadOrder(itoa(o.ID), int(o.CompanyID))
	require.NoError(t, err)
	assert.Equal(t, "Canceled", got.Status)
	st := stockOf(t, db, prodID)
	assert.Equal(t, 3, st.AvailableStock) // unchanged
	assert.Equal(t, 0, st.ReservedStock)
}

func TestService_DeleteOrder(t *testing.T) {
	svc, db := setupService(t)
	o := seedOrder(t, db, "sale", "Pending")

	require.NoError(t, svc.DeleteOrder(o.ID, int(o.CompanyID)))
	_, err := svc.ReadOrder(itoa(o.ID), int(o.CompanyID))
	assert.Error(t, err)
}

func TestService_DeleteOrder_NotFound(t *testing.T) {
	svc, _ := setupService(t)
	err := svc.DeleteOrder(99999, 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Not found")
}

func TestService_AllOrders(t *testing.T) {
	svc, db := setupService(t)
	o := seedOrder(t, db, "sale", "Pending")

	got, err := svc.AllOrders(int(o.CompanyID))
	require.NoError(t, err)
	require.Len(t, got.Orders, 1)
	assert.Equal(t, "sale", got.Orders[0].OrderType)
	assert.Equal(t, "Pending", got.Orders[0].Status)
}

func TestService_UpdateOrder_NoChanges(t *testing.T) {
	svc, db := setupService(t)
	o := seedOrder(t, db, "sale", "Pending")

	// Same values as seeded -> no diff -> "no changes detected".
	req := &UpdateOrderRequest{
		ID:                o.ID,
		OrderType:         o.OrderType,
		BusinessPartnerID: o.BusinessPartnerID,
		CurrencyID:        o.CurrencyID,
		ExchangeRate:      o.ExchangeRate,
	}
	err := svc.UpdateOrder(req, 1, int(o.CompanyID))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no changes detected")
}

func TestService_UpdateOrder(t *testing.T) {
	svc, db := setupService(t)
	o := seedOrder(t, db, "sale", "Pending")

	// create a second currency to switch to
	cur2 := Currency{Name: "EUR", PersianName: "یورو"}
	require.NoError(t, db.Create(&cur2).Error)

	req := &UpdateOrderRequest{
		ID:                o.ID,
		OrderType:         o.OrderType,
		BusinessPartnerID: o.BusinessPartnerID,
		CurrencyID:        cur2.ID, // changed
		ExchangeRate:      o.ExchangeRate,
	}
	require.NoError(t, svc.UpdateOrder(req, 1, int(o.CompanyID)))

	var stored Order
	require.NoError(t, db.First(&stored, o.ID).Error)
	assert.Equal(t, cur2.ID, stored.CurrencyID)
}

func TestService_ModifiedFields(t *testing.T) {
	svc, db := setupService(t)
	o := seedOrder(t, db, "sale", "Pending")
	cur2 := Currency{Name: "EUR", PersianName: "یورو"}
	require.NoError(t, db.Create(&cur2).Error)

	changes := svc.modifiedFields(&UpdateOrderRequest{
		ID:                o.ID,
		OrderType:         "purchase", // changed
		CurrencyID:        cur2.ID,    // changed
		BusinessPartnerID: o.BusinessPartnerID,
		ExchangeRate:      o.ExchangeRate,
	}, int(o.CompanyID))
	assert.Contains(t, changes, "OrderType")
	assert.Contains(t, changes, "CurrencyID")
	assert.NotContains(t, changes, "ExchangeRate")
}

func TestService_ModifiedFields_NotFound(t *testing.T) {
	svc, _ := setupService(t)
	changes := svc.modifiedFields(&UpdateOrderRequest{ID: 99999}, 1)
	assert.Nil(t, changes)
}

// itoa is a tiny local helper to keep imports lean.
func itoa(n uint) string {
	if n == 0 {
		return "0"
	}
	var digits [20]byte
	i := len(digits)
	for n > 0 {
		i--
		digits[i] = byte('0' + n%10)
		n /= 10
	}
	return string(digits[i:])
}
