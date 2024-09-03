<div align="center" width="100%">
  <img src="./assets/logo.png" alt="Git sync logo" width="150" />
</div>
<div align="center" width="100%">
    <h2>🔄 git-sync</h2>
    <p>A simple tool to backup and sync your git repositories</p>
    <a target="_blank" href="https://github.com/AkashRajpurohit/git-sync/actions"><img src="https://github.com/AkashRajpurohit/git-sync/actions/workflows/release.yml/badge.svg?event=push" /></a>
    <a href="https://goreportcard.com/report/github.com/AkashRajpurohit/git-sync"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/AkashRajpurohit/git-sync">
    <a target="_blank" href="https://github.com/AkashRajpurohit/git-sync/releases"><img src="https://img.shields.io/github/downloads/AkashRajpurohit/git-sync/total" /></a>
    <img alt="Visitors count" src="https://visitor-badge.laobi.icu/badge?page_id=@akashrajpurohit~git-sync.visitor-badge&style=flat-square">
    <a target="_blank" href="https://github.com/AkashRajpurohit/git-sync/releases"><img src="https://img.shields.io/github/go-mod/go-version/AkashRajpurohit/git-sync?filename=go.mod" /></a>
    <a target="_blank" href="https://github.com/AkashRajpurohit/git-sync/releases"><img src="https://img.shields.io/github/v/release/AkashRajpurohit/git-sync?display_name=tag" /></a>
    <a href="#-contributors"><img alt="All Contributors" src="https://img.shields.io/github/all-contributors/AkashRajpurohit/git-sync?color=1f85bf"></a>
    <a target="_blank" href="https://github.com/AkashRajpurohit/git-sync"><img src="https://img.shields.io/github/stars/AkashRajpurohit/git-sync" /></a>
    <br />
    <br />
    <p align="center">
      <a href="https://github.com/AkashRajpurohit/git-sync/issues/new?template=bug_report.yml">Bug report</a>
      ·
      <a href="https://github.com/AkashRajpurohit/git-sync/issues/new?template=feature_request.yml">Feature request</a>
      ·
      <a href="https://github.com/AkashRajpurohit/git-sync/wiki">Read Docs</a>
    </p>
</div>
<hr />

`git-sync` is a CLI tool designed to help you back up your Git repositories. This tool ensures you have a local copy of your repositories, safeguarding against potential issues such as account bans or data loss.

## 📺 Demo

[![asciicast](./assets/asciinema.svg)](https://asciinema.org/a/674282)

## 🤔 Why `git-sync`?

Remember when `@defunkt` [GitHub account got banned?](https://twitter.com/defunkt/status/1754610843361362360) Well, he is the co-founder of GitHub so he did get his account un-banned pretty quickly but what if you are not that lucky?

Recently I have seen many developers [getting their GitHub account banned](https://www.reddit.com/r/github/search/?q=account+got+banned&sort=new) and losing access to their repositories. Some may be able to recover their account (but there is delay in that) and some may not be able to recover their account at all. What would you do if you lose access to your repositories? What if GitHub goes down? What if you accidentally delete your repositories? The answer is simple, you should have a backup of your repositories.

`git-sync` provides a straightforward way to back up all your repositories locally, ensuring you have access to your code whenever you need it. It does this by doing a bare clone of all your repositories in a specified directory so that you can recover your code in case of any unforeseen circumstances as well as the file size of your backups is minimal.

## ✨ Features

- **Backup All Repositories:** Automatically clone or update all your GitHub repositories to a local directory.
- **Periodic Sync:** Keep your backups in sync with your remote repositories by running `git-sync` [periodically](https://github.com/AkashRajpurohit/git-sync/wiki/Setup-Periodic-Backups).
- **Bare Clone:** Efficiently back up repositories using bare clones to save space and speed up the process.
- **Concurrency:** Sync multiple repositories concurrently to reduce the time required for backup.
- **Configuration File:** Easily manage your settings through a YAML configuration file.
- **Custom Backup Directory:** Specify the directory where you want to store your repositories.
- **Multi Platform:** Currently this project supports backing up repositories from all major Git hosting services like GitHub, GitLab, and Bitbucket.

## 🚀 Getting Started

We have a thorough guide on how to set up and get started with `git-sync` in our [documentation](https://github.com/AkashRajpurohit/git-sync/wiki).

## 🙏🏻 Support

If you found the project helpful, consider giving it a star ⭐️. If you would like to support the project in other ways, you can [buy me a coffee](https://ko-fi.com/akashrajpurohit) or [sponsor me on GitHub](https://github.com/sponsors/AkashRajpurohit).

## 🐛 Bugs or Requests

If you encounter any problems feel free to open an [issue](https://github.com/AkashRajpurohit/git-sync/issues/new?template=bug_report.yml). If you feel the project is missing a feature, please raise a [ticket](https://github.com/AkashRajpurohit/git-sync/issues/new?template=feature_request.yml) on GitHub and I'll look into it. Pull requests are also welcome.

## 🫱🏻‍🫲🏼 Contributors

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://akashrajpurohit.com/?ref=git-sync"><img src="https://avatars.githubusercontent.com/u/30044630?v=4?s=100" width="100px;" alt="Akash Rajpurohit"/><br /><sub><b>Akash Rajpurohit</b></sub></a><br /><a href="#code-AkashRajpurohit" title="Code">💻</a> <a href="#ideas-AkashRajpurohit" title="Ideas, Planning, & Feedback">🤔</a> <a href="#review-AkashRajpurohit" title="Reviewed Pull Requests">👀</a> <a href="#doc-AkashRajpurohit" title="Documentation">📖</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://joao.bonadiman.dev"><img src="https://avatars.githubusercontent.com/u/18357636?v=4?s=100" width="100px;" alt="João Vitor Bonadiman"/><br /><sub><b>João Vitor Bonadiman</b></sub></a><br /><a href="#code-jbonadiman" title="Code">💻</a> <a href="#ideas-jbonadiman" title="Ideas, Planning, & Feedback">🤔</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://qlaffont.com"><img src="https://avatars.githubusercontent.com/u/10044790?v=4?s=100" width="100px;" alt="Quentin Laffont"/><br /><sub><b>Quentin Laffont</b></sub></a><br /><a href="#code-qlaffont" title="Code">💻</a></td>
    </tr>
  </tbody>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

## 👀 Who am I?

[![Website Badge](https://img.shields.io/badge/-akashrajpurohit.com-3b5998?logo=google-chrome&logoColor=white)](https://akashrajpurohit.com/)
[![Linkedin Badge](https://img.shields.io/badge/-@AkashRajpurohit-0e76a8?logo=Linkedin&logoColor=white)](https://linkedin.com/in/AkashRajpurohit)
[![Twitter Badge](https://img.shields.io/twitter/follow/akashwhocodes)](https://twitter.com/AkashWhoCodes)
[![Mastodon Follow](https://img.shields.io/mastodon/follow/112372456922065040)](https://mastodon.social/@akashrajpurohit)
