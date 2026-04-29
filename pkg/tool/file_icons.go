package tool

import (
	"bytes"
	"errors"
	"github.com/energye/designer/pkg/svg/oksvg"
	"github.com/energye/designer/pkg/svg/rasterx"
	"image"
	"image/png"
)

// FileExtensionToIcon 文件扩展名到实际 SVG 文件的映射
var FileExtensionToIcon = map[string]string{
	"png":      "image.svg",
	"jpeg":     "image.svg",
	"jpg":      "image.svg",
	"ico":      "image.svg",
	"tif":      "image.svg",
	"tiff":     "image.svg",
	"psd":      "image.svg",
	"psb":      "image.svg",
	"ami":      "image.svg",
	"apx":      "image.svg",
	"avif":     "image.svg",
	"bmp":      "image.svg",
	"bpg":      "image.svg",
	"brk":      "image.svg",
	"cur":      "image.svg",
	"dds":      "image.svg",
	"dng":      "image.svg",
	"exr":      "image.svg",
	"fpx":      "image.svg",
	"gbr":      "image.svg",
	"img":      "image.svg",
	"jbig2":    "image.svg",
	"jb2":      "image.svg",
	"jng":      "image.svg",
	"jxr":      "image.svg",
	"pgf":      "image.svg",
	"pic":      "image.svg",
	"raw":      "image.svg",
	"webp":     "image.svg",
	"eps":      "image.svg",
	"afphoto":  "image.svg",
	"ase":      "image.svg",
	"aseprite": "image.svg",
	"clip":     "image.svg",
	"cpt":      "image.svg",
	"heif":     "image.svg",
	"heic":     "image.svg",
	"kra":      "image.svg",
	"mdp":      "image.svg",
	"ora":      "image.svg",
	"pdn":      "image.svg",
	"reb":      "image.svg",
	"sai":      "image.svg",
	"tga":      "image.svg",
	"xcf":      "image.svg",
	"jfif":     "image.svg",
	"ppm":      "image.svg",
	"pbm":      "image.svg",
	"pgm":      "image.svg",
	"pnm":      "image.svg",
	"svgx":     "image.svg",

	"svg": "svg.svg",
	"gif": "gif.svg",

	// 视频格式
	"webm": "video.svg",
	"mkv":  "video.svg",
	"flv":  "video.svg",
	"vob":  "video.svg",
	"ogv":  "video.svg",
	"ogg":  "video.svg",
	"gifv": "video.svg",
	"avi":  "video.svg",
	"mov":  "video.svg",
	"qt":   "video.svg",
	"wmv":  "video.svg",
	"yuv":  "video.svg",
	"rm":   "video.svg",
	"rmvb": "video.svg",
	"mp4":  "video.svg",
	"m4v":  "video.svg",
	"mpg":  "video.svg",
	"mp2":  "video.svg",
	"mpeg": "video.svg",
	"mpe":  "video.svg",
	"mpv":  "video.svg",
	"m2v":  "video.svg",

	// 音频格式
	"mp3":  "audio.svg",
	"flac": "audio.svg",
	"m4a":  "audio.svg",
	"wma":  "audio.svg",
	"aiff": "audio.svg",
	"wav":  "audio.svg",

	// 文档和文本
	"pdf":   "pdf.svg",
	"md":    "markdown.svg",
	"txt":   "text.svg",
	"csv":   "csv.svg",
	"xlsx":  "csv.svg",
	"xlsm":  "csv.svg",
	"xls":   "csv.svg",
	"tsv":   "csv.svg",
	"psv":   "csv.svg",
	"ods":   "csv.svg",
	"ipynb": "notebook.svg",

	// Web
	"html":   "code-orange.svg",
	"htm":    "code-orange.svg",
	"shtml":  "code-orange.svg",
	"css":    "brackets-purple.svg",
	"scss":   "sass.svg",
	"sass":   "sass.svg",
	"less":   "brackets-sky.svg",
	"styl":   "stylus.svg",
	"pcss":   "postcss.svg",
	"sss":    "postcss.svg",
	"jsx":    "react.svg",
	"tsx":    "react-ts.svg",
	"vue":    "vue.svg",
	"svelte": "svelte.svg",
	"astro":  "astro.svg",
	"mdx":    "mdx.svg",
	"svx":    "svx.svg",

	// JavaScript/TypeScript
	"js":    "js.svg",
	"mjs":   "js.svg",
	"cjs":   "js.svg",
	"ts":    "ts.svg",
	"d.ts":  "ts-types.svg",
	"d.cts": "ts-types.svg",
	"d.mts": "ts-types.svg",

	// Python
	"py":     "python.svg",
	"python": "python.svg",

	// Go
	"go":      "go.svg",
	"go.mod":  "go-pink.svg", // 注意: 使用 go-pink.svg
	"go.sum":  "go-pink.svg",
	"go.work": "go-pink.svg",

	// Java/Kotlin
	"java": "java.svg",
	"jsp":  "java.svg",
	"kt":   "kotlin.svg",
	"kts":  "kotlin.svg",

	// C/C++
	"c":   "c.svg",
	"h":   "h.svg",
	"i":   "c.svg",
	"mi":  "c.svg",
	"cc":  "cplus.svg",
	"cpp": "cplus.svg",
	"cxx": "cplus.svg",
	"c++": "cplus.svg",
	"cp":  "cplus.svg",
	"mm":  "cplus.svg",
	"mii": "cplus.svg",
	"ii":  "cplus.svg",

	// C#
	"cs":     "csharp.svg",
	"cshtml": "razor.svg",
	"csx":    "csharp.svg",

	// Ruby
	"rb":  "ruby.svg",
	"erb": "ruby.svg",

	// PHP
	"php":       "php.svg",
	"blade.php": "laravel.svg",

	// Rust
	"rs":  "rust.svg",
	"ron": "rust.svg",

	// Swift
	"swift": "swift.svg",

	// Dart
	"dart":         "dart.svg",
	"freezed.dart": "dart.svg",
	"g.dart":       "dart.svg",

	// Haskell
	"hs":  "haskell.svg",
	"lhs": "haskell.svg",

	// Elixir
	"ex":   "elixir.svg",
	"exs":  "elixir.svg",
	"eex":  "elixir.svg",
	"leex": "elixir.svg",
	"heex": "elixir.svg",

	// Erlang
	"erl": "erlang.svg",
	"hrl": "erlang.svg",

	// Clojure
	"clj":  "clojure.svg",
	"cljs": "clojure.svg",
	"cljc": "clojure.svg",

	// Scala
	"scala": "scala.svg",
	"sc":    "scala.svg",
	"sbt":   "sbt.svg",

	// R
	"r":   "r.svg",
	"rmd": "r.svg",

	// Lua
	"lua":  "lua.svg",
	"luau": "luau.svg",

	// Perl
	"pl": "perl.svg",
	"pm": "perl.svg",

	// Julia
	"jl": "julia.svg",

	// Fortran
	"f90": "fortran.svg",
	"f95": "fortran.svg",
	"f03": "fortran.svg",
	"f":   "fortran.svg",
	"for": "fortran.svg",

	// OCaml
	"ml":  "ocaml.svg",
	"mli": "ocaml.svg",
	"cmx": "ocaml.svg",

	// F#
	"fs":     "fsharp.svg",
	"fsx":    "fsharp.svg",
	"fsi":    "fsharp.svg",
	"fsproj": "fsharp.svg",

	// Crystal
	"cr":  "crystal.svg",
	"ecr": "crystal.svg",

	// Zig
	"zig": "zig.svg",

	// Nim
	"nim": "nim.svg",

	// V
	"v": "v.svg",

	// Nix
	"nix": "nix.svg",

	// Func
	"fc": "func.svg",

	// CUDA
	"cu":  "cuda.svg",
	"cuh": "cuda.svg",

	// Shell
	"sh":         "shell.svg",
	"bash":       "shell.svg",
	"zsh":        "shell.svg",
	"fish":       "shell.svg",
	"ksh":        "shell.svg",
	"csh":        "shell.svg",
	"tcsh":       "shell.svg",
	"bat":        "shell.svg",
	"cmd":        "shell.svg",
	"ps1":        "shell.svg",
	"psm1":       "shell.svg",
	"psd1":       "shell.svg",
	"nu":         "shell.svg",
	"awk":        "shell.svg",
	"exp":        "shell.svg",
	"ssh_config": "shell.svg",

	// 配置文件
	"json": "brackets-yellow.svg",
	"yaml": "yaml.svg",
	"yml":  "yaml.svg",
	"toml": "gear.svg",
	"env":  "gear.svg",

	// Git
	"gitignore":     "git.svg",
	"gitconfig":     "git.svg",
	"gitattributes": "git.svg",
	"gitmodules":    "git.svg",
	"gitkeep":       "git.svg",

	// Docker
	"dockerfile":      "docker.svg",
	"dockerignore":    "docker.svg",
	"containerignore": "docker.svg",

	// XML
	"xml":      "xml.svg",
	"plist":    "xml.svg",
	"xsd":      "xml.svg",
	"dtd":      "xml.svg",
	"xsl":      "xml.svg",
	"xslt":     "xml.svg",
	"resx":     "xml.svg",
	"iml":      "xml.svg",
	"xquery":   "xml.svg",
	"manifest": "xml.svg",

	// 数据库
	"sql":      "database.svg",
	"db":       "database.svg",
	"sqlite":   "database.svg",
	"sqlite3":  "database.svg",
	"pgsql":    "database.svg",
	"postgres": "database.svg",
	"psql":     "database.svg",
	"pdb":      "database.svg",

	// GraphQL
	"graphql": "graphql.svg",
	"gql":     "graphql.svg",

	// Terraform
	"tf":      "terraform.svg",
	"tf.json": "terraform.svg",
	"tfvars":  "terraform.svg",
	"tfstate": "terraform.svg",

	// Font
	"woff":  "font.svg",
	"woff2": "font.svg",
	"ttf":   "font.svg",
	"eot":   "font.svg",
	"otf":   "font.svg",
	"fnt":   "font.svg",

	// Template
	"pug":    "pug.svg",
	"jade":   "pug.svg",
	"haml":   "haml.svg",
	"twig":   "twig.svg",
	"liquid": "liquid.svg",
	"njk":    "nunjucks.svg",

	// Angular
	"component.ts": "angular-component.svg",
	"component.js": "angular-component.svg",
	"service.ts":   "angular-service.svg",
	"directive.ts": "angular-directive.svg",
	"module.ts":    "angular-module.svg",
	"guard.ts":     "angular-guard.svg",
	"pipe.ts":      "angular-pipe.svg",

	// NestJS
	"controller.ts": "nest-controller.svg",
	"decorator.ts":  "nest-decorator.svg",
	"middleware.ts": "nest-middleware.svg",

	// Redux
	"actions.ts":  "redux-actions.svg",
	"effects.ts":  "redux-effects.svg",
	"facade.ts":   "redux-facade.svg",
	"reducer.ts":  "redux-reducer.svg",
	"selector.ts": "redux-selector.svg",

	// Test
	"test.js":  "js-test.svg",
	"spec.js":  "js-test.svg",
	"test.ts":  "ts-test.svg",
	"spec.ts":  "ts-test.svg",
	"test.jsx": "react-test.svg",
	"spec.jsx": "react-test.svg",
	"test.tsx": "react-test.svg",
	"spec.tsx": "react-test.svg",

	// Lock
	"lock": "lock.svg",

	// License
	"license":     "license.svg",
	"license.md":  "license.svg",
	"license.txt": "license.svg",

	// HTTP
	"http": "http.svg",
	"rest": "http.svg",
	"bru":  "http.svg",

	// Protocol Buffer
	"proto": "proto.svg",

	// LaTeX
	"tex": "tex.svg",
	"sty": "tex.svg",
	"dtx": "tex.svg",

	// Draw.io
	"drawio": "drawio.svg",
	"dio":    "drawio.svg",

	// Patch
	"patch": "patch.svg",

	// Minecraft
	"mcfunction":  "minecraft.svg",
	"mcmeta":      "minecraft.svg",
	"mcworld":     "minecraft.svg",
	"mcstructure": "minecraft.svg",
	"mcpack":      "minecraft.svg",
	"mcaddon":     "minecraft.svg",

	// i18n
	"lang": "i18n.svg",
	"mo":   "i18n.svg",
	"po":   "i18n.svg",
	"pot":  "i18n.svg",

	// ColdFusion
	"cfml": "coldfusion.svg",
	"cfc":  "coldfusion.svg",
	"cfm":  "coldfusion.svg",

	// Rescript
	"res":  "rescript.svg",
	"resi": "rescript-interface.svg",

	// Solidity
	"sol": "solidity.svg",

	// Markdoc
	"mdoc": "markdoc.svg",

	// CoffeeScript
	"coffee": "coffeescript.svg",

	// Gradle
	"gradle": "gradle.svg",

	// CMake
	"cmake": "cmake.svg",

	// Gulp
	"gulpfile.js": "gulp.svg",
	"gulpfile.ts": "gulp.svg",

	// ESLint
	"eslintrc": "eslint.svg",

	// Prettier
	"prettierrc": "prettier.svg",

	// Jest
	"jest.config": "jest.svg",

	// Babel
	"babelrc": "babel.svg",

	// Webpack
	"webpack.config": "webpack.svg",

	// Vite
	"vite.config": "vite.svg",

	// Vitest
	"vitest.config": "vitest.svg",

	// Storybook
	"stories": "storybook.svg",
	"story":   "storybook.svg",

	// Cypress
	"cypress": "cypress.svg",

	// Jenkins
	"jenkinsfile": "jenkins.svg",

	// Firebase
	"firebase.json": "firebase.svg",
	"firebaserc":    "firebase.svg",

	// Prisma
	"prisma": "prisma.svg",

	// Tailwind
	"tailwind.config": "tailwind.svg",

	// PostCSS
	"postcss.config": "postcss.svg",

	// Stylelint
	"stylelintrc": "stylelint.svg",

	// Node
	"package.json": "node.svg",

	// NPM
	"npmignore": "npm.svg",
	"npmrc":     "npm.svg",

	// Yarn
	"yarn.lock": "yarn.svg",
	"yarnrc":    "yarn.svg",

	// PNPM
	"pnpm-lock.yaml":      "pnpm.svg",
	"pnpm-workspace.yaml": "pnpm.svg",

	// Bun
	"bun.lockb":   "bun.svg",
	"bunfig.toml": "bun.svg",

	// Deno
	"deno.json":  "deno.svg",
	"deno.jsonc": "deno.svg",

	// GitLab
	".gitlab-ci.yml": "gitlab.svg",

	// GitHub
	"github": "github.svg",

	// VSCode
	"vscodeignore": "ignore.svg",
	"vsixmanifest": "puzzle.svg",
	"vsix":         "puzzle.svg",

	// Cursor
	"cursorrules": "cursor.svg",

	// Claude
	"CLAUDE.md":   "claude.svg",
	".claude":     "claude.svg",
	"claude.json": "claude.svg",

	// Tauri
	"tauri.conf.json": "tauri.svg",

	// Vercel
	"vercel.json": "vercel.svg",

	// Netlify
	"netlify.json": "netlify.svg",
	"netlify.yml":  "netlify.svg",
	"netlify.toml": "netlify.svg",

	// Next.js
	"next.config.js":  "next.svg",
	"next.config.mjs": "next.svg",

	// Nuxt
	"nuxt.config.js": "nuxt.svg",
	"nuxtrc":         "nuxt.svg",

	// Gatsby
	"gatsby-config.js":  "gatsby.svg",
	"gatsby-node.js":    "gatsby.svg",
	"gatsby-browser.js": "gatsby.svg",

	// Hugo
	"hugo": "hugo.svg",

	// Docusaurus
	"docusaurus.config.js": "docusaurus.svg",

	// Fresh
	"fresh.config.ts": "fresh.svg",

	// SvelteKit
	"svelte.config.js": "svelte.svg",

	// Astro
	"astro.config.js": "astro.svg",

	// Contentlayer
	"contentlayer.config.ts": "contentlayer.svg",

	// Drizzle
	"drizzle.config.ts": "drizzle.svg",

	// Sanity
	"sanity.cli.ts":    "sanity.svg",
	"sanity.config.ts": "sanity.svg",

	// Panda CSS
	"panda.config.ts": "panda.svg",

	// UnoCSS
	"uno.config.js":    "unocss.svg",
	"unocss.config.js": "unocss.svg",

	// Vanilla Extract
	"css.ts": "vanilla-extract.svg",

	// Shadcn
	"components.json": "shadcn.svg",

	// YummaCSS
	"yumma.css":    "yummacss.svg",
	"yummacss.css": "yummacss.svg",

	// Knip
	"knip.json":   "knip.svg",
	"knip.config": "knip.svg",

	// Orval
	"orval.config.js": "orval.svg",
	"orvalrc":         "orval.svg",

	// Lunaria
	"lunaria.config.json": "lunaria.svg",

	// Expressive Code
	"ec.config.mjs": "expressive-code.svg",

	// Rome
	"rome.json": "rome.svg",

	// Biome
	"biome.json":  "biome.svg",
	"biome.jsonc": "biome.svg",

	// Oxlint
	"oxlintrc.json": "oxlint.svg",

	// SWC
	"swcrc": "swc.svg",

	// RSBuild
	"rsbuild.config.ts": "rsbuild.svg",

	// RSPack
	"rspack.config.ts": "rspack.svg",

	// RSLib
	"rslib.config.ts": "rslib.svg",

	// Turborepo
	"turbo.json": "turborepo.svg",

	// Nx
	"nx.json":  "nx.svg",
	"nxignore": "nx.svg",

	// Pulumi
	"pulumi.yaml": "pulumi.svg",

	// Serverless
	"serverless.yml": "severless.svg",

	// Cloudflare Workers
	"dev.vars": "cloudflare-workers.svg",

	// Earthly
	"earthlyignore": "earthfile.svg",

	// Bruno
	"bruno.json": "bruno.svg",

	// Capacitor
	"capacitor.config.json": "capacitor.svg",

	// Ionic
	"ionic.config.json": "ionic.svg",

	// Dune
	"dune":         "dune.svg",
	"dune-project": "dune.svg",

	// Keystatic
	"keystatic.config.ts": "keystatic.svg",
	"keystatic.page.ts":   "keystatic.svg",

	// Statamic Antlers
	"antlers.html": "statamic-antlers.svg",

	// Pkl
	"pkl": "pkl.svg",

	// Gleam
	"gleam": "gleam.svg",

	// Visual Studio
	"csproj": "visual-studio.svg",
	"sln":    "visual-studio.svg",
	"slnx":   "visual-studio.svg",

	// Robot
	"robot": "robot.svg",

	// Puppet
	"pp": "puppet.svg",

	// Nodemon
	"nodemon.json": "nodemon.svg",

	// Exe
	"exe": "exe.svg",
	"msi": "exe.svg",

	// Docker Compose (特殊颜色)
	"docker-compose.yml": "docker-pink.svg",
	"compose.yml":        "docker-pink.svg",

	"mongo":   "folder-mongo",
	"mongodb": "folder-mongo",
}

