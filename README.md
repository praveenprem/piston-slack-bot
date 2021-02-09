<h1 align="center">
    <a href="https://github.com/praveenprem/testbed-slack-bot"><img src="images/code.png" width="25" height="25" alt="Bot icon"></a>
  Test Bed
</h1>

## About

Test Bed is a Slack code execution plugin powered by [Piston - general purpose code execution engine](https://github.com/engineer-man/piston).

This repository only contains the code for the Slack integration

## Installation

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
