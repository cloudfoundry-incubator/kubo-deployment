# Configuring CF routing for Kubo

CF routing requires that a CloudFoundry installation with TCP routing enabled is accessible,
and the credentials for various services can be obtained.

Follow the steps below in order to configure kubo to use CF routing:

1. Uncomment and fill in `routing-cf-client-secret` and `routing-cf-nats-password` in 
  `<KUBO_ENV>/director-secrets.yml`

2. Configure the CF routing settings in `<KUBO_ENV>/director.yml`:

  - Comment out all the lines grouped underneath the **IaaS routing mode settings** comment
  
  - Uncomment all the lines grouped underneath the **CF routing mode settings** comment 
    and fill in all the values.
  
