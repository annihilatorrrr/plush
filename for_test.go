package plush_test

import (
	"strings"
	"testing"

	"github.com/gobuffalo/plush/v5"
	"github.com/stretchr/testify/require"
)

func Test_Render_For_Array(t *testing.T) {
	r := require.New(t)
	input := `<% for (i,v) in ["a", "b", "c"] {return v} %>`
	s, err := plush.Render(input, plush.NewContext())
	r.NoError(err)
	r.Equal("", s)
}

func Test_Render_For_Update_Global_Scope(t *testing.T) {
	r := require.New(t)
	input := `<% let varTest = "" %><% for (i,v) in ["a", "b", "c"] {varTest =  v} %><%= varTest %>`
	s, err := plush.Render(input, plush.NewContext())
	r.NoError(err)
	r.Equal("c", s)
}

func Test_Render_For_Hash(t *testing.T) {
	r := require.New(t)
	input := `<%= for (k,v) in myMap { %><%= k + ":" + v%><% } %>`
	s, err := plush.Render(input, plush.NewContextWith(map[string]interface{}{
		"myMap": map[string]string{
			"a": "A",
			"b": "B",
		},
	}))
	r.NoError(err)
	r.Contains(s, "a:A")
	r.Contains(s, "b:B")
}

func Test_Render_For_Global_Scope_Key_Access(t *testing.T) {
	r := require.New(t)
	input := `<%= for (k,v) in myMap { %><%= k + ":" + v%><% } %> <%= k %>`
	_, err := plush.Render(input, plush.NewContextWith(map[string]interface{}{
		"myMap": map[string]string{
			"a": "A",
		},
	}))
	r.Error(err)
	r.Errorf(err, `line 1: "k": unknown identifier`)
}

func Test_Render_For_Global_Scope_Value_Access(t *testing.T) {
	r := require.New(t)
	input := `<%= for (k,v) in myMap { %><%= k + ":" + v%><% } %> <%= v %>`
	_, err := plush.Render(input, plush.NewContextWith(map[string]interface{}{
		"myMap": map[string]string{
			"a": "A",
		},
	}))
	r.Error(err)
	r.Errorf(err, `line 1: "v": unknown identifier`)
}

func Test_Render_For_Nested_For_With_Same_Iterators_Keys(t *testing.T) {
	r := require.New(t)
	input := `<%= for (k,v) in myMap { %><%=  for (k,v) in myMap2 { %><%= k + ":" + v%><% } %>%><%}%>`
	s, err := plush.Render(input, plush.NewContextWith(map[string]interface{}{
		"myMap": map[string]string{
			"a": "A",
		},
		"myMap2": map[string]string{
			"b": "B",
		},
	}))
	r.NoError(err)
	r.Contains(s, "b:B")
}

func Test_Render_For_Array_Return(t *testing.T) {
	r := require.New(t)
	input := `<%= for (i,v) in ["a", "b", "c"] {return v} %>`
	s, err := plush.Render(input, plush.NewContext())
	r.NoError(err)
	r.Equal("abc", s)
}

func Test_Render_For_Array_Continue(t *testing.T) {
	r := require.New(t)
	input := `<%= for (i,v) in [1, 2, 3,4,5,6,7,8,9,10] {
		%>Start<%
		if (v == 1 || v ==3 || v == 5 || v == 7 || v == 9) {


			%>Odd<%
			continue
		}

		return v
		} %>`
	s, err := plush.Render(input, plush.NewContext())

	r.NoError(err)
	r.Equal("StartOddStart2StartOddStart4StartOddStart6StartOddStart8StartOddStart10", s)
}

func Test_Render_For_Array_WithNoOutput(t *testing.T) {
	r := require.New(t)
	input := `<%= for (i,v) in [1, 2, 3,4,5,6,7,8,9,10] {

		if (v == 1 || v == 2 || v ==3 || v == 4|| v == 5 || v == 6 || v == 7 || v == 8 || v == 9 || v == 10) {

			continue
		}

		return v
		} %>`
	s, err := plush.Render(input, plush.NewContext())

	r.NoError(err)
	r.Equal("", s)
}

func Test_Render_For_Array_WithoutContinue(t *testing.T) {
	r := require.New(t)
	input := `<%= for (i,v) in [1, 2, 3,4,5,6,7,8,9,10] {
		if (v == 1 || v ==3 || v == 5 || v == 7 || v == 9) {
		}
		return v
		} %>`
	s, err := plush.Render(input, plush.NewContext())

	r.NoError(err)
	r.Equal("12345678910", s)
}

func Test_Render_For_Array_ContinueNoControl(t *testing.T) {
	r := require.New(t)
	input := `<%= for (i,v) in [1, 2, 3,4,5,6,7,8,9,10] {
		continue
		return v
		} %>`
	s, err := plush.Render(input, plush.NewContext())

	r.NoError(err)
	r.Equal("", s)
}

func Test_Render_For_Array_Break_String(t *testing.T) {
	r := require.New(t)
	input := `<%= for (i,v) in [1, 2, 3,4,5,6,7,8,9,10] {
		%>Start<%
		if (v == 5) {


			%>Odd<%
			break
		}

		return v
		} %>`
	s, err := plush.Render(input, plush.NewContext())

	r.NoError(err)
	r.Equal("Start1Start2Start3Start4StartOdd", s)
}

