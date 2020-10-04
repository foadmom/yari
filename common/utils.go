package common

import (
	"math/rand"
	"net"
	"time"

	"github.com/google/uuid"
)

// ========================================================
// returns today's date (yyyy-mm-dd) as string
// ========================================================
func GetTodaysDateString () string {
	var _today  time.Time = time.Now ();
	var _date   string = _today.Format("2006-01-02");

	return _date;
}


func GetMyNetInterfaceDetails () (net.IP, string, net.HardwareAddr) {
	var _ipAddress  			net.IP;
	var _networkHardwareName 	string;
	var _macAddress				net.HardwareAddr;

	interfaces, _ := net.Interfaces()
    for _, _interf := range interfaces {
        if addrs, err := _interf.Addrs(); err == nil {
            for _, _address := range addrs {

				_ipnet, ok := _address.(*net.IPNet);
				if ( ok && !_ipnet.IP.IsLoopback() ) {
					if _ipnet.IP.To4() != nil {
						_ipAddress = _ipnet.IP;
						_networkHardwareName = _interf.Name;
						_macAddress = _interf.HardwareAddr;
					}
				}
            }
        }
    }
	return _ipAddress, _networkHardwareName, _macAddress;
}


func GenerateUuid_V4 () string {
	var _uuid 	uuid.UUID = uuid.New ();
	return _uuid.String();
}


func GenerateRandomInt (lower int, upper int) int {
	rand.Seed (time.Now().UnixNano());
	var _rand int = rand.Intn(upper - lower) + lower;
	return _rand;
}
	
func TimeFormat (this time.Time) string {
	return this.Format (time.RFC3339Nano);
}