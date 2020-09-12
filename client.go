package main

// Mostly copied and modified from https://github.com/vmware/govmomi/tree/master/examples

import (
	"context"
	"errors"
	"net/url"
	"os"
	"strings"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/soap"
)

// getEnvString returns string from environment variable.
func getEnvString(v string, def string) string {
	r := os.Getenv(v)
	if r == "" {
		return def
	}

	return r
}

// getEnvBool returns boolean from environment variable.
func getEnvBool(v string, def bool) bool {
	r := os.Getenv(v)
	if r == "" {
		return def
	}

	switch strings.ToLower(r[0:1]) {
	case "t", "y", "1":
		return true
	}

	return false
}

const (
	envURL      = "GOVC_URL"
	envUserName = "GOVC_USERNAME"
	envPassword = "GOVC_PASSWORD"
	envInsecure = "GOVC_INSECURE"
)

func processOverride(u *url.URL) {
	envUsername := os.Getenv(envUserName)
	envPassword := os.Getenv(envPassword)

	// Override username if provided
	if envUsername != "" {
		var password string
		var ok bool

		if u.User != nil {
			password, ok = u.User.Password()
		}

		if ok {
			u.User = url.UserPassword(envUsername, password)
		} else {
			u.User = url.User(envUsername)
		}
	}

	// Override password if provided
	if envPassword != "" {
		var username string

		if u.User != nil {
			username = u.User.Username()
		}

		u.User = url.UserPassword(username, envPassword)
	}
}

// NewClient creates a govmomi.Client for use in the examples
func NewClient(ctx context.Context) (*govmomi.Client, error) {
	// Parse URL from string
	u, err := soap.ParseURL(getEnvString(envURL, ""))
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("No URL specified")
	}

	// Override username and/or password as required
	processOverride(u)

	insecure := getEnvBool(envInsecure, false)

	// Connect and log in to ESX or vCenter
	return govmomi.NewClient(ctx, u, insecure)
}
