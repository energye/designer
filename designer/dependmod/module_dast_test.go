package dependmod

import (
	"github.com/energye/designer/pkg/dast"
	"github.com/energye/designer/resources/frameworks"
	"path/filepath"
	"testing"
)

func TestInitDependencyModule(t *testing.T) {
	data, err := frameworks.LCL("lcl/callback_event_def.go")
	t.Log(len(data) > 0, err)
	lclSRCEventDef := filepath.Join("lcl", "callback_event_def.go")
	GLCLFuncTypeAliases := dast.GetAllFuncTypeAliasesByCode(lclSRCEventDef, data)
	t.Log(GLCLFuncTypeAliases != nil)
}
