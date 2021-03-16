package ctxtags

import "github.com/gin-gonic/gin"

func CtxTags(opts ...Option) gin.HandlerFunc {
	return func(c *gin.Context) {
		o := evaluateOptions(opts...)
		c = ToContext(c, NewTags())
		if o.requestFieldsFunc != nil {
			setRequestFieldTags(c, o.requestFieldsFunc)
		}
		c.Next()
	}
}

func setRequestFieldTags(c *gin.Context, f RequestFieldExtractorFunc) {
	if fields := f(c); fields != nil {
		t := FromContext(c)
		for k, v := range fields {
			t.Set(k, v)
		}
	}
}
