# ci-github-actions-go

[![Coverage Status](https://coveralls.io/repos/github/syunkitada/ci-github-actions-go/badge.svg?branch=mod-actions)](https://coveralls.io/github/syunkitada/ci-github-actions-go?branch=mod-actions)

- github actions のテスト用
- [GitHub Actions ドキュメント](https://docs.github.com/en/actions)

## メモ

- coverallsapp
  - 手順
    - github でサインアップだけしておく
    - coverallsapp/github-action を利用してカバレッジデータを保存する
    - 使い方は、[こちら](https://github.com/coverallsapp/github-action)
    - 設定の中で、github-token: ${{ secrets.GITHUB_TOKEN }} を設定する必要があるが、これはおまじないとして書くだけ
      - GITHUB_TOKEN はデフォルトで組み込まれてるので、明示的に GitHub 側で secret を設定する必要はない
  - badge ステータスの追加方法は以下を参照
    - https://github.com/OpenSourceHelpCommunity/OpenSourceHelpCommunity.github.io/issues/83)
