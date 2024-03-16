# GithubSecrets
Update the Github secrets for a given repository

## Purpose
Right now credentials stored in AWS Parameter Store will be fetched, encoded 
and send to Github. AWS Profile is necessary, but the Github user and Github API
key can be overwritten. Otherwise it uses a default user and fetches the API key
from the parameter store.

This is no general purpose tool and might need some more modification to make it
work for for edge cases.

## Example
```shell
githubsecrets --github_repo inventory --aws-profile developer
```
If a custom text is given it will either encrypted or decrypted and printed to
stdout instead of sending it to Github.

