package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 对象查看器

var (
	componentTreeHeight int32 = 150
	inspector           *Inspector
)

// 组件树和对象查看器
type Inspector struct {
	componentTreeFilter     lcl.ITreeFilterEdit     // 组件树过滤框
	componentTree           lcl.ITreeView           // 组件树
	componentTreeBox        lcl.IPanel              // 组件树盒子
	componentPropertyFilter lcl.ITreeFilterEdit     // 组件属性过滤框
	componentProperty       lcl.ILazVirtualDrawTree // 组件属性
	componentPropertyBox    lcl.IPanel              // 组件属性盒子
	boxSplitter             lcl.ISplitter           // 分割线
}

func (m *BottomBox) createInspectorLayout() *Inspector {
	ins := new(Inspector)
	// 对象查看器面板分隔
	{
		ins.boxSplitter = lcl.NewSplitter(m.leftBox)
		ins.boxSplitter.SetParent(m.leftBox)
		ins.boxSplitter.SetAlign(types.AlTop)

		ins.componentTreeBox = lcl.NewPanel(m.leftBox)
		ins.componentTreeBox.SetParent(m.leftBox)
		ins.componentTreeBox.SetBevelOuter(types.BvNone)
		ins.componentTreeBox.SetDoubleBuffered(true)
		ins.componentTreeBox.SetWidth(m.leftBox.Width())
		ins.componentTreeBox.SetHeight(componentTreeHeight)
		ins.componentTreeBox.Constraints().SetMinWidth(50)
		ins.componentTreeBox.Constraints().SetMinHeight(50)
		ins.componentTreeBox.SetAlign(types.AlTop)
		//ins.componentTreeBox.SetColor(colors.ClAliceblue)

		ins.componentPropertyBox = lcl.NewPanel(m.leftBox)
		ins.componentPropertyBox.SetParent(m.leftBox)
		ins.componentPropertyBox.SetBevelOuter(types.BvNone)
		ins.componentPropertyBox.SetDoubleBuffered(true)
		ins.componentPropertyBox.SetWidth(m.leftBox.Width())
		ins.componentPropertyBox.SetHeight(componentTreeHeight)
		ins.componentPropertyBox.Constraints().SetMinWidth(50)
		ins.componentPropertyBox.Constraints().SetMinHeight(50)
		ins.componentPropertyBox.SetAlign(types.AlClient)
		//ins.componentPropertyBox.SetColor(colors.Cl3DDkShadow)
	}
	// 组件树
	{
		componentTreeTitle := lcl.NewLabel(ins.componentTreeBox)
		componentTreeTitle.SetParent(ins.componentTreeBox)
		componentTreeTitle.SetCaption("组件")
		componentTreeTitle.Font().SetStyle(types.NewSet(types.FsBold))
		componentTreeTitle.SetTop(8)
		componentTreeTitle.SetLeft(5)

		ins.componentTreeFilter = lcl.NewTreeFilterEdit(ins.componentTreeBox)
		ins.componentTreeFilter.SetParent(ins.componentTreeBox)
		ins.componentTreeFilter.SetTop(5)
		ins.componentTreeFilter.SetLeft(30)
		ins.componentTreeFilter.SetWidth(m.leftBox.Width() - ins.componentTreeFilter.Left())
		ins.componentTreeFilter.SetAlign(types.AlCustom)
		ins.componentTreeFilter.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight))

		ins.componentTree = lcl.NewTreeView(ins.componentTreeBox)
		ins.componentTree.SetParent(ins.componentTreeBox)
		ins.componentTree.SetTop(35)
		ins.componentTree.SetWidth(m.leftBox.Width())
		ins.componentTree.SetHeight(componentTreeHeight - ins.componentTree.Top())
		ins.componentTree.SetAlign(types.AlCustom)
		ins.componentTree.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkBottom, types.AkRight))
	}
	// 组件属性
	{
		componentPropertyTitle := lcl.NewLabel(ins.componentPropertyBox)
		componentPropertyTitle.SetParent(ins.componentPropertyBox)
		componentPropertyTitle.SetCaption("属性")
		componentPropertyTitle.Font().SetStyle(types.NewSet(types.FsBold))
		componentPropertyTitle.SetTop(8)
		componentPropertyTitle.SetLeft(5)

		ins.componentPropertyFilter = lcl.NewTreeFilterEdit(ins.componentPropertyBox)
		ins.componentPropertyFilter.SetParent(ins.componentPropertyBox)
		ins.componentPropertyFilter.SetTop(5)
		ins.componentPropertyFilter.SetLeft(30)
		ins.componentPropertyFilter.SetWidth(m.leftBox.Width() - ins.componentPropertyFilter.Left())
		ins.componentPropertyFilter.SetAlign(types.AlCustom)
		ins.componentPropertyFilter.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight))

		ins.componentProperty = lcl.NewLazVirtualDrawTree(ins.componentPropertyBox)
		ins.componentProperty.SetParent(ins.componentPropertyBox)
		ins.componentProperty.SetTop(35)
		ins.componentProperty.SetWidth(m.leftBox.Width())
		ins.componentProperty.SetHeight(ins.componentPropertyBox.Height() - (ins.componentTree.Top()))
		ins.componentProperty.SetAlign(types.AlCustom)
		ins.componentProperty.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkBottom, types.AkRight))
	}
	return ins
}
