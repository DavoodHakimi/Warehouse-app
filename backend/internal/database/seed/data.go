package seed

type UserType struct {
	ID          uint
	Name        string
	PersianName string
	Description string
}

var userTypes = []UserType{
	{Name: "CEO", PersianName: "مدیر", Description: "مدیر همه"},
	{Name: "Warehouse Manager", PersianName: "سرپرست انبار", Description: ""},
	{Name: "Storeman-Full", PersianName: "انبار دار", Description: ""},
	{Name: "Storeman-EnterOnly", PersianName: "انباردار- ورود", Description: ""},
	{Name: "Storeman-ExitOnly", PersianName: "انباردار- خروج", Description: ""},
}

var permissionTypes = map[string][]string{
	"users":    {"read", "create", "delete", "update"},
	"products": {"read", "create", "delete", "update"},
	"partners": {"read", "create", "delete", "update"},
	"orders":   {"read", "create", "delete", "update", "receive", "ship", "pack"},
}

var rolePermissions = map[int][]string{
	1: {
		"users.read", "users.create", "users.delete", "users.update",
		"products.read", "products.create", "products.delete", "products.update",
		"partners.read", "partners.create", "partners.delete", "partners.update",
		"orders.read", "orders.create", "orders.delete",
		"orders.update", "orders.receive", "orders.ship", "orders.pack",
	},
	2: {
		"products.read", "products.create", "products.delete", "products.update",
		"orders.read", "orders.create", "orders.delete",
		"orders.update", "orders.receive", "orders.ship", "orders.pack",
	},
	3: {
		"products.read",
		"orders.read", "orders.receive", "orders.ship", "orders.pack",
	},
	4: {
		"products.read",
		"orders.read", "orders.receive",
	},
	5: {
		"products.read",
		"orders.read", "orders.ship",
	},
}
