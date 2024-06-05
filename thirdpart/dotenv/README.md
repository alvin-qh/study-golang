# DotEnv

## 2. 管理 `.env` 文件

本例中的 `.env` 文件通过 DotEnv 工具管理, 且内容已被加密, 所以 `.env` 源文件并不会提交到 Git 中

如果 `.env` 文件还不存在, 则需要通过如下方式拉取该文件

```bash
npx dotenv-vault@latest pull
```

如果对 `.env` 文件进行了修改, 则需要重新提交

```bash
npx dotenv-vault@latest push
```
