Sample `docker-compose.yml` to illustrate how to set up a service (ie: `whoami`) with this middleware plugin.

`setup` and `test` are only defined to be able to run tests and do not make sense out of it, but the rest of the compose can be used as is.

`docker compose up` will bring up the whole stack:
* `setup` creates the required CAs and certificates (for server and for client tests)
* `test` will run [tests.bats](tests.bats) when all the others are up and running
    * It can be re-executed with `docker compose run test`

After bringing it up, `whoami` service can also be called from host:

```
curl https://whoami-plugin.7f000001.nip.io:8889 --cert traefik_config/certs/auth-client/cert.pem \
                                                --key traefik_config/certs/auth-client/key.pem \
                                                --cacert traefik_config/certs/good-one.pem
```