<h1 align="center">
    <a href="https://github.com/praveenprem/testbed-slack-bot"><img src="images/code.png" width="25" height="25" alt="Bot icon"></a>
  Test Bed
</h1>

## About

Test Bed is a Slack code execution plugin powered by [Piston - general purpose code execution engine](https://github.com/engineer-man/piston).

This repository only contains the code for the Slack integration

## Installation

<a href="https://slack.com/oauth/v2/authorize?client_id=1722472343379.1722479188979&scope=commands&user_scope=">
    <img alt="Add to Slack" height="40" width="139" src="https://platform.slack-edge.com/img/add_to_slack.png" srcSet="https://platform.slack-edge.com/img/add_to_slack.png 1x, https://platform.slack-edge.com/img/add_to_slack@2x.png 2x" />
</a>

## Usage

### Inline usage prompt

- `/tb-help` command will provide guide on how to use with in Slack
- `/tb-lang` will provide a list of languages supported by the bot
- `/tb` will executed the give command and respond with the output

>/tb [lag]
>```
>code
>```
> [arg] [arg] [arg]...

>/tb python3
>```
>from sys import argv
>print(f'Hello {argv[1]} {argv[2]}')
>```
>John Smith
>
> _Output_
>
