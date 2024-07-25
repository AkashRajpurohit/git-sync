<h1 align="center" style="border-bottom: none;">🔄 git-sync</h1>
<h3 align="center">A tool to backup and sync your git repositories</h3>
<br />
<p align="center">
  <a href="https://github.com/AkashRajpurohit/git-sync/actions/workflows/release.yml">
    <img alt="Build states" src="https://github.com/AkashRajpurohit/git-sync/actions/workflows/release.yml/badge.svg">
  </a>
  <a href="https://goreportcard.com/report/github.com/AkashRajpurohit/git-sync">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/AkashRajpurohit/git-sync">
  </a>
  <img alt="Visitors count" src="https://visitor-badge.laobi.icu/badge?page_id=@akashrajpurohit~git-sync.visitor-badge&style=flat-square">
  <a href="#contributors-">
    <img alt="All Contributors" src="https://img.shields.io/github/all-contributors/AkashRajpurohit/git-sync?color=1f85bf">
  </a>
  <a href="https://twitter.com/akashwhocodes">
    <img alt="follow on twitter" src="https://img.shields.io/twitter/follow/akashwhocodes.svg?style=social&label=@akashwhocodes">
  </a>

  <p align="center">
    <a href="https://github.com/AkashRajpurohit/git-sync/issues/new?template=bug_report.yml">Bug report</a>
    ·
    <a href="https://github.com/AkashRajpurohit/git-sync/issues/new?template=feature_request.yml">Feature request</a>
  </p>
</p>
<br />
<hr />

`git-sync` is a CLI tool designed to help you back up your GitHub repositories. This tool ensures you have a local copy of your repositories, safeguarding against potential issues such as account bans or data loss on GitHub.

By using `git-sync`, you can easily clone or update your repositories to a specified local directory.

