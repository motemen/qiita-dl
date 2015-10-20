# qiita-dl

qiita-dl is a simple tool that donwloads snippets published on Qiita <http://qiita.com>.

    qiita-dl [-x] [-o <name>] [-d <directory>] <url>

## Example

To download a snippet under ~/bin and make it executable:

    $ qiita-dl -x -d ~/bin http://qiita.com/uasi/items/57da2e4268d348b371fb
    Title: "git commit --fixup で fixup する対象を peco/fzf で選べるスクリプト書いた"
    Saved to ~/bin/git-fixup

## Installation

    go get -u github.com/motemen/qiita-dl

## Author

motemen <https://motemen.github.io/>
