package config

//  -a=http://localhost:8081/api/shorten -b=http://localhost:8081/api/shorten
import (
	"bytes"
	"cmp"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"flag"
	"log"
	"math/big"
	"net"
	"os"
	"strconv"
	"time"
)

// Config структура
type Config struct {
	AddressStart string `json:"server_address"`
	BaseShort    string `json:"base_url"`
	FileStorage  string `json:"file_storage_path"`
	DatabaseURL  string `json:"database_dsn"`
	StorageType  string
	FullLog      bool
	HTTPS        struct {
		Enabled    bool `json:"enable_https"`
		SSLKeyPath string
		SSLPemPath string
	}
}

// StorageTypeMemory место зранения
const StorageTypeMemory = "Memory"

// StorageTypeFile место зранения
const StorageTypeFile = "File"

// StorageTypeDB место зранения
const StorageTypeDB = "DB"

// NewConfig получение конфигов
func NewConfig() *Config {
	var configPath string
	flag.StringVar(&configPath, "c", "", "Path to config file")

	cfg := &Config{}

	flag.StringVar(&cfg.AddressStart, "a", "", "start url and port")
	flag.StringVar(&cfg.BaseShort, "b", "", "url redirect")
	flag.StringVar(&cfg.FileStorage, "f", "", "file storage path")
	flag.StringVar(&cfg.DatabaseURL, "d", "", "database storage path")
	flag.BoolVar(&cfg.HTTPS.Enabled, "s", false, "ssl enabled")
	flag.StringVar(&cfg.HTTPS.SSLKeyPath, "sslk", "./auto_server.key", "Path to ssl key file")
	flag.StringVar(&cfg.HTTPS.SSLPemPath, "sslp", "./auto_server.pem", "Path to ssl pem file")
	flag.Parse()

	preConfig := &Config{}
	if configPath != "" {
		rawContent, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatal(err)
		}

		if err = json.Unmarshal(rawContent, preConfig); err != nil {
			log.Fatal(err)
		}
	}

	cfg.AddressStart = cmp.Or(cfg.AddressStart, os.Getenv("SERVER_ADDRESS"), preConfig.AddressStart, "127.0.0.1:8080")
	cfg.BaseShort = cmp.Or(cfg.BaseShort, os.Getenv("BASE_URL"), preConfig.BaseShort, "127.0.0.1:8080")
	cfg.FileStorage = cmp.Or(cfg.FileStorage, os.Getenv("FILE_STORAGE_PATH"), preConfig.FileStorage)
	cfg.DatabaseURL = cmp.Or(cfg.DatabaseURL, os.Getenv("DATABASE_DSN"), preConfig.DatabaseURL)

	var storageType = StorageTypeMemory

	if cfg.FileStorage != "" || os.Getenv("FILE_STORAGE_PATH") != "" {
		storageType = StorageTypeFile
	}

	if cfg.DatabaseURL != "" || os.Getenv("DATABASE_DSN") != "" {
		storageType = StorageTypeDB
	}

	cfg.StorageType = storageType
	cfg.FullLog = true

	if envHTTPSStr := os.Getenv("ENABLE_HTTPS"); envHTTPSStr != "" {
		envHTTPS, err := strconv.ParseBool(envHTTPSStr)
		if err == nil {
			cfg.HTTPS.Enabled = envHTTPS
		}
	}

	if cfg.HTTPS.Enabled && cfg.HTTPS.SSLKeyPath == "./auto_server.key" {
		CreateTLSCert("./auto_server.pem", "./auto_server.key")
	}

	return cfg
}

// CreateTLSCert - generate TLS certificate and key for run server HTTPS
func CreateTLSCert(certPath string, keyPath string) error {
	cert := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"artur.turundaev"},
			Country:      []string{"RU"},
			Province:     []string{"Moscow"},
			Locality:     []string{"Moscow"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}

	privateKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	certBytes, _ := x509.CreateCertificate(rand.Reader, &cert, &cert, &privateKey.PublicKey, privateKey)
	err := saveToFile(certPath, "CERTIFICATE", certBytes)
	if err != nil {
		return err
	}

	err = saveToFile(keyPath, "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(privateKey))
	if err != nil {
		return err
	}

	return nil
}

func saveToFile(filePath string, cypherType string, cypher []byte) error {
	var (
		buf  bytes.Buffer
		file *os.File
	)

	_ = pem.Encode(&buf, &pem.Block{
		Type:  cypherType,
		Bytes: cypher,
	})

	file, _ = os.Create(filePath)
	defer file.Close()

	_, err := buf.WriteTo(file)
	if err != nil {
		return err
	}

	return nil
}
