
# meme-web-go

_✨ 运行在浏览器的表情包生成器，用于制作各种沙雕表情包 ✨_
_✨ 本程序基于 [meme-generator-rs](https://github.com/MemeCrafters/meme-generator-rs) 与  [meme-generator](https://github.com/MemeCrafters/meme-generator) 编写 ✨_

## 安装方法

1. 安装并运行任意版本 meme-generator 
   - [meme-generator-rs](https://github.com/MemeCrafters/meme-generator-rs) （rust 版本）
   - [meme-generator](https://github.com/MemeCrafters/meme-generator) （python 版本）
2. 下载对应系统的 [release](https://github.com/zouXH-god/meme-web-go/releases) 文件
3. 初次运行程序会创建 `.env` 文件：
    ```env
    HOST=127.0.0.1 
    PORT=5222
    
    # Meme Point 填写 meme 的api地址，同时兼容 python 与 rust 版本的meme，有多个则用英文逗号隔开
    MEME_POINT=http://127.0.0.1:2233,http://127.0.0.2:2233,...
    ```
4. 修改 `.env` 文件后重新运行程序

