# sbbs_b

简单的 bbs 系统

## 系统功能

基本功能

1. 用户登录和注册
2. 用户发帖
3. 用户评论

注意事项

1. 用户登录时的密码不能保存为明文
2. 用户登录后，对于自己的帖子可以重新编辑或删除，不能操作其他用户的帖子
3. 用户表、文章表、评论表

修改

1. 图片验证码
2. 上传图片
3. 阻止用户在发文章或评论时输入 html 或 js 的内容
4. 防止 sql 注入

## 接口

1. 用户登录和注册接口
2. 文章列表接口
3. 评论列表接口
4. 发布文章接口
5. 评论接口
6. 上传图片接口