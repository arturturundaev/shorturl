package config

//  -a=http://localhost:8081/api/shorten -b=http://localhost:8081/api/shorten
import (
	"cmp"
	"flag"
	"fmt"
	"net/netip"
	"os"
	"strconv"
	"strings"
)

// Config структура
type Config struct {
	AddressStart AddressStartType
	BaseShort    BaseShortURLType
	FileStorage  FileStorageType
	DatabaseURL  DatabaseURLType
	StorageType  string
	FullLog      bool
}

// StorageTypeMemory место зранения
const StorageTypeMemory = "Memory"

// StorageTypeFile место зранения
const StorageTypeFile = "File"

// StorageTypeDB место зранения
const StorageTypeDB = "DB"

// FileStorageType стуктура
type FileStorageType struct {
	Path string
}

// AddressStartType стуктура
type AddressStartType struct {
	URL  string
	Port string
}

// BaseShortURLType стуктура
type BaseShortURLType struct {
	URL string
}

// DatabaseURLType стуктура
type DatabaseURLType struct {
	URL string
}

// NewConfig получение конфигов
func NewConfig() *Config {
	var ServerAddress AddressStartType
	var BaseURL BaseShortURLType
	var FileStorage FileStorageType
	var databaseURL DatabaseURLType

	flag.Var(&ServerAddress, "a", "start url and port")
	flag.Var(&BaseURL, "b", "url redirect")
	flag.Var(&FileStorage, "f", "file storage path")
	flag.Var(&databaseURL, "d", "database storage path")
	flag.Parse()
	var URL, port string
	data := strings.Split(cmp.Or(ServerAddress.String(), os.Getenv("SERVER_ADDRESS"), "localhost:8080"), ":")
	URL = data[0]
	port = data[1]

	BaseURLFinal := cmp.Or(BaseURL.String(), os.Getenv("BASE_URL"), "http://localhost:8080")
	FileStorageFinal := cmp.Or(FileStorage.String(), os.Getenv("FILE_STORAGE_PATH"), "/tmp/db.txt")
	databaseURLFinal := cmp.Or(databaseURL.String(), os.Getenv("DATABASE_DSN"), "postgres://postgres:postgres@localhost:5432/shorturl?sslmode=disable")

	var storageType = StorageTypeMemory

	if FileStorage.String() != "" || os.Getenv("FILE_STORAGE_PATH") != "" {
		storageType = StorageTypeFile
	}

	if databaseURL.String() != "" || os.Getenv("DATABASE_DSN") != "" {
		storageType = StorageTypeDB
	}

	return &Config{
		AddressStart: AddressStartType{URL: URL, Port: port},
		BaseShort:    BaseShortURLType{URL: BaseURLFinal},
		FileStorage:  FileStorageType{Path: FileStorageFinal},
		DatabaseURL:  DatabaseURLType{URL: databaseURLFinal},
		StorageType:  storageType,
		FullLog:      true,
	}
}

// String AddressStartType
func (d *AddressStartType) String() string {
	arr := make([]string, 0)
	arr = append(arr, d.URL, d.Port)

	if arr[0] != "" && arr[1] != "" {
		return fmt.Sprint(strings.Join(arr, ":"))
	}

	return ""
}

// Set AddressStartType
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

// String BaseShortURLType
func (d *BaseShortURLType) String() string {
	return d.URL
}

// Set BaseShortURLType
func (d *BaseShortURLType) Set(flagValue string) error {
	d.URL = flagValue

	return nil
}

// String FileStorageType
func (d *FileStorageType) String() string {
	return d.Path
}

// Set FileStorageType
func (d *FileStorageType) Set(flagValue string) error {
	d.Path = flagValue

	return nil
}

// String DatabaseURLType
func (d *DatabaseURLType) String() string {
	return d.URL
}

// Set DatabaseURLType
func (d *DatabaseURLType) Set(flagValue string) error {
	d.URL = flagValue

	return nil
}
