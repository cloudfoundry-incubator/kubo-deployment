# Creating a Kubo environment

A Kubo environment is a set of configuration files used to deploy and update 
both KuBOSH and Kubo. 

Run `./bin/generate_env_config <ENV_PATH> <ENV_NAME> gcp` to generate a Kubo
environment. The environment will be referred to as `KUBO_ENV` in this guide,
and will be located at `<ENV_PATH>/<ENV_NAME>`. 

Follow the comments in `<KUBO_ENV>/director.yml` to fill in the values. 
You might need to fill in the values in `<KUBO_ENV>/director-secrets.yml` as 
well.
 
> Run `bin/generate_env_config --help` for more detailed information.