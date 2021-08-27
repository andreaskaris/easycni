package main

import (
	"fmt"
	"github.com/andreaskaris/easycni/cni"
	"strings"
)

func getIpAddress(subnet string) string {
	return strings.Replace(subnet, ".0/24", ".10/24", -1)
}

func getGateway(subnet string) string {
	return strings.Replace(subnet, ".0/24", ".1", -1)
}

func getMac() string {
	return "aa:aa:aa:aa:aa:aa"
}

func getInterfaceName() string {
	return "veth0"
}

func main() {
	c, err := cni.NewCni()
	if err != nil {
		cni.PrintError(err.Error())
	}

	switch cmd := c.GetCommand(); cmd {
	case "VERSION":
		fmt.Println(c.PrintVersion())
		return
	case "DEL":
		// do nothing
		fmt.Println("")
		return
	case "ADD":
		parameters, err := c.GetPluginParameters()
		if err != nil {
			fmt.Println(cni.PrintError(err.Error()))
			return
		}
		subnet := parameters["subnet"].(string)
		ipAddress := getIpAddress(subnet)
		gateway := getGateway(subnet)
		mac := getMac()
		interfaceName := getInterfaceName()
		netns := c.GetNetns()

		intf := cni.CniInterface{
			Name:    interfaceName,
			Mac:     mac,
			Sandbox: netns,
		}
		ip := cni.CniIp{
			Version:   "4",
			Address:   ipAddress,
			Gateway:   gateway,
			Interface: 0,
		}
		c.AddInterface(intf)
		c.AddIp(ip)
		fmt.Println(c.PrintSuccess())
	default:
		s := cni.PrintError("Not Implemented: " + cmd)
		fmt.Println(s)
	}
}
