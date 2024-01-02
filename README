<div>
svart
</div>

____

svart is a small command-line utility for working with environment variables in Terraform. Tools like [direnv](https://direnv.net) or [Chamber](https://github.com/segmentio/chamber) allow you to automatically export environment variables for your application configuration, e.g. from a dotenv file or from a secrets managers.

If you want to be able to read these values inside your Terraform configuration, you would manually have to re-export 


## Usage

The contents of a .env file can be re-exported by piping the contents to svart:

```shell
cat .env | svart
```

In order to use svart in conjunction with direnv, it suffices to add the following to your .envrc file file:

```shell
<(cat .envrc | svart)
```

svart will ignore all the input