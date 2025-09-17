package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

// 下面设计器

var (
	bottomColor        = colors.RGBToColor(100, 120, 140)
	leftBoxWidth int32 = 250
)

type BottomBox struct {
	box      lcl.IPanel
	leftBox  lcl.IPanel    // 左侧-面板组件对象查看器
	splitter lcl.ISplitter // 分割线
	rightBox lcl.IPanel    // 右侧-窗体设计器
}

func (m *TAppWindow) createBottomBox() *BottomBox {
	box := &BottomBox{}
	box.box = lcl.NewPanel(m)
	box.box.SetParent(m)
	box.box.SetBevelOuter(types.BvNone)
	box.box.SetDoubleBuffered(true)
	box.box.SetTop(toolbarHeight)
	box.box.SetWidth(m.Width())
	box.box.SetHeight(m.Height() - toolbarHeight)
	box.box.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight, types.AkBottom))
	//box.box.SetColor(bottomColor)
	m.box = box

	// 工具栏-分隔线
	box.splitter = lcl.NewSplitter(box.box)
	box.splitter.SetParent(box.box)
	box.splitter.SetAlign(types.AlLeft)
	box.splitter.SetWidth(5)
	box.splitter.SetMinSize(50)

	// 左侧-面板组件对象查看器
	box.leftBox = lcl.NewPanel(box.box)
	box.leftBox.SetParent(box.box)
	box.leftBox.SetBevelOuter(types.BvNone)
	box.leftBox.SetDoubleBuffered(true)
	box.leftBox.SetWidth(leftBoxWidth)
	box.leftBox.SetHeight(box.box.Height())
	box.leftBox.Constraints().SetMinWidth(50)
	box.leftBox.SetAlign(types.AlLeft)
	box.leftBox.SetColor(colors.ClAliceblue)

	// 右侧-窗体设计器
	box.rightBox = lcl.NewPanel(box.box)
	box.rightBox.SetParent(box.box)
	box.rightBox.SetBevelOuter(types.BvNone)
	box.rightBox.SetDoubleBuffered(true)
	box.rightBox.SetAlign(types.AlClient)
	box.rightBox.SetColor(colors.Cl3DLight)

	// 测试属性
	//btn := lcl.NewButton(m)
	//btn.SetParent(m)
	//testProp := lcl.GetComponentProperties(btn)
	//for _, prop := range testProp {
	//	fmt.Printf("%+v\n", prop)
	//}
	//
	//font := lcl.NewFont()
	//fmt.Println("font.GetNamePath():", font.GetNamePath())
	//fs := lcl.Screen.Fonts()
	//for i := 0; i < int(fs.Count()); i++ {
	//	fmt.Println(fs.Strings(int32(i)))
	//}
	return box
}
