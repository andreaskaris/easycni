package cni

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// based on https://github.com/containernetworking/cni/blob/spec-v0.3.1/SPEC.md

const (
	cniVersion = "0.3.1"
)

var cniSupportedVersions = []string{
	"0.1.0",
	"0.2.0",
	"0.3.0",
	"0.3.1",
	"0.4.0",
	"1.0.0",
}

type CniInput map[string]string

type Cni struct {
	input   CniInput
	version CniVersion
	success CniSuccess
}

type CniError struct {
	CniVersion string `json:"cniVersion"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	Details    string `json:"details"`
}

type CniVersion struct {
	CniVersion        string   `json:"cniVersion"`
	SupportedVersions []string `json:"supportedVersions"`
}

type CniInterface struct {
	Name    string `json:"name"`
	Mac     string `json:"mac"`
	Sandbox string `json:"sandbox"`
}

type CniIp struct {
	Version   string `json:"version"`
	Address   string `json:"address"`
	Gateway   string `json:"gateway"`
	Interface uint   `json:interface"`
}

type CniRoute struct {
	Dst string `json:"dst"`
	Gw  string `json:"gw"`
}

type CniDns struct {
	Nameservers []string `json:"nameservers"`
	Domain      string   `json:"domain"`
	Search      []string `json:"search"`
	Options     []string `json:"options"`
}

type CniSuccess struct {
	CniVersion string         `json:"cniVersion"`
	Interfaces []CniInterface `json:"interfaces"`
	Ips        []CniIp        `json:"ips"`
	Routes     []CniRoute     `json:"routes"`
	Dns        []CniDns       `json:"dns"`
}

func (c *Cni) AddInterface(i CniInterface) {
	c.success.Interfaces = append(c.success.Interfaces, i)
}

func (c *Cni) AddIp(i CniIp) {
	c.success.Ips = append(c.success.Ips, i)
}

func (c *Cni) AddRoute(r CniRoute) {
	c.success.Routes = append(c.success.Routes, r)
}

func (c *Cni) AddDns(d CniDns) {
	c.success.Dns = append(c.success.Dns, d)
}

func NewCni() (*Cni, error) {
	cni := &Cni{}
	var err error
	cni.input, err = readCniInput()
	if err != nil {
		return nil, err
	}
	cni.version = CniVersion{
		CniVersion:        cniVersion,
		SupportedVersions: cniSupportedVersions,
	}
	cni.success = CniSuccess{
		CniVersion: cniVersion,
	}

	return cni, nil
}

func (c *Cni) PrintVersion() string {
	s, err := json.Marshal(c.version)
	if err != nil {
		return PrintError(err.Error())
	}
	return string(s)
}

func (c *Cni) PrintSuccess() string {
	s, err := json.Marshal(c.success)
	if err != nil {
		return PrintError(err.Error())
	}
	return string(s)
}

func PrintError(msg string) string {
	cniErr := CniError{
		CniVersion: cniVersion,
		Code:       7,
		Msg:        msg,
		Details:    msg,
	}
	s, err := json.Marshal(cniErr)
	if err != nil {
		panic(err)
	}
	return string(s)
}

func (c *Cni) GetPath() string {
	return c.input["CNI_PATH"]
}

func (c *Cni) GetArgs() string {
	return c.input["CNI_ARGS"]
}

func (c *Cni) GetContainerId() string {
	return c.input["CNI_CONTAINERID"]
}

func (c *Cni) GetIfName() string {
	return c.input["CNI_IFNAME"]
}

func (c *Cni) GetCommand() string {
	return c.input["CNI_COMMAND"]
}

func (c *Cni) GetNetns() string {
	return c.input["CNI_NETNS"]
}

func (c *Cni) GetRawInput() map[string]string {
	return c.input
}

func (c *Cni) GetStdin() string {
	return c.input["STDIN"]
}

func (c *Cni) GetPluginParameters() (map[string]interface{}, error) {
	var result map[string]interface{}

	s := c.GetStdin()
	if s == "" {
		return nil, errors.New("Plugin parameters are empty. No STDIN")
	}

	err := json.Unmarshal([]byte(c.GetStdin()), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func readCniInput() (CniInput, error) {
	var err error

	cniInput := CniInput{}
	cniInput["CNI_COMMAND"] = os.Getenv("CNI_COMMAND")
	cniInput["CNI_CONTAINERID"] = os.Getenv("CNI_CONTAINERID")
	cniInput["CNI_NETNS"] = os.Getenv("CNI_NETNS")
	cniInput["CNI_IFNAME"] = os.Getenv("CNI_IFNAME")
	cniInput["CNI_ARGS"] = os.Getenv("CNI_ARGS")
	cniInput["CNI_PATH"] = os.Getenv("CNI_PATH")
	cniInput["STDIN"], err = readStdin()
	if err != nil {
		return nil, err
	}

	return cniInput, nil
}

func readStdin() (string, error) {
	var stdin string

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Mode()&os.ModeNamedPipe != 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			stdin = stdin + scanner.Text()
		}

		if scanner.Err() != nil {
			return "", errors.New("Cannot read from stdin")
		}
	}
	return stdin, nil
}

func Error(errMsg string) {
	fmt.Println(errMsg)
	os.Exit(1)
}
