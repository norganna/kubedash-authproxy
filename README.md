# kdash - authenticating kubernetes dashboard proxy.

AWS IAM tokens used to authenticate to Kubernetes have a notoriously short duration, and it can be extremely frustrating to use the dashboard because of this.

This is because every 15 minutes you need to re-run the `aws-iam-authentication` command to output a token, copy the token, switch to your browser, click logout on the kubernetes dashboard, select the token radio box, select the token input field, paste in your clipboard token, submit, and then re-navigate back to where you were before this whole ordeal began (because you get directed back to the home page after loging in and your namespace is set back to default).

This application is a proxy I had to write because this whole process annoyed me so deeply I was not able to do my actual job anymore because of the constant searches on the internet to try and find a solution to this deeply upsetting problem.

kdash will start up a server on a local port (8002 by default) and when you visit it, will retrieve and automatically refresh your AWS IAM token before it expires (i.e. every 10 minutes).

When it has authenticated to the kubernetes dashboard app for you, it automatically injects the authentication details into the requests the webpage makes that it is forwarding to the actual dashboard proxy.

## Installing

```bash
go get github.com/norganna/kdash
go install github.com/norganna/kdash
```

## Running

First start up your kubernetes proxy:

```bash
kubectl proxy
```

Now run the kdash proxy, substituting the cluster and role you would normally supply to `aws-iam-authentcation` command:

```bash
kdash --cluster clusterName --role arn:aws:iam::12345678:role/roleName
```

If you can't find the kdash application, you may not have the $GOPATH/bin folder in your search path, you can copy or
link the binary to a suitable place in your path.

Once kdash is running, open your browser to [http://localhost:8002](http://localhost:8002)

## Options

```
kdash --help
Usage of kdash:
      --authenticator string   The path the the AWS IAM Authenticator binary (default "/usr/local/bin/aws-iam-authenticator")
      --cluster string         The name of the cluster to pass to the authentication
      --listen string          Where to listen for connections (default "localhost:8002")
      --proxy string           The proxy's location (default "http://localhost:8001")
      --role string            The role ARN to pass to the authenticator
```

You can also create a `~/.kdash/config.yaml` file which contains these options to save you having to specify them every time, for example:

```yaml
cluster: clusterName
role: arn:aws:iam::12345678:role/roleName
```

Alternatively any of these options can be supplied via an environment variable prefixed with `KDASH_`, eg:

```bash
export KDASH_CLUSTER=clusterName
export KDASH_ROLE=arn:aws:iam::12345678:role/roleName
```
