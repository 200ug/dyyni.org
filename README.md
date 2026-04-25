## development

Development setup with [stormdrain](https://codeberg.org/2ug/stormdrain):

```shell
stormdrain new -f profile.json
stomdrain enter

# inside container
npm install
hugo server --ignoreCache --bind 0.0.0.0
```

