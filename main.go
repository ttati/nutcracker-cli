package main

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"crypto/tls"
)

const (
	version  = "0.0.1"
	base64ID = "$base64$"
)

var (
	server string
	key    string
	id     string
)

func getSecret(c *cli.Context) {
	u := url.URL{
		Host:   server,
		Scheme: "https",
		Path:   "/secrets/view/" + c.String("name"),
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("X-Secret-ID", id)
	req.Header.Add("X-Secret-Key", key)

	client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
    }
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatal(resp.Status)
	}

	secret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Decode secret if base64 encoded
	if bytes.Compare(secret[0:8], []byte(base64ID)) == 0 {
		var decoded []byte
		_, err = base64.StdEncoding.Decode(decoded, secret)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("%s\n", decoded)
	} else {
		log.Printf("%s\n", secret)
	}
}

func listSecrets(c *cli.Context) {
	u := url.URL{
		Host:   server,
		Scheme: "https",
		Path:   "/secrets/list/secrets",
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("X-Secret-ID", id)
	req.Header.Add("X-Secret-Key", key)

	client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
    }
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatal(resp.Status)
	}

	for {
		data := make([]byte, 4<<10) // Read 4KB at a time
		_, err := resp.Body.Read(data)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s\n", data)
	}

	err = resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "nutcracker-cli"
	app.Usage = "CLI interface for nutcracker"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "server, s",
			Usage:       "Nutcracker server.  e.g localhost:443",
			Destination: &server,
		},
		cli.StringFlag{
			Name:        "id, i",
			Usage:       "Nutcracker API ID",
			Destination: &id,
		},
		cli.StringFlag{
			Name:        "key, k",
			Usage:       "Nutcracker API key",
			Destination: &key,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "get",
			Aliases: []string{"g"},
			Usage:   "get a secret",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Usage: "name of the secret",
				},
			},
			Action: getSecret,
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list all secrets",
			Action:  listSecrets,
		},
	}

	app.Run(os.Args)

}