需求：针对编辑器功能，在保证现有功能完整情况下，优化代码逻辑，去除冗余代码，包括Go端和Web端
将Web端编辑器功能拆分为 JS和html两部分，并优化代码逻辑，保证原有功能不丢失把代码改的更容易理解和解读
将Go端编辑器功能抽出通用接口扩展能力为适配其它模式的编辑器插件加进来，同时整理和规范代码功能和代码文件
我们的目的去除啰嗦的代码逻辑让更通俗易懂，前提是不能影响现有功能。

## 目录说明

编辑器Go代码实现：C:\app\workspace\designer\designer\editor

前端编辑器实现使用的 monaco：C:\app\workspace\designer\resources\editor\assets

主要的编辑器所有功能目录：C:\app\workspace\designer\designer

gopls的Go代码，结合monaco编辑器：C:\app\workspace\designer\designer\editor\gopls 

Go依赖模块在这个工作空间 C:\app\workspace 不需要在go env 的环境 go mod 目录去找
C:\app\workspace\energy
C:\app\workspace\lcl
C:\app\workspace\wv