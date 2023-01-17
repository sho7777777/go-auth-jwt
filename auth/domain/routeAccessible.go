package domain

var routeAccessibleForUser = []string{"getResource"}
var routeAccessibleForAdmin = []string{"getResource", "adminPage"}

func GetRouteForUser() []string {
	return routeAccessibleForUser
}

func GetRouteForAdmin() []string {
	return routeAccessibleForAdmin
}