func Test_Render_For_Array_WithBreakFirstValue(t *testing.T) {
	r := require.New(t)
	input := `<%= for (i,v) in [1, 2, 3,4,5,6,7,8,9,10] {
		if (v == 1 || v ==3 || v == 5 || v == 7 || v == 9) {
			break
		}
		return v
		} %>`
	s, err := plush.Render(input, plush.NewContext())

	r.NoError(err)
	r.Equal("", s)
}

func Test_Render_For_Array_WithBreakFirstValueWithReturn(t *testing.T) {
	r := require.New(t)
	input := `<%= for (i,v) in [1, 2, 3,4,5,6,7,8,9,10] {
		if (v == 1 || v ==3 || v == 5 || v == 7 || v == 9) {
			%><%=v%><%
			break
		}
		return v
		} %>`
	s, err := plush.Render(input, plush.NewContext())

	r.NoError(err)
	r.Equal("1", s)
}
func Test_Render_For_Array_Break(t *testing.T) {
	r := require.New(t)
	input := `<%= for (i,v) in [1, 2, 3,4,5,6,7,8,9,10] {
		break
		return v
		} %>`
	s, err := plush.Render(input, plush.NewContext())

	r.NoError(err)
	r.Equal("", s)
}

func Test_Render_For_Array_Key_Only(t *testing.T) {
	r := require.New(t)
	input := `<%= for (v) in ["a", "b", "c"] {%><%=v%><%} %>`
	s, err := plush.Render(input, plush.NewContext())
	r.NoError(err)
	r.Equal("abc", s)
}

func Test_Render_For_Func_Range(t *testing.T) {
	r := require.New(t)
	input := `<%= for (v) in range(3,5) { %><%=v%><% } %>`
	s, err := plush.Render(input, plush.NewContext())
	r.NoError(err)
	r.Equal("345", s)
}

func Test_Render_For_Func_Between(t *testing.T) {
	r := require.New(t)
	input := `<%= for (v) in between(3,6) { %><%=v%><% } %>`
	s, err := plush.Render(input, plush.NewContext())
	r.NoError(err)
	r.Equal("45", s)
}

func Test_Render_For_Func_Until(t *testing.T) {
	r := require.New(t)
	input := `<%= for (v) in until(3) { %><%=v%><% } %>`
	s, err := plush.Render(input, plush.NewContext())
	r.NoError(err)
	r.Equal("012", s)
}

func Test_Render_For_Array_Key_Value(t *testing.T) {
	r := require.New(t)
	input := `<%= for (i,v) in ["a", "b", "c"] {%><%=i%><%=v%><%} %>`
	s, err := plush.Render(input, plush.NewContext())
	r.NoError(err)
	r.Equal("0a1b2c", s)
}

func Test_Render_For_Array_Key_Global_Scope_Same_Identifier(t *testing.T) {
	r := require.New(t)
	input := `<%   let i = 10000 %><%= for (i,v) in ["a", "b", "c"] {%><%=i%><%=v%><%} %><%= i %>`
	s, err := plush.Render(input, plush.NewContext())
	r.NoError(err)
	r.Equal("0a1b2c10000", s)
}

func Test_Render_For_Array_Key_Not_Defined(t *testing.T) {
	r := require.New(t)
	input := `<%= for (i,v) in ["a", "b", "c"] {%><%=i%><%=v%><%} %><%= i %>`
	_, err := plush.Render(input, plush.NewContext())
	r.Error(err)
	r.Errorf(err, `line 1: "i": unknown identifier`)
}

func Test_Render_For_Array_Value_Global_Scope(t *testing.T) {
	r := require.New(t)
	input := ` <%= for (i,v) in ["a", "b", "c"] {%><%=i%><%=v%><%} %><%= v %>`
	_, err := plush.Render(input, plush.NewContext())
	r.Error(err)
	r.Errorf(err, `line 1: "v": unknown identifier`)
}

func Test_Render_For_Nil(t *testing.T) {
	r := require.New(t)
	input := `<% for (i,v) in nilValue {return v} %>`
	ctx := plush.NewContext()
	ctx.Set("nilValue", nil)
	s, err := plush.Render(input, ctx)
	r.Error(err)
	r.Equal("", s)
}

func Test_Render_For_Map_Nil_Value(t *testing.T) {
	r := require.New(t)
	input := `
	<%= for (k, v) in flash["errors"] { %>
		Flash:
			<%= k %>:<%= v %>
	<% } %>
`
	ctx := plush.NewContext()
	ctx.Set("flash", map[string][]string{})
	s, err := plush.Render(input, ctx)
	r.NoError(err)
	r.Equal("", strings.TrimSpace(s))
}

type Category struct {
	Products []Product
}
type Product struct {
	Name []string
}

func Test_Render_For_Array_OutofBoundIndex(t *testing.T) {
	r := require.New(t)
	ctx := plush.NewContext()
	product_listing := Category{}
	ctx.Set("product_listing", product_listing)
	input := `<%= for (i, names) in product_listing.Products[0].Name { %>
				<%= splt %>
			<% } %>`
	_, err := plush.Render(input, ctx)
	r.Error(err)
}
