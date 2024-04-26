# teu

GitOps, made easy.

> **Note:** teu is currently in v0. This means that it is still under development and may have some bugs or incomplete features. Please use it with caution and report any issues you encounter. It's not recommended yet to use teu in production.



## What is teu?
teu is a program writen in go which allows you to orchestrate containers in a dead simple way.

You write a teu.yml file which points to many docker compose files, which have everything needed to run your application. This all fits inside a single repository.

Changing anything inside the repository will automatically update any nodes that are joined to the repository. (GitOps)

Teu also supports storing secrets in the repository, by encrypting them using age. Secrets get decrypted at runtime.

## How does it work?

1. You write a teu.yml file which points to many docker compose files. (More on this later)
2. You publish everything inside of a public Git repository.
3. You run `teu join` as root on a node, which will clone the repository and start the containers.
4. Every 5 minutes, teu will pull the repository and restart any containers if needed.
5. Teu is smart! It will only restart containers if the deployment has changed.
6. You can run `teu status` to see the status of all deployments.
7. Once you are done, you can run `teu leave` to remove the containers and the repository.


## What's inside a teu.yml file?

A teu.yml file is a simple yaml file which points to a directory containing a docker-compose file.

Here's a really basic example:

```yaml
slate:
  name: production
  description: Production environment
  age_secret_key: /root/secrets/age-secret.key

applications:
  - name: SearXNG
    deploy: ./deploy/searxng

  - name: Nginx
    deploy: ./deploy/nginx
    secrets:
      - tls.crt.age
      - tls.key.age
```

In this example, we have two applications: SearXNG and Nginx. SearXNG is a simple application which doesn't require any secrets, while Nginx requires a TLS certificate and key. The TLS certificate and key are stored in the repository as encrypted files.

The application name is the identifier for the application.
The deploy key points to a directory containing a docker-compose file and typically all other files needed to run the application.
The secrets key is optional and points to a list of encrypted files which will be decrypted at runtime.

### How do I encrypt secrets?

You can encrypt secrets using the `age` command line tool. Here's an example:

```bash
age-keygen
cat tls.crt | age -r age1... > tls.crt.age
cat tls.key | age -r age1... > tls.key.age
```

You can then commit the `tls.crt.age` and `tls.key.age` files to your repository.
