module github.com/energye/designer

go 1.20

replace (
	github.com/energye/lcl => ../lcl
	github.com/energye/widget => ../widget
	github.com/energye/energy/v3 => ../energy
)

require (
	github.com/energye/lcl v0.0.0-beta
	github.com/energye/widget v0.0.0-beta
	github.com/energye/energy/v3 v3.0.0-beta
)
