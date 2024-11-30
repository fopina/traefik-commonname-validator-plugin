# CommonName Validator

CommonName Validator is a middleware plugin for [Traefik](https://github.com/traefik/traefik) which allows authorizing mTLS requests based on subject CN.

## Configuration

### Static

```toml
[experimental.plugins.cnvalidator]
  modulename = "github.com/fopina/traefik-commonname-validator-plugin"
  version = "v0.0.1"
```

Or via cli flag

```yaml
...
--experimental.plugins.cnvalidator.modulename=github.com/fopina/traefik-commonname-validator-plugin \
--experimental.plugins.cnvalidator.version=v0.0.1 \
...
```

### Dynamic

To configure the `CommonName Validator` plugin you should create a [middleware](https://doc.traefik.io/traefik/middlewares/overview/) in 
your dynamic configuration as explained [here](https://doc.traefik.io/traefik/middlewares/overview/).
The following example creates and uses the `cnvalidator` middleware plugin to authorize requests with valid mTLS certificates that have the subject CN of `auth-client` and `auth2-client`. Other subjects will get a 403.

```toml
[http.routers]
  [http.routers.my-router]
    rule = "Host(`localhost`)"
    middlewares = ["allow-cn"]
    service = "my-service"

[http.middlewares]
  [http.middlewares.allow-cn.plugin.cnvalidator]
    allowed = ["auth-client", "auth2-client"]
    # uncomment to enable debug to print out CNs of rejected requests
    # debug = true

[http.services]
  [http.services.my-service]
    [http.services.my-service.loadBalancer]
      [[http.services.my-service.loadBalancer.servers]]
        url = "http://127.0.0.1"
```

Or via compose labels

```yaml
whoami:
    image: traefik/whoami
    labels:
      # uncomment to enable debug to print out CNs of rejected requests
      # traefik.http.middlewares.allow-cn.plugin.cnvalidator.debug: 'true'
      traefik.http.middlewares.allow-cn.plugin.cnvalidator.allowed[0]: auth-client
      traefik.http.middlewares.allow-cn.plugin.cnvalidator.allowed[1]: auth2-client
      traefik.http.routers.whoami-plugin.rule: Host(`localhost`)
      traefik.http.routers.whoami-plugin.entrypoints: webtls
      traefik.http.routers.whoami-plugin.tls: "true"
      traefik.http.routers.whoami-plugin.tls.options: ...
      traefik.http.routers.whoami-plugin.middlewares: allow-cn
```