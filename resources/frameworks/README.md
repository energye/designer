## 框架内嵌包

 存放源码及二进制库

## 目录结构
```text
frameworks  包
└── cef
    └── cef.zip             CEF 源码
└── lcl
    └── lcl.zip             LCL 源码
└── wv
    └── darwin
        └── wv.zip          MacOS WebView 源码
    └── linux
        └── wv.zip          Linux WebView 源码 
    └── windows
        └── wv.zip          Windows WebView 源码    
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