# Terraform azurerm-restapi Provider Language Server

Experimental version of [terraform-provider-azurerm-restapi](https://github.com/Azure/terraform-provider-azurerm-restapi) language server.

## What is LSP

Read more about the Language Server Protocol at https://microsoft.github.io/language-server-protocol/

## Introduction

This project only supports completion/hover/diagnostics for `terraform-provider-azurerm-restapi`,
not targeting support all language features for `HCL` or `Terraform`. To get the best user experience, 
it's recommended to use it with language server for `Terraform`.

## Features

- [x] Completion when input `type`
- [x] Completion when input `body`, limitation: it only works when use `jsonencode` function to build the JSON
- [X] Show hint when hover on `type`, `body` and properties defined inside `body`
- [X] Show diagnostics for properties defined inside `body`
- [x] Better completion for discriminated object

## Installation

1. Clone this project to local
2. Run `go install` under the project folder.

## Usage

The most reasonable way you will interact with the language server
is through a client represented by an IDE, or a plugin of an IDE.

VSCode extension: [azurerm-restapi-vscode](https://github.com/ms-henglu/azurerm-restapi-vscode)
