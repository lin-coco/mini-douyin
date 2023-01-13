## 前言

有的时候我们需要跟别人合作进行开发，然后分别使用不同的Git分支，等项目完成时，需要进行代码合并，就需要知道Git如何合并远程分支。

团队开发：

main（master分支）：不作为开发分支

dev-a：开发分支，由开发个人掌控

dev-b：开发分支，由开发个人掌控

## 步骤

当你的代码提交前或者开发前，要先进行与主分支的合并，因为可能主分支先合并了其他分支拥有了新代码，假设你本地在使用的分支为dev-a，需要合并的远程分支为main（主分支）

## 第一步

在本地新建一个与远程的分支main相同(被合并的版本)的分支main

```
git checkout -b main origin/main
```

该指令的意思：创建一个本地分支，并将远程分支放到该分支里面去。

## 第二步

将远程代码pull到本地

```
git pull origin main
```

2

## 第三步

返回到你的分支a

```
git checkout dev-a
```

1

## 第四步

合并分支a与分支b

```
git merge main
```

该指令的意思：当前所在分支与main进行合并。dev-a就拥有了主分支的代码

## 第五步

把本地的分支a同步到远程

```
git push origin dev-a
```

## 第六步

pull Request  在github上请求合并

![image-20230113143645730](https://typora-img-xue.oss-cn-beijing.aliyuncs.com/img/image-20230113143645730.png)
选择自己的分支

![image-20230113143737285](https://typora-img-xue.oss-cn-beijing.aliyuncs.com/img/image-20230113143737285.png)

![image-20230113143811885](https://typora-img-xue.oss-cn-beijing.aliyuncs.com/img/image-20230113143811885.png)

![image-20230113143834948](https://typora-img-xue.oss-cn-beijing.aliyuncs.com/img/image-20230113143834948.png)

## 第七步

（原则上项目经理）main分支 同意合并

![image-20230113143934189](https://typora-img-xue.oss-cn-beijing.aliyuncs.com/img/image-20230113143934189.png)

## 其他

如果你不需要本地或者远程的分支，你可以查询并删除多余分支。

本地 查询本地分支：

git branch 1 删除本地分支:

git branch -D br 1 远程 查询远程分支：

git branch 1 删除远程分支:

git push origin :br (origin 后面有空格) 1 修改本地分支与远程分支的关联：

git branch --set-upstream-to=origin/remote_branch your_branch 1

git clone -b [分支名] [地址] 克隆指定分支

git branch 查询本地分支





