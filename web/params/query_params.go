package params

import (
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/strs/strcase"
	"github.com/mlogclub/simple/sqls"
	"github.com/tidwall/gjson"
	"strings"
)

type QueryParams struct {
	Ctx iris.Context
	sqls.Cnd
	gJsonResult *gjson.Result
}

func NewQueryParams(ctx iris.Context) *QueryParams {
	return &QueryParams{
		Ctx: ctx,
	}
}

func (q *QueryParams) getValueByColumn(column string) string {
	if q.Ctx.GetContentTypeRequested() == "application/json" {
		if q.gJsonResult == nil {
			body, err := q.Ctx.GetBody()
			if err != nil {
				panic(err)
			}
			result := gjson.ParseBytes(body)
			q.gJsonResult = &result
		}
		if q.gJsonResult.Get(strcase.ToLowerCamel(column)).Exists() {
			return q.gJsonResult.Get(strcase.ToLowerCamel(column)).String()
		}
	}

	if q.Ctx == nil {
		return ""
	}
	fieldName := strcase.ToLowerCamel(column)
	return q.Ctx.FormValue(fieldName)
}

func (q *QueryParams) EqByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.Eq(column, value)
	}
	return q
}

func (q *QueryParams) NotEqByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.NotEq(column, value)
	}
	return q
}

func (q *QueryParams) GtByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.Gt(column, value)
	}
	return q
}

func (q *QueryParams) GteByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.Gte(column, value)
	}
	return q
}

func (q *QueryParams) LtByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.Lt(column, value)
	}
	return q
}

func (q *QueryParams) LteByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.Lte(column, value)
	}
	return q
}

func (q *QueryParams) LikeByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	if len(value) > 0 {
		q.Like(column, value)
	}
	return q
}

func (q *QueryParams) PageByReq() *QueryParams {
	if q.Ctx == nil {
		return q
	}
	paging := GetPaging(q.Ctx)
	q.Page(paging.Page, paging.Limit)
	return q
}

func (q *QueryParams) Asc(column string) *QueryParams {
	q.Orders = append(q.Orders, sqls.OrderByCol{Column: column, Asc: true})
	return q
}

func (q *QueryParams) Desc(column string) *QueryParams {
	q.Orders = append(q.Orders, sqls.OrderByCol{Column: column, Asc: false})
	return q
}

func (q *QueryParams) Limit(limit int) *QueryParams {
	q.Page(1, limit)
	return q
}

func (q *QueryParams) Page(page, limit int) *QueryParams {
	if q.Paging == nil {
		q.Paging = &sqls.Paging{Page: page, Limit: limit}
	} else {
		q.Paging.Page = page
		q.Paging.Limit = limit
	}
	return q
}

func (q *QueryParams) BetweenByReq(column string) *QueryParams {
	value := q.getValueByColumn(column)
	result := gjson.Parse(value).Array()
	if len(result) > 1 {
		q.Cnd.Where(column+" BETWEEN ? AND ?", result[0].String(), result[1].String())
	}
	return q
}

func (q *QueryParams) SortByReq() *QueryParams {
	//q.Orders = append(q.Orders, sqls.OrderByCol{Column: column, Asc: false})
	params := q.Ctx.URLParams()
	for k, p := range params {
		if strings.HasPrefix(k, "sort[") && strings.HasSuffix(k, "]") {
			name := k[5 : len(k)-1]
			if p == "asc" || p == "ascend" {
				q.Asc(name)
			} else {
				q.Desc(name)
			}
		}
	}

	return q
}
