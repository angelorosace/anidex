package googleApp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

type AppCredentials struct {
	Type                string `json:"type"`
	ProjectId           string `json:"project_id"`
	PrivateKeyId        string `json:"private_key_id"`
	PrivateKey          string `json:"private_key"`
	ClientEmail         string `json:"client_email"`
	ClientId            string `json:"client_id"`
	AuthUri             string `json:"auth_uri"`
	TokenURI            string `json:"token_uri"`
	AuthProviderCertUrl string `json:"auth_provider_x509_cert_url"`
	ClientCertUrl       string `json:"client_x509_cert_url"`
	UniverseDomain      string `json:"universe_domain"`
}

func BuildApp(ctx context.Context) (*firebase.App, error) {

	opt, err := buildClientOptions()
	if err != nil {
		return &firebase.App{}, fmt.Errorf("Client Options could not be built: %v", err)
	}

	fmt.Sprintln("buildClientOptions")
	fmt.Sprintln(opt)

	conf := buildConfig()

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		return &firebase.App{}, fmt.Errorf("error initializing app: %v", err)
	}

	return app, nil
}

func buildConfig() *firebase.Config {
	port := os.Getenv("PORT")

	var conf *firebase.Config
	if port == "" {
		conf = &firebase.Config{
			DatabaseURL: "https://anidex-db184-default-rtdb.firebaseio.com/",
		}
	} else {
		conf = &firebase.Config{
			DatabaseURL: os.Getenv("DatabaseUrl"),
		}
	}

	return conf
}

func buildClientOptions() (option.ClientOption, error) {
	port := os.Getenv("PORT")

	var o option.ClientOption
	if port == "" {
		o = option.WithCredentialsFile("aindex_firebase_config.json")
	} else {
		appCredentials := buildAppCredentials()
		fmt.Println(appCredentials)
		p, err := json.Marshal(appCredentials)
		if err != nil {
			return nil, fmt.Errorf("App authentication failed: %v", err)
		}
		o = option.WithCredentialsJSON(p)
	}
	return o, nil
}

func buildAppCredentials() AppCredentials {
	return AppCredentials{
		Type:                "service_account",
		ProjectId:           os.Getenv("ProjectId"),
		PrivateKeyId:        os.Getenv("PrivateKeyId"),
		PrivateKey:          os.Getenv("GooglePrivateKey"),
		ClientEmail:         os.Getenv("GoogleClientEmail"),
		ClientId:            os.Getenv("GoogleClientId"),
		ClientCertUrl:       os.Getenv("ClientCertUrl"),
		AuthUri:             "https://accounts.google.com/o/oauth2/auth",
		TokenURI:            "https://oauth2.googleapis.com/token",
		AuthProviderCertUrl: "https://www.googleapis.com/oauth2/v1/certs",
		UniverseDomain:      "googleapis.com",
	}
}
