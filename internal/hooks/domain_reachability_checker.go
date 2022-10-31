package hooks

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/eleven-sh/agent/proto"
	"github.com/eleven-sh/cli/internal/agent"
	cliEntities "github.com/eleven-sh/cli/internal/entities"
	"github.com/eleven-sh/eleven/entities"
	"github.com/google/uuid"
)

type DomainReachabilityChecker struct {
	agentClientBuilder agent.ClientBuilder
}

func NewDomainReachabilityChecker(
	agentClientBuilder agent.ClientBuilder,
) DomainReachabilityChecker {

	return DomainReachabilityChecker{
		agentClientBuilder: agentClientBuilder,
	}
}

func (d DomainReachabilityChecker) Check(
	env *entities.Env,
	domain string,
) (reachable bool, redirToHTTPS bool, returnedError error) {

	agentClient := d.agentClientBuilder.Build(
		agent.NewDefaultClientConfig(
			[]byte(env.SSHKeyPairPEMContent),
			env.InstancePublicIPAddress,
		),
	)

	uniqueID := uuid.NewString()
	servedPorts := cliEntities.BuildProtoEnvServedPortsFromEnv(env)

	err := agentClient.CheckDomainReachability(&proto.CheckDomainReachabilityRequest{
		Domain:      domain,
		ServedPorts: servedPorts,
		UniqueId:    uniqueID,
	}, func(stream agent.CheckDomainReachabilityStream) error {
		_, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		return err
	})

	if err != nil {
		returnedError = err
		return
	}

	pollTimeoutChan := time.After(20 * time.Second)
	pollSleepDuration := 1 * time.Second

	for {
		select {
		case <-pollTimeoutChan:
			return
		default:
			redirectedToHTTPS := false

			httpClient := &http.Client{
				Timeout: 4 * time.Second,
				CheckRedirect: func(nextReq *http.Request, prevReq []*http.Request) error {

					if nextReq.URL.Scheme == "https" &&
						nextReq.URL.Hostname() == domain {

						redirectedToHTTPS = true
					}

					return nil
				},
			}

			httpResp, err := httpClient.Get("http://" + domain)

			if err != nil {
				returnedError = nil
				break
			}

			if httpResp.StatusCode != 200 {

				if redirectedToHTTPS { // proxy

					if httpResp.StatusCode >= 520 && httpResp.StatusCode <= 527 { // Cloudflare

						// Cloudflare tries to connect through
						// HTTPS (SSL/TLS encryption mode is set to "Full")

						returnedError = entities.ErrCloudflareSSLFull{
							Domain: domain,
						}
						return
					}

					returnedError = entities.ErrProxyForceHTTPS{
						Domain: domain,
					}
					return
				}

				returnedError = nil
				break
			}

			httpBody, err := io.ReadAll(httpResp.Body)
			httpResp.Body.Close()

			if err != nil {
				returnedError = err
				break
			}

			returnedError = nil

			if string(httpBody) == uniqueID {

				reachable = true

				// We want to check that we don't have a proxy
				// in front of the sandbox that passes
				// all requests as HTTP (like Cloudflare with SSL mode set to "flexible")
				httpsClient := &http.Client{
					Timeout: 4 * time.Second,
				}

				httpsResp, err := httpsClient.Get("https://" + domain)

				if err == nil && httpsResp.StatusCode == 200 {
					httpsBody, err := io.ReadAll(httpsResp.Body)
					httpsResp.Body.Close()

					if err == nil && string(httpsBody) == uniqueID {
						// OK. We have a proxy in front of the sandbox that passes
						// all requests as HTTP (like Cloudflare with SSL mode set to "flexible").
						// If we redirect HTTP to HTTPS at the server level,
						// we will get a redirect loop.
						redirToHTTPS = false
						return
					}
				}

				redirToHTTPS = true
				return
			}
		} // <- end of select

		time.Sleep(pollSleepDuration)
	}
}

func WaitUntilDomainIsReachableViaHTTPS(
	domain string,
	timeout time.Duration,
) (returnedError error) {

	pollTimeoutChan := time.After(timeout)
	pollSleepDuration := 1 * time.Second

	for {
		select {
		case <-pollTimeoutChan:
			return
		default:
			httpsClient := &http.Client{
				Timeout: 4 * time.Second,
			}

			httpsResp, err := httpsClient.Get("https://" + domain)

			if err != nil {
				returnedError = err
				break
			}

			if httpsResp.StatusCode >= 520 && httpsResp.StatusCode <= 527 {
				returnedError = fmt.Errorf(
					"cloudflare cannot reach origin (HTTP status code %d)",
					httpsResp.StatusCode,
				)
				break
			}

			returnedError = nil
			return
		} // <- end of select

		time.Sleep(pollSleepDuration)
	}
}