// GetIconForFile 根据文件名获取对应的 SVG 图标
func GetIconForFile(filename string) string {
	if icon, ok := FileExtensionToIcon[filename]; ok {
		return icon
	}
	for ext, icon := range FileExtensionToIcon {
		if len(ext) > 0 && ext[0] == '.' {
			if len(filename) > len(ext) && filename[len(filename)-len(ext):] == ext {
				return icon
			}
		} else if !contains(ext, ".") {
			if len(filename) > len(ext)+1 && filename[len(filename)-len(ext)-1] == '.' {
				if filename[len(filename)-len(ext):] == ext {
					return icon
				}
			}
		}
	}
	return "document.svg"
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s != substr
}

func SVGToPNG(svgData []byte, targetWidth, targetHeight int) ([]byte, error) {
	if svgData == nil {
		return nil, errors.New("svgData is nil")
	}
	icon, err := oksvg.ReadIconStream(bytes.NewBuffer(svgData))
	if err != nil {
		return nil, err
	}
	w, h := float64(targetWidth), float64(targetHeight)
	icon.SetTarget(0, 0, w, h)
	rgba := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	scannerGV := rasterx.NewScannerGV(targetWidth, targetHeight, rgba, rgba.Bounds())
	rasterizer := rasterx.NewDasher(targetWidth, targetHeight, scannerGV)
	icon.Draw(rasterizer, 1.0)
	pngBuf := bytes.NewBuffer(nil)
	err = png.Encode(pngBuf, rgba)
	return pngBuf.Bytes(), err
}
