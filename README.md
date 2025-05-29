# 🌃 Thoth (Go)

A language agnostic tool for turning TODO lines / comments in a code base into open Github issues.

## 🚀 Features

- Finds all the TODO lines in the current folder 
- Finds all the open issues in your github - using git remote 
- Checks to see whether or not the issue is in github 
    - If it is not on GitHub in will add a issue number to the start of the todo line
    - If it is on GitHub it will ignore the issue 

## 🛠️ Prerequisites

- [Go](https://golang.org/dl/) installed (version 1.16+ recommended)
- A github token for the repository with permission to read / edit issues 

## 📁 Setup

1. Clone this repository:

   ```bash
   git clone https://github.com/jonathon-chew/Thoth.git
   cd Thoth 
   ```

2. Compile the script:

    `go build .`

## 📂 Output

This will make Github issues for you automatically and edit your codebase - just the todo line, to save the number of the issue for easily finding which issue is the right issue.

## 🧠 Notes

This is inspired by the project here: https://github.com/tsoding/snitch

## 📜 License

This project is licensed under the MIT License. See the LICENSE file for details.
