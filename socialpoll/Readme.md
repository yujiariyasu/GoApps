* 実行

- nsqlookupdを起動し、nsqdインスタンスを発見できるようにする

```
nsqlookupd
```

- nsqdを起動し、どのnsqlookupdを利用するか指定する

```
nsqd --lookupd-tcp-address=localhost:4160
```

- mongodを起動してデータ関連のサービスを実行する(データの保管先には好きな場所を指定)

```
mongod --dbpath ./db
```

学習中につき、整理のためにコメント手厚く書くことにします