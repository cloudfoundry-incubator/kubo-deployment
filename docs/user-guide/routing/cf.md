# Configuring CF routing for Kubo

CF routing requires that a Cloud Foundry installation with TCP routing enabled is accessible,
and the credentials for various services can be obtained. In order to configure kubo to use
Cloud Foundry routing, the following CF settings should be available:
  
  - Cloud Foundry TCP Router hostname
  - Cloud Foundry apps domain
  - Cloud Foundry API URL
  - Cloud Foundry UAA URL and credentials for a client that can access the TCP router
  - Cloud Foundry NATS bus information - ip addresses, username, port and access credentials

Follow the steps below in order to configure Kubo to use CF routing from the bastion box:

1. Uncomment and fill in the values for `routing-cf-client-secret` and `routing-cf-nats-password` in 
  `<KUBO_ENV>/director-secrets.yml`

1. Configure the CF routing settings in `<KUBO_ENV>/director.yml`:

  - Comment out all the lines grouped underneath the **IaaS routing mode settings** comment
  
  - Uncomment all the lines grouped underneath the **CF routing mode settings** comment 
    and fill in all the values.
  
