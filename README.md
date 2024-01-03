<div align="center">
    <picture>
        <source media="(prefers-color-scheme: dark; max-width: 100px)" srcset="https://github.com/michielvermeir/svart/blob/main/logo/svart.jpg">
        <source media="(prefers-color-scheme: light; max-width: 100px)" srcset="https://github.com/michielvermeir/svart/blob/main/logo/svart.jpg">
        <img alt="The svart logo, a stylized letter S " width="100" src="https://github.com/michielvermeir/svart/blob/main/logo/svart.jpg">
    </picture>
    <h1 style="text-transform: lowercase;">Svart</h1>
</div>

____

svart is a small command-line utility for working with environment variables in Terraform. Tools like [direnv](https://direnv.net) or [Chamber](https://github.com/segmentio/chamber) allow you to automatically export environment variables for your application configuration, e.g. from a dotenv file or from a secrets managers.

Terraform and OpenTofu require environment variables to be prefixed with `TF_VAR_` before they can be read from within the configuration. If you need environment variables inside your Terraform configuration, exported by tools like direnv or chamber, they would manually need to be re-exported. 

svart pre

## Usage

> [!NOTE] Strict Mode
> Svart runs in strict mode by default unless told differently. This means svart will not output anything unless allowed patterns are configured configured. Strict mode can be (temporarily) disabled by using `--relaxed` mode.

Imagine the following `.env` file:

```
AWS_FOO=egg
AWS_BAR=spam
```
_An example .env file_



```shell

```

The contents of a .env file can be re-exported by piping the contents to svart:

```shell
cat .env | svart --from-stdin --relaxed
```

In order to use svart in conjunction with direnv, it suffices to add the following to your .envrc file file:

```shell
<(cat .envrc | svart)
```

svart will ignore all the input