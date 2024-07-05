package config

import (
	"fmt"
	"net/netip"
	"strconv"
	"strings"
)

type Config struct {
	AddressStart AddressStartType
	BaseShort    BaseShortUrlType
}

type AddressStartType struct {
	Url  string
	Port string
}

type BaseShortUrlType struct {
	Url string
}

func NewConfig() *Config {
	return &Config{
		AddressStart: AddressStartType{Url: "localhost", Port: "8080"},
		BaseShort:    BaseShortUrlType{Url: "http://localhost:8080"},
	}
}

func (d *AddressStartType) String() string {
	arr := make([]string, 0)
	arr = append(arr, d.Url, d.Port)

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

	d.Url = ip
	d.Port = data[1]

	return nil
}

func (d *BaseShortUrlType) String() string {
	return d.Url
}

func (d *BaseShortUrlType) Set(flagValue string) error {
	d.Url = flagValue

	return nil
}
