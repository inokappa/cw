# cw

## これなに

* YAML で定義したパラメータの CloudWatch メトリクスを取得するだけのツールです

## 使い方

### cw の取得

以下の通り, cw バイナリをパスが通ったディレクトリにぽいっと置くだけです.

```
wget https://github.com/inokappa/cw/releases/download/v0.0.1/cw_{{.OS}}_{{.Arch}} -O /path/to/dir/cw
chmod +x /path/to/dir/cw
```

### ヘルプ

```sh
$ ./cw -help
Usage of ./cw:
  -config string
        YAML ファイルを指定.
  -endpoint string
        AWS API のエンドポイントを指定.
  -profile string
        Profile 名を指定.
  -region string
        Region 名を指定. (default "ap-northeast-1")
  -target string
        YAML ファイル内のターゲット名を指定.
  -version
        バージョンを出力.
```

### YAML ファイルの書き方

YAML ファイルは以下のように記載することが出来ます.

```yaml
dev:                          # --target オプションで指定するキー
  start_time: -3600           # 取得開始時間を秒で指定（600 秒前）
  metric_name: CPUUtilization # メトリクス名を指定
  namespace: AWS/EC2          # Namespace 名を指定
  period: 300                 # Period を指定
  statistics: Maximum         # 統計を指定(Minimum, Maximum, Sum, Average, SampleCount)
  dimensions:                 # ディメンション を指定
    - name: InstanceId
      value: i-xxxxxxxx
  unit: Percent               # 単位を指定(Seconds, Bytes, Percent, Count... etc)
stg:
  start_time: -600
  metric_name: CPUUtilization
  namespace: AWS/EC2
  period: 300
  statistics: Maximum
  dimensions:
    - name: InstanceId
      value: i-xxxxxxxxxxxxx
  unit: Percent
```

* 1 つの YAML ファイルに複数の YAML 定義を記載することが可能です
* ターゲットで指定するキーは YAML ファイル内で重複しないように記述する必要があります
* YAML ファイルは任意のパスに設置が可能です

### 実行

YAML ファイルを cpu.yml というファイル名で作成し configs/sample/cpu.yml に保存している場合, 以下のように実行します.

```sh
$ ./cw -profile=you_profile -config=config/sample/cpu.yml -target=dev
```
