# go-module-diagnostic-gateway
Gateway/Proxy to go module protocol proxies with diagnostic logging.

## Rationale
The diagnostic proxy program was created to help find errors in JFrog Artifactory installed at my employer.

This program accepts http requests that conform to go module protocol. Then it sends an https request to the proxy servers on the list, logs what it got in response and sends back to the client what it downloaded.

This way you can see the request parameters and the response content in the log. This will hopefully allow us to get a better look at JFrog Artifactory.
