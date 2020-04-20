# GitHub/GitLab Labels

This tiny command line help you to configure your GitHub label across your org.

For now only for GitHub, GitLab comes next

There are two commands: `set` and `list`

To list and save the labels in a YAML format file you can run

```shell
$ github-gitlab-labels list --github-token xoxoxoxoxo --repo github-labels --org cpanato --save
```

To configure the labels in a specific repo

```shell
$ github-gitlab-labels set  --github-token xoxoxoxoxo --label-file labels-sample.yaml --repo github-labels --org cpanato
```

