//go:build ignore

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

func main() {
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("Failed to generate private key: %v\n", err)
		os.Exit(1)
	}

	// 创建证书模板
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Trashcan System"},
			Country:       []string{"CN"},
			Province:      []string{""},
			Locality:      []string{""},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour), // 1年有效期
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:     []string{"localhost"},
	}

	// 创建证书
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		fmt.Printf("Failed to create certificate: %v\n", err)
		os.Exit(1)
	}

	// 保存证书文件
	certFile, err := os.Create("cert.pem")
	if err != nil {
		fmt.Printf("Failed to create cert.pem: %v\n", err)
		os.Exit(1)
	}
	defer certFile.Close()

	err = pem.Encode(certFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})
	if err != nil {
		fmt.Printf("Failed to encode certificate: %v\n", err)
		os.Exit(1)
	}

	// 保存私钥文件
	keyFile, err := os.Create("key.pem")
	if err != nil {
		fmt.Printf("Failed to create key.pem: %v\n", err)
		os.Exit(1)
	}
	defer keyFile.Close()

	privateKeyDER := x509.MarshalPKCS1PrivateKey(privateKey)
	err = pem.Encode(keyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyDER,
	})
	if err != nil {
		fmt.Printf("Failed to encode private key: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ SSL证书生成成功！")
	fmt.Println("  - cert.pem: SSL证书文件")
	fmt.Println("  - key.pem: SSL私钥文件")
	fmt.Println("\n注意：这是自签名证书，浏览器会显示安全警告。")
	fmt.Println("生产环境请使用由CA签发的正式证书。")
}

