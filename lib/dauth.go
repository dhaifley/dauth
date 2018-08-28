// Package lib provides data structures and functions for interacting with
// the authentication services.
package lib

import "github.com/dhaifley/dlib"

// ServiceInfo provides information about this service.
var ServiceInfo dlib.ServiceInfo

func init() {
	ServiceInfo = dlib.ServiceInfo{
		Name:    "dauth",
		Short:   "Authentication services",
		Long:    `Provides authentication and authorization services.`,
		Version: "0.1.1",
	}
}
