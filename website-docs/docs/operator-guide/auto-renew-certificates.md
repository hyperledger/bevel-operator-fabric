---
id: auto-renew-certificates
title: Automatic renewal of certificates
---

## Auto renewal of certificates

The following command line parameters are available for the operator:

- `auto-renew-orderer-certificates-delta`
- `auto-renew-peer-certificates-delta`
- `auto-renew-identity-certificates-delta`
- `auto-renew-peer-certificates`
- `auto-renew-orderer-certificates`
- `auto-renew-identity-certificates`


### Auto renewal of certificates for orderers

This operator supports the auto renewal of certificates for orderers. The following command line parameters are available for the operator:

- `auto-renew-orderer-certificates-delta`
- `auto-renew-orderer-certificates`

The `auto-renew-orderer-certificates-delta` parameter specifies the number of days before the expiration of the orderer certificates when the operator should start the renewal process. The default value is `15` days.

The `auto-renew-orderer-certificates` parameter enables the auto renewal of certificates for orderers. The default value is `false`.

### Auto renewal of certificates for peers

This operator supports the auto renewal of certificates for peers. The following command line parameters are available for the operator:

- `auto-renew-peer-certificates-delta`
- `auto-renew-peer-certificates`

The `auto-renew-peer-certificates-delta` parameter specifies the number of days before the expiration of the peer certificates when the operator should start the renewal process. The default value is `15` days.

The `auto-renew-peer-certificates` parameter enables the auto renewal of certificates for peers. The default value is `false`.


### Auto renewal of certificates for identities

This operator supports the auto renewal of certificates for identities. The following command line parameters are available for the operator:

- `auto-renew-identity-certificates-delta`
- `auto-renew-identity-certificates`

The `auto-renew-identity-certificates-delta` parameter specifies the number of days before the expiration of the identity certificates when the operator should start the renewal process. The default value is `15` days.

The `auto-renew-identity-certificates` parameter enables the auto renewal of certificates for identities. The default value is `false`.

