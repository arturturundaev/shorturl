package config

//  -a=http://localhost:8081/api/shorten -b=http://localhost:8081/api/shorten
import (
	"cmp"
	"fmt"
	"net/netip"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AddressStart AddressStartType
	BaseShort    BaseShortURLType
	FileStorage  FileStorageType
	DatabaseURL  DatabaseURLType
	StorageType  string
}

const StorageTypeMemory = "Memory"
const StorageTypeFile = "File"
const StorageTypeDB = "DB"

type FileStorageType struct {
	Path string
}

type AddressStartType struct {
	URL  string
	Port string
}

type BaseShortURLType struct {
	URL string
}

type DatabaseURLType struct {
	URL string
}

func NewConfig(ServerAddress, BaseURL, FileStorage, databaseURL string) *Config {
	fmt.Println("ServerAddress: " + ServerAddress + " ; BaseURL:" + BaseURL + " ; FileStorage:" + FileStorage + " ; databaseURL:" + databaseURL)

	var storageType = StorageTypeMemory

	if FileStorage != "" {
		storageType = StorageTypeFile
	}

	if databaseURL != "" {
		storageType = StorageTypeDB
	}

	var URL, port string
	data := strings.Split(cmp.Or(ServerAddress, os.Getenv("SERVER_ADDRESS"), "localhost:8080"), ":")
	URL = data[0]
	port = data[1]

	BaseURL = cmp.Or(BaseURL, os.Getenv("BASE_URL"), "http://localhost:8080")
	FileStorage = cmp.Or(FileStorage, os.Getenv("FILE_STORAGE_PATH"), "/tmp/db.txt")
	databaseURL = cmp.Or(databaseURL, os.Getenv("DATABASE_DSN"), "postgres://postgres:postgres@localhost:5432/shorturl?sslmode=disable")

	return &Config{
		AddressStart: AddressStartType{URL: URL, Port: port},
		BaseShort:    BaseShortURLType{URL: BaseURL},
		FileStorage:  FileStorageType{Path: FileStorage},
		DatabaseURL:  DatabaseURLType{URL: databaseURL},
		StorageType:  storageType,
	}
}

func (d *AddressStartType) String() string {
	arr := make([]string, 0)
	arr = append(arr, d.URL, d.Port)

	if arr[0] != "" && arr[1] != "" {
		return fmt.Sprint(strings.Join(arr, ":"))
	}

	return ""
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

func (d *FileStorageType) String() string {
	return d.Path
}

func (d *FileStorageType) Set(flagValue string) error {
	d.Path = flagValue

	return nil
}

func (d *DatabaseURLType) String() string {
	return d.URL
}

func (d *DatabaseURLType) Set(flagValue string) error {
	d.URL = flagValue

	return nil
}
