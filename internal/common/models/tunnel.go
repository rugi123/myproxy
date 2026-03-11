package models

type TunnelStatus string

type Tunnel struct {
	LocalPort  int
	RemotePort int
	Status     TunnelStatus
	Protocol   string
}
