# CONTRIBUTING
## Makefile
### コンテナ

```$ make up```

```$ make down```

### ドキュメント生成

新しいエンドポイント作成時に必要

```$ make swag```

### 依存性注入(Dependency Injection)

新しいcontroller, usecase, repository, middleware作成時に必要

intrnal/app/wire.go に Newhoge() を含めてから実行

```$ make wire```
