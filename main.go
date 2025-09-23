package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"time"

	"github.com/oracle/oci-go-sdk/common"
)

const name = "oci-cluster-token"

const version = "0.0.2"

var revision = "HEAD"

var outputTemplate = `{
    "apiVersion": "client.authentication.k8s.io/v1beta1",
    "kind": "ExecCredential",
    "status": {
        "token": "%s",
        "expirationTimestamp": "%s"
    }
}`

func ociConfigFile() string {
	var home string
	current, e := user.Current()
	if e != nil {
		//Give up and try to return something sensible
		home = os.Getenv("HOME")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
	} else {
		home = current.HomeDir
	}
	return filepath.Join(home, ".oci", "config")
}

func main() {
	var clusterID string
	var region string
	var showVersion bool
	flag.StringVar(&clusterID, "cluster-id", os.Getenv("OCI_CLUSTER"), "OCI Container Engine for Kubernetes Cluster OCID")
	flag.StringVar(&region, "region", os.Getenv("OCI_REGION"), "OCI Region (e.g., us-ashburn-1)")
	flag.BoolVar(&showVersion, "V", false, "Print the version")
	flag.Parse()

	if showVersion {
		fmt.Printf("%s %s (rev: %s/%s)\n", name, version, revision, runtime.Version())
		return
	}

	if clusterID == "" {
		log.Fatal("cluster-id is required")
	}

	defaultFileProvider, _ := common.ConfigurationProviderFromFileWithProfile(ociConfigFile(), "DEFAULT", "")
	environmentProvider := common.NewRawConfigurationProvider(
		os.Getenv("OCI_TENANCY"),
		os.Getenv("OCI_USER"),
		os.Getenv("OCI_REGION"),
		os.Getenv("OCI_FINGERPRINT"),
		os.Getenv("OCI_KEY_FILE"),
		common.String((os.Getenv("OCI_PASSPHRASE"))))

	provider, err := common.ComposingConfigurationProvider([]common.ConfigurationProvider{
		defaultFileProvider,
		environmentProvider,
	})
	if err != nil {
		log.Fatal(err)
	}

	if region == "" {
		if tmp, err := provider.Region(); err == nil {
			region = tmp
		}
	}
	_url := fmt.Sprintf("https://containerengine.%s.oraclecloud.com/cluster_request/%s", region, clusterID)
	req, err := http.NewRequest("GET", _url, nil)
	if err != nil {
		log.Fatal(err)
	}

	date := time.Now().UTC().Format(http.TimeFormat)
	req.Header.Set("Date", date)

	signer := common.DefaultRequestSigner(provider)
	err = signer.Sign(req)
	if err != nil {
		log.Fatal(err)
	}

	params := url.Values{}
	params.Add("date", req.Header.Get("Date"))
	params.Add("authorization", req.Header.Get("Authorization"))

	fmt.Printf(outputTemplate,
		base64.StdEncoding.EncodeToString([]byte(req.URL.String()+"?"+params.Encode())),
		time.Now().UTC().Add(4*time.Minute).Format("2006-01-02T15:04:05Z"),
	)
}
