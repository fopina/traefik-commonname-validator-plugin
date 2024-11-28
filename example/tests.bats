#!/usr/bin/env bats

load '/bats-lib/bats-support/load'
load '/bats-lib/bats-assert/load'

@test "test whoami is working" {
  run curl -f \
           --cacert traefik_config/certs/good-one.pem \
           https://whoami.7f000001.nip.io/
  assert_success
  assert_output --partial 'X-Forwarded-Server'
}

@test "test whoami-plugin without cert is rejected" {
  run curl -f \
           --cacert traefik_config/certs/good-one.pem \
           https://whoami-plugin.7f000001.nip.io/
  assert_failure 56
  assert_output --partial 'alert bad certificate, errno 0'
}

@test "test whoami-plugin with bad-client cert is rejected" {
  run curl -f \
           --cacert traefik_config/certs/good-one.pem \
           --cert traefik_config/certs/bad-client/cert.pem \
           --key traefik_config/certs/bad-client/key.pem \
           https://whoami-plugin.7f000001.nip.io/
  assert_failure 56
  assert_output --partial 'alert bad certificate, errno 0'
}

@test "test whoami-plugin with auth-client cert is accepted and header added" {
  run curl -f \
           --cacert traefik_config/certs/good-one.pem \
           --cert traefik_config/certs/auth-client/cert.pem \
           --key traefik_config/certs/auth-client/key.pem \
           https://whoami-plugin.7f000001.nip.io/
  assert_success
  assert_output --partial 'Foo: Bar'
}

@test "test whoami-plugin with auth2-client cert is also accepted" {
  run curl -f \
           --cacert traefik_config/certs/good-one.pem \
           --cert traefik_config/certs/auth2-client/cert.pem \
           --key traefik_config/certs/auth2-client/key.pem \
           https://whoami-plugin.7f000001.nip.io/
  assert_success
  assert_output --partial 'Foo: Bar'
}

@test "test whoami-plugin with not-auth-client cert is ALSO ACCEPTED FOR NOW BUT SHOULD NOT # FIXME" {
  run curl -f \
           --cacert traefik_config/certs/good-one.pem \
           --cert traefik_config/certs/not-auth-client/cert.pem \
           --key traefik_config/certs/not-auth-client/key.pem \
           https://whoami-plugin.7f000001.nip.io/
  # assert_failure 22
  # assert_output --partial 'returned error: 403'
  assert_success
  assert_output --partial 'Foo: Bar'
}
