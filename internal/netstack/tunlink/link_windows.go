package tunlink

import (
	"fmt"
	"net"
)

func New() (int, string, error) {
	return 0, "", fmt.Errorf("not implmeneted for this architecture")
}

func Remove(ID int) error {
	return fmt.Errorf("not implemented for this architecture")
}

func AddRoute(ID int, CIDR *net.IPNet) error {
	return fmt.Errorf("not implemented for this architecture")
}

func RemoveAllRoutes(ID int) error {
	return fmt.Errorf("not implemented for this architecture")
}

func GetName(ID int) (string, error) {
	return "", fmt.Errorf("not implemented for this architecture")
}
