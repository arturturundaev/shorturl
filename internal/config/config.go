package config

import (
	"fmt"
	"net/netip"
	"strconv"
	"strings"
)

type Config struct {
	AddressStart AddressStartType
	BaseShort    BaseShortURLType
}

type AddressStartType struct {
	URL  string
	Port string
}

type BaseShortURLType struct {
	URL string
}

func NewConfig() *Config {
	return &Config{
		AddressStart: AddressStartType{URL: "localhost", Port: "8080"},
		BaseShort:    BaseShortURLType{URL: "http://localhost:8080"},
	}
}

func (d *AddressStartType) String() string {
	arr := make([]string, 0)
	arr = append(arr, d.URL, d.Port)

	return fmt.Sprint(strings.Join(arr, ":"))
}

func (d *AddressStartType) Set(flagValue string) error {
	data := strings.Split(flagValue, ":")

	var ip string

	if data[0] == "localhost" {
		ip = data[0]
	} else {
		add, err := netip.ParseAddr(data[0])
		if err != nil {
			return err
		}

		ip = add.String()
	}

	port, err2 := strconv.Atoi(data[1])
	if err2 != nil {
		return err2
	}

	if port < 1 || port > 65535 {
		return fmt.Errorf("PORT incorrected")
	}

	d.URL = ip
	d.Port = data[1]

	return nil
}

func (d *BaseShortURLType) String() string {
	return d.URL
}

func (d *BaseShortURLType) Set(flagValue string) error {
	d.URL = flagValue

	return nil
}
