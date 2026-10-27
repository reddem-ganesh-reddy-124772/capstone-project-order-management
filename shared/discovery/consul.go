package discovery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"order-management/config"
)

// Client communicates with Consul to register and find services.
type Client struct {
	// baseURL stores the address of the Consul server.
	baseURL string

	// httpClient sends requests to Consul with a timeout.
	httpClient *http.Client
}

// serviceRegistration contains the details Consul needs to register a service.
type serviceRegistration struct {
	ID      string       `json:"ID"`
	Name    string       `json:"Name"`
	Address string       `json:"Address"`
	Port    int          `json:"Port"`
	Check   serviceCheck `json:"Check"`
}

// serviceCheck tells Consul how to check whether a service is healthy.
type serviceCheck struct {
	HTTP     string `json:"HTTP"`
	Interval string `json:"Interval"`
	Timeout  string `json:"Timeout"`
}

// discoveredService contains the address and port returned by Consul.
type discoveredService struct {
	Service struct {
		Address string `json:"Address"`
		Port    int    `json:"Port"`
	} `json:"Service"`
}

// New creates a client using the Consul settings from configuration.
func New(cfg config.Config) *Client {
	// Use HTTP when no Consul scheme is configured.
	scheme := strings.TrimSpace(cfg.ConsulScheme)
	if scheme == "" {
		scheme = "http"
	}

	// Use the local Consul address when none is configured.
	address := strings.TrimSpace(cfg.ConsulAddress)
	if address == "" {
		address = "127.0.0.1:8500"
	}

	// Store the Consul URL and set a timeout for requests.
	return &Client{
		baseURL: strings.TrimRight(scheme+"://"+address, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// RegisterService adds a running service to Consul.
// The healthURL must point to the service's health endpoint.
func (c *Client) RegisterService(
	serviceName string,
	serviceID string,
	address string,
	port int,
	healthURL string,
) error {
	// Prepare the service details for Consul.
	registration := serviceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: address,
		Port:    port,
		Check: serviceCheck{
			HTTP:     healthURL,
			Interval: "40s",
			Timeout:  "2s",
		},
	}

	// Convert the service details into JSON.
	body, err := json.Marshal(registration)
	if err != nil {
		return fmt.Errorf("prepare service registration: %w", err)
	}

	// Send the registration request to Consul.
	req, err := http.NewRequest(
		http.MethodPut,
		c.baseURL+"/v1/agent/service/register",
		strings.NewReader(string(body)),
	)
	if err != nil {
		return fmt.Errorf("create registration request: %w", err)
	}

	// Tell Consul that the request body contains JSON.
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("register service with Consul: %w", err)
	}
	defer resp.Body.Close()

	// Consul returns a successful status when registration succeeds.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Consul registration failed with status %d", resp.StatusCode)
	}

	return nil
}

// DiscoverService finds a healthy service in Consul and returns its URL.
func (c *Client) DiscoverService(serviceName string) (string, error) {
	// Ask Consul for healthy instances of the requested service.
	endpoint := c.baseURL + "/v1/health/service/" +
		url.PathEscape(serviceName) + "?passing=true"

	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("discover service %q: %w", serviceName, err)
	}
	defer resp.Body.Close()

	// Stop if Consul cannot provide the service information.
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Consul discovery failed with status %d", resp.StatusCode)
	}

	// Decode the list of healthy service instances.
	var services []discoveredService
	if err := json.NewDecoder(resp.Body).Decode(&services); err != nil {
		return "", fmt.Errorf("read service information: %w", err)
	}

	// Return an error if Consul has no healthy instance.
	if len(services) == 0 {
		return "", fmt.Errorf("no healthy instance found for service %q", serviceName)
	}

	// Use the first healthy instance returned by Consul.
	service := services[0].Service

	if service.Address == "" || service.Port <= 0 {
		return "", fmt.Errorf("service %q has an invalid address", serviceName)
	}

	// Build the URL that the calling service can use.
	return fmt.Sprintf("http://%s:%d", service.Address, service.Port), nil
}