[![asciicast](https://asciinema.org/a/664462.svg)](https://asciinema.org/a/664462)

## Why `git-sync`?

Remember when `@defunkt` [GitHub account got banned?](https://twitter.com/defunkt/status/1754610843361362360) Well, he is the co-founder of GitHub so he did get his account un-banned pretty quickly but what if you are not that lucky?

Recently I have seen many developers [getting their GitHub account banned](https://www.reddit.com/r/github/search/?q=account+got+banned&sort=new) and losing access to their repositories. Some may be able to recover their account (but there is delay in that) and some may not be able to recover their account at all. What would you do if you lose access to your repositories? What if GitHub goes down? What if you accidentally delete your repositories? The answer is simple, you should have a backup of your repositories.

`git-sync` provides a straightforward way to back up all your repositories locally, ensuring you have access to your code whenever you need it. It does this by doing a bare clone of all your repositories in a specified directory so that you can recover your code in case of any unforeseen circumstances as well as the file size of your backups is minimal.

## Features

- **Backup All Repositories:** Automatically clone or update all your GitHub repositories to a local directory.
- **Bare Clone:** Efficiently back up repositories using bare clones to save space and speed up the process.
- **Concurrency:** Sync multiple repositories concurrently to reduce the time required for backup.
- **Configuration File:** Easily manage your settings through a YAML configuration file.

## Installation

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.22 or later)
- [Git](https://git-scm.com/downloads)

### Using `go get`

```bash
go get github.com/AkashRajpurohit/git-sync
```

### Build from source

```bash
git clone https://github.com/AkashRajpurohit/git-sync.git
cd git-sync
go install
go build
```

### With Docker

```bash
docker run --rm -v /path/to/config/:/git-sync -v /path/to/backups:/backups ghcr.io/akashrajpurohit/git-sync:latest
```

Or you can use the `docker-compose.yml` file to run the container.

```yaml
services:
  git-sync:
    image: ghcr.io/akashrajpurohit/git-sync:latest
    volumes:
      - /path/to/config/:/git-sync
      - /path/to/backups:/backups
```

```bash
docker-compose up
```

### Download Pre-built Binaries

Pre-built binaries are available for various platforms. You can download the latest release from the [Releases](https://github.com/AkashRajpurohit/git-sync/releases) page.

## Usage

### Configuration

Before using `git-sync`, you need to create a configuration file named `config.yml`. The default path for the configuration file is `~/.config/git-sync/config.yml`.

Here's an example configuration file:

```yaml
# Configuration file for git-sync
# Default path: ~/.config/git-sync/config.yml
username: your-github-username
token: your-personal-access-token
include_repos: []
exclude_repos: []
include_orgs: []
exclude_orgs: []
backup_dir: /path/to/backup
include_forks: false
```

| Field          | Description                                                                                                      |
|----------------|------------------------------------------------------------------------------------------------------------------|
| `username`     | Your GitHub username.                                                                                            |
| `token`        | Your GitHub personal access token. You can create a new token [here](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#about-personal-access-tokens). Ensure that the token has the `repo` scope. |
| `include_repos` | A list of repositories to include in the backup. If set then only these repositories will be backed up.          |
| `exclude_repos` | A list of repositories to exclude from the backup. If set then all repositories except these will be backed up.  |
| `include_orgs` | A list of orgs to include in the backup. If set then repositories only from these organizations will be backed up          |
| `exclude_orgs` | A list of orgs to exclude from the backup. If set then all organization repositories except these orgs will be backed up.  |
| `backup_dir`   | The directory where the repositories will be backed up. Default is `~/git-backups`.                              |
| `include_forks`| If set to `true`, forks of the user's repositories will also be backed up. Default is `false`.                   |

> Note: The `include_repos`, `exclude_repos`, `include_orgs` and `exclude_orgs` fields accept repository/organization name as well as glob patterns. The patterns supported are those defined by [filepath.Match](https://pkg.go.dev/path/filepath#Match).

---

Here are some dry run examples for different configurations:

#### Include all repositories that I own (even forks)

```yaml
username: your-github-username
token: your-personal-access-token
include_repos: []
exclude_repos: []
include_orgs: []
exclude_orgs: []
backup_dir: /path/to/backup
include_forks: true
```

#### Include all repositories that I own excluding forks

```yaml
username: your-github-username
token: your-personal-access-token
include_repos: []
exclude_repos: []
include_orgs: []
exclude_orgs: []
backup_dir: /path/to/backup
include_forks: false
```

#### Include all repositories that I own excluding forks and some repositories

```yaml
username: your-github-username
token: your-personal-access-token
include_repos: []
exclude_repos: ["repo1", "repo2"]
include_orgs: []
exclude_orgs: []
backup_dir: /path/to/backup
include_forks: false
```

#### Include only specific repositories

```yaml
username: your-github-username
token: your-personal-access-token
include_repos: ["repo1"]
exclude_repos: []
include_orgs: []
exclude_orgs: []
backup_dir: /path/to/backup
include_forks: false
```

**NOTE:** In the above example if you passed `exclude_repos` as well then it will be ignored since you have passed `include_repos`. So an example like below will still have same effect as above and only "repo1" will be backed up.

```yaml
username: your-github-username
token: your-personal-access-token
include_repos: ["repo1"]
exclude_repos: ["repo2", "repo3"]
include_orgs: []
exclude_orgs: []
backup_dir: /path/to/backup
include_forks: false
```

#### Include / Exclude repositories only from specific organizations

```yaml
username: your-github-username
token: your-personal-access-token
include_repos: []
exclude_repos: []
include_orgs: ["org1", "org2"]
exclude_orgs: []
backup_dir: /path/to/backup
include_forks: false
```

Similarly you can exclude repositories from specific organizations:

```yaml
username: your-github-username
token: your-personal-access-token
include_repos: []
exclude_repos: []
include_orgs: []
exclude_orgs: ["org1", "org2"]
backup_dir: /path/to/backup
include_forks: false
```

> **NOTE:** When you are using `exclude_orgs` option, it means that all _organization_ repositories from the organizations mentioned in the list will be excluded from the backup. If you want to exclude only specific repositories from the organizations then you can use `exclude_repos` option.

#### Include / Exclude repositories using glob patterns

To include all repos which start with `repo`, you can use the following configuration:

```yaml
username: your-github-username
token: your-personal-access-token
include_repos: ["repo*"]
exclude_repos: []
include_orgs: []
exclude_orgs: []
backup_dir: /path/to/backup
include_forks: false
```

Same way you can exclude all repos which start with `repo`:

```yaml
username: your-github-username
token: your-personal-access-token
include_repos: []
exclude_repos: ["repo*"]
include_orgs: []
exclude_orgs: []
backup_dir: /path/to/backup
include_forks: false
```

### Commands

#### Sync Repositories

To sync your repositories, run the following command:

```bash
git-sync
```

This command will clone or update all your repositories to the specified backup directory.

#### Version

To check the version of `git-sync`, run the following command:

```bash
git-sync version
```

#### Help

To view the help message, run the following command:

```bash
git-sync help
```

### Setup Periodic Backups

#### Unix-based Systems

You can set up periodic backups using [cron jobs or systemD timers](https://akashrajpurohit.com/blog/systemd-timers-vs-cron-jobs/?ref=git-sync). For example, to back up your repositories every day at 12:00 AM, you can add the following cron job:

```bash
0 0 * * * /path/to/git-sync
```

Replace `/path/to/git-sync` with the path to the `git-sync` binary.

#### Windows

You can set up periodic backups using Task Scheduler. Here's how you can do it:

1. Open Task Scheduler.
2. Click on `Create Basic Task`.
3. Enter a name and description for the task.
4. Choose the trigger (e.g., `Daily`).
5. Set the time for the trigger.
6. Choose `Start a program` as the action.
7. Browse to the `git-sync` binary.
8. Click `Finish` to create the task.
9. Right-click on the task and select `Run` to test it.
10. Your repositories will now be backed up periodically.

Or you can use Powershell script to run the `git-sync` binary.

```powershell
$action = New-ScheduledTaskAction -Execute "path\to\git-sync.exe"
$trigger = New-ScheduledTaskTrigger -Daily -At "12:00AM"
Register-ScheduledTask -Action $action -Trigger $trigger -TaskName "GitSyncTask" -Description "Daily Git Sync"
```

Replace `path\to\git-sync.exe` with the path to the `git-sync` binary.

## FAQs

### How do I create a GitHub personal access token?

You can create a new personal access token by following the steps mentioned here: [Creating a personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token). If you don't want to rotating the tokens after expiry, select the `(Classic)` option while creating the token and set expiration to `Never`.

### How do I restore from the backups?

The backups generated by `git-sync` are essentially bare clones of your repositories. To restore a repository, you can simply clone it from the backup directory. For example:

```bash
git clone /path/to/backup/repo.git
```

## Bugs or Requests 🐛

If you encounter any problems feel free to open an [issue](https://github.com/AkashRajpurohit/git-sync/issues/new?template=bug_report.yml). If you feel the project is missing a feature, please raise a [ticket](https://github.com/AkashRajpurohit/git-sync/issues/new?template=feature_request.yml) on GitHub and I'll look into it. Pull requests are also welcome.

## Where to find me? 👀

[![Website Badge](https://img.shields.io/badge/-akashrajpurohit.com-3b5998?logo=google-chrome&logoColor=white)](https://akashrajpurohit.com/)
[![Twitter Badge](https://img.shields.io/badge/-@akashwhocodes-00acee?logo=Twitter&logoColor=white)](https://twitter.com/AkashWhoCodes)
[![Linkedin Badge](https://img.shields.io/badge/-@AkashRajpurohit-0e76a8?logo=Linkedin&logoColor=white)](https://linkedin.com/in/AkashRajpurohit)
[![Instagram Badge](https://img.shields.io/badge/-@akashwho.codes-e4405f?logo=Instagram&logoColor=white)](https://instagram.com/akashwho.codes/)
[![Telegram Badge](https://img.shields.io/badge/-@AkashRajpurohit-0088cc?logo=Telegram&logoColor=white)](https://t.me/AkashRajpurohit)

## Contributors 🫱🏻‍🫲🏼

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