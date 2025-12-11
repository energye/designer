## 框架内嵌包

 存放源码及二进制库

## 目录结构
```text
备注: 压缩包 xxx.zip 源码库直接包含源码, 没有 lcl, cef, wv 目录

frameworks  包
└── cef
    └── cef.zip             CEF 源码库
└── lcl
    └── lcl.zip             LCL 源码库
└── wv
    └── darwin
        └── wv.zip          MacOS WebView 源码库
    └── linux
        └── wv.zip          Linux WebView 源码库
    └── windows
        └── wv.zip          Windows WebView 源码库
└── lib
    └── darwin
        └── libenergy.zip   libenergy.dylib   动态链接库
            └── x86_64 | ARM64                架构
            └── LCL | LCL+CEF | LCL+WebView   三个运行时库
    └── linux
        └── libenergy.zip   libenergy.so      动态链接库
            └── x86_64 | ARM64                架构
            └── LCL | LCL+CEF | LCL+WebView   三个运行时库
    └── windows
        └── libenergy.zip   libenergy.dll     动态链接库 
            └── x86_64 | i386                 架构
            └── LCL | LCL+CEF | LCL+WebView   三个运行时库
    
```  