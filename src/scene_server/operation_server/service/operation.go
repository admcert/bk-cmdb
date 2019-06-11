/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (o *OperationServer) CreateStatisticChart(ctx *rest.Contexts) {
	chartInfo := new(metadata.ChartConfig)
	if err := ctx.DecodeInto(chartInfo); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// 自定义报表
	if chartInfo.ReportType == common.OperationCustom {
		result, err := o.Engine.CoreAPI.CoreService().Operation().CreateOperationChart(ctx.Kit.Ctx, ctx.Kit.Header, chartInfo)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrOperationNewAddStatisticFail, "new add statistic fail, err: %v", err)
			return
		}

		blog.Debug("count: %v", result.Data)
		ctx.RespEntity(result.Data)
		return
	}

	blog.Debug("create inner chart")
	// 内置报表
	srvData := o.newSrvComm(ctx.Kit.Header)
	resp, err := srvData.lgc.CreateInnerChart(ctx.Kit, chartInfo)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationNewAddStatisticFail, "new add statistic fail, err: %v", err)
		return
	}

	ctx.RespEntity(resp)
}

func (o *OperationServer) DeleteStatisticChart(ctx *rest.Contexts) {
	opt := mapstr.MapStr{}
	if err := ctx.DecodeInto(&opt); err != nil {
		ctx.RespAutoError(err)
		return
	}
	_, err := o.Engine.CoreAPI.CoreService().Operation().DeleteOperationChart(ctx.Kit.Ctx, ctx.Kit.Header, opt)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationDeleteStatisticFail, "search chart info fail, err: %v, id: %v", err)
		return
	}

	ctx.RespEntity(nil)
}

func (o *OperationServer) SearchStatisticChart(ctx *rest.Contexts) {
	opt := make(map[string]interface{})

	blog.Debug("here")
	result, err := o.Engine.CoreAPI.CoreService().Operation().SearchOperationChart(ctx.Kit.Ctx, ctx.Kit.Header, opt)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationSearchStatisticsFail, "search chart info fail, err: %v", err)
		return
	}

	blog.Debug("result: %v", result)
	ctx.RespEntity(result.Data)
}

func (o *OperationServer) UpdateStatisticChart(ctx *rest.Contexts) {
	opt := mapstr.MapStr{}
	if err := ctx.DecodeInto(&opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := o.Engine.CoreAPI.CoreService().Operation().UpdateOperationChart(ctx.Kit.Ctx, ctx.Kit.Header, opt)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationSearchStatisticsFail, "update statistic info fail, err: %v", err)
		return
	}

	ctx.RespEntity(result.Data)
}

func (o *OperationServer) SearchChartData(ctx *rest.Contexts) {
	inputParams := mapstr.MapStr{}
	if err := ctx.DecodeInto(&inputParams); err != nil {
		ctx.RespAutoError(err)
		return
	}

	chart, err := o.CoreAPI.CoreService().Operation().SearchChartByID(ctx.Kit.Ctx, ctx.Kit.Header, inputParams)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationGetChartDataFail, "search chart data fail, err: %v", err)
		return
	}

	// 判断模型是否存在，不存在返回nil
	cond := make(map[string]interface{}, 0)
	cond[common.BKObjIDField] = chart.Data.ObjID
	query := metadata.QueryCondition{Condition: cond}
	models, err := o.CoreAPI.CoreService().Model().ReadModel(ctx.Kit.Ctx, ctx.Kit.Header, &query)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationGetChartDataFail, "search chart data fail, err: %v", err)
		return
	}
	if models.Data.Count <= 0 {
		ctx.RespEntity(nil)
		return
	}

	innerChart := []string{
		"host_change_biz_chart", "model_inst_chart", "model_inst_change_chart",
		"biz_module_host_chart", "model_and_inst_count",
	}

	srvData := o.newSrvComm(ctx.Kit.Header)
	if !util.InStrArr(innerChart, chart.Data.ReportType) {
		result, err := srvData.lgc.CommonStatisticFunc(ctx.Kit, chart.Data)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrOperationGetChartDataFail, "search chart data fail, err: %v", err)
			return
		}
		ctx.RespEntity(result)
		return
	}

	result, err := srvData.lgc.GetInnerChartData(ctx.Kit, chart.Data)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationGetChartDataFail, "search chart data fail, err: %v", err)
		return
	}

	ctx.RespEntity(result)
}

func (o *OperationServer) UpdateChartPosition(ctx *rest.Contexts) {
	opt := mapstr.MapStr{}
	if err := ctx.DecodeInto(&opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	_, err := o.CoreAPI.CoreService().Operation().UpdateOperationChartPosition(ctx.Kit.Ctx, ctx.Kit.Header, opt)
	if err != nil {
		blog.Errorf("update chart position fail, err: %v", err)
		return
	}

	ctx.RespEntity(nil)
}
