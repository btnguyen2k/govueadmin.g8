/*
Package gvabe provides backend API for GoVueAdmin Frontend.

@author Thanh Nguyen <btnguyen2k@gmail.com>
@since template-v0.1.0
*/
package gvabe

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/btnguyen2k/goyai"
	"main/src/goapi"
	blogv2 "main/src/gvabe/bov2/blog"
	userv2 "main/src/gvabe/bov2/user"
)

var (
	userDaov2        userv2.UserDao
	blogPostDaov2    blogv2.BlogPostDao
	blogCommentDaov2 blogv2.BlogCommentDao
	blogVoteDaov2    blogv2.BlogVoteDao
)

// MyBootstrapper implements goapi.IBootstrapper
type MyBootstrapper struct {
	name string
}

var Bootstrapper = &MyBootstrapper{name: "gvabe"}

// Bootstrap implements goapi.IBootstrapper.Bootstrap
//
// Bootstrapper usually does the following:
// - register api-handlers with the global ApiRouter
// - other initializing work (e.g. creating DAO, initializing database, etc)
func (b *MyBootstrapper) Bootstrap() error {
	DEBUG = os.Getenv("DEBUG") != ""
	go routineUpdateSystemInfo()

	initRsaKeys()
	initI18n()
	initExter()
	initDaos()
	initApiHandlers(goapi.ApiRouter)
	initApiFilters(goapi.ApiRouter)
	return nil
}

// available since template-v0.2.0
func initExter() {
	if exterAppId = goapi.AppConfig.GetString("gvabe.exter.app_id"); exterAppId == "" {
		log.Printf("[WARN] No Exter app-id configured at [gvabe.exter.app_id], Exter login is disabled.")
	} else if exterBaseUrl = goapi.AppConfig.GetString("gvabe.exter.base_url"); exterBaseUrl == "" {
		log.Printf("[WARN] No Exter base-url configured at [gvabe.exter.base_url], default value will be used.")
		exterBaseUrl = "https://exteross.gpvcloud.com"
	}
	exterBaseUrl = strings.TrimSuffix(exterBaseUrl, "/") // trim trailing slashes
	if exterAppId != "" {
		exterClient = NewExterClient(exterAppId, exterBaseUrl)
	}
	log.Printf("[INFO] Exter app-id: %s / Base Url: %s", exterAppId, exterBaseUrl)

	go goFetchExterInfo(60)
}

// available since template-v0.4.0
func initI18n() {
	i18nConfigFileOrDir := goapi.AppConfig.GetString("gvabe.i18n.i18n_file_or_directory")
	if i18nConfigFileOrDir == "" {
		log.Println("[INFO] No i18n config specified at [gvabe.i18n.i18n_file_or_directory].")
		i18n = goyai.NullI18n()
		return
	}
	_, err := os.Stat(i18nConfigFileOrDir)
	if err != nil {
		panic(fmt.Sprintf("error loading i18n config [%s]: %e", i18nConfigFileOrDir, err))
	}

	defaultLocale := goapi.AppConfig.GetString("gvabe.i18n.default_locale")
	i18n, err = goyai.BuildI18n(goyai.I18nOptions{
		ConfigFileOrDir: i18nConfigFileOrDir,
		DefaultLocale:   defaultLocale,
		I18nFileFormat:  goyai.Auto,
	})
	if err != nil {
		panic(fmt.Sprintf("error loading i18n config [%s]: %e", i18nConfigFileOrDir, err))
	}

	locales := i18n.AvailableLocales()
	if len(locales) == 0 {
		log.Printf("[WARN] i18n config loaded from file [%s] but no locale configuration found", i18nConfigFileOrDir)
		return
	}

	if DEBUG {
		if defaultLocale == "" {
			defaultLocale = locales[0].Id
		}
		log.Printf("[DEBUG] i18n config loaded from [%s], available locales: %v / default: [%s]", i18nConfigFileOrDir, locales, defaultLocale)
	}
}

// available since template-v0.2.0
func initRsaKeys() {
	rsaPrivKeyFile := goapi.AppConfig.GetString("gvabe.keys.rsa_privkey_file")
	if rsaPrivKeyFile == "" {
		log.Println("[WARN] No RSA private key file configured at [gvabe.keys.rsa_privkey_file], generating one...")
		privKey, err := genRsaKey(2048)
		if err != nil {
			panic(err)
		}
		rsaPrivKey = privKey
	} else {
		log.Println(fmt.Sprintf("[INFO] Loading RSA private key from [%s]...", rsaPrivKeyFile))
		content, err := ioutil.ReadFile(rsaPrivKeyFile)
		if err != nil {
			panic(err)
		}
		block, _ := pem.Decode(content)
		if block == nil {
			panic(fmt.Sprintf("cannot decode PEM from file [%s]", rsaPrivKeyFile))
		}
		var der []byte
		passphrase := goapi.AppConfig.GetString("gvabe.keys.rsa_privkey_passphrase")
		if passphrase != "" {
			log.Println("[INFO] RSA private key is pass-phrase protected")
			if decrypted, err := x509.DecryptPEMBlock(block, []byte(passphrase)); err != nil {
				panic(err)
			} else {
				der = decrypted
			}
		} else {
			der = block.Bytes
		}
		if block.Type == "RSA PRIVATE KEY" {
			if privKey, err := x509.ParsePKCS1PrivateKey(der); err != nil {
				panic(err)
			} else {
				rsaPrivKey = privKey
			}
		} else if block.Type == "PRIVATE KEY" {
			if privKey, err := x509.ParsePKCS8PrivateKey(der); err != nil {
				panic(err)
			} else {
				rsaPrivKey = privKey.(*rsa.PrivateKey)
			}
		}
	}

	rsaPubKey = &rsaPrivKey.PublicKey

	if DEBUG {
		log.Printf("[DEBUG] Exter public key: {Size: %d / Exponent: %d / Modulus: %x}",
			rsaPubKey.Size()*8, rsaPubKey.E, rsaPubKey.N)

		pubBlockPKCS1 := pem.Block{
			Type:    "RSA PUBLIC KEY",
			Headers: nil,
			Bytes:   x509.MarshalPKCS1PublicKey(rsaPubKey),
		}
		rsaPubKeyPemPKCS1 := pem.EncodeToMemory(&pubBlockPKCS1)
		log.Printf("[DEBUG] Exter public key (PKCS1): %s", string(rsaPubKeyPemPKCS1))

		pubPKIX, _ := x509.MarshalPKIXPublicKey(rsaPubKey)
		pubBlockPKIX := pem.Block{
			Type:    "PUBLIC KEY",
			Headers: nil,
			Bytes:   pubPKIX,
		}
		rsaPubKeyPemPKIX := pem.EncodeToMemory(&pubBlockPKIX)
		log.Printf("[DEBUG] Exter public key (PKIX): %s", string(rsaPubKeyPemPKIX))
	}
}
