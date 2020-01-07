package goapi

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"io"
	"main/grpc"
	"main/src/itineris"
)

/*
PApiServiceServer is gRPC gateway server

@author Thanh Nguyen <btnguyen2k@gmail.com>
@since template-v0.4.r2
*/
type PApiServiceServer struct {
}

func (s *PApiServiceServer) Ping(context.Context, *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (s *PApiServiceServer) Check(context.Context, *grpc.PApiAuth) (*grpc.PApiResult, error) {
	return &grpc.PApiResult{
		Status:     itineris.StatusOk,
		Message:    "Ok",
		Encoding:   grpc.PDataEncoding_JSON_STRING,
		ResultData: nil,
		DebugData:  nil,
	}, nil
}

func (s *PApiServiceServer) Call(_ context.Context, gctx *grpc.PApiContext) (*grpc.PApiResult, error) {
	ctx := itineris.NewApiContext().SetApiName(gctx.ApiName).SetGateway("GRPC")
	auth := itineris.NewApiAuth(gctx.ApiAuth.AppId, gctx.ApiAuth.AccessToken)
	params := parseParams(gctx.ApiParams)
	if params == nil {
		result := itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Cannot parse request parameters.")
		return toPApiResult(grpc.PDataEncoding_JSON_STRING, result), nil
	}
	result := ApiRouter.CallApi(ctx, auth, params)
	resultEncoding := gctx.ApiParams.ExpectedReturnEncoding
	if resultEncoding == grpc.PDataEncoding_JSON_DEFAULT {
		resultEncoding = gctx.ApiParams.Encoding
		if resultEncoding == grpc.PDataEncoding_JSON_DEFAULT {
			resultEncoding = grpc.PDataEncoding_JSON_STRING
		}
	}
	return toPApiResult(resultEncoding, result), nil
}

func newGrpcGateway() *PApiServiceServer {
	return &PApiServiceServer{}
}

func toPApiResult(encoding grpc.PDataEncoding, apiResult *itineris.ApiResult) *grpc.PApiResult {
	result := &grpc.PApiResult{
		Status:   int32(apiResult.Status),
		Message:  apiResult.Message,
		Encoding: encoding,
	}
	{
		js, _ := json.Marshal(apiResult.Data)
		if encoding == grpc.PDataEncoding_JSON_GZIP {
			js, _ = gzipEncode(js)
		}
		result.ResultData = js
	}
	{
		js, _ := json.Marshal(apiResult.DebugInfo)
		if encoding == grpc.PDataEncoding_JSON_GZIP {
			js, _ = gzipEncode(js)
		}
		result.DebugData = js
	}

	return result
}

func parseParams(gparams *grpc.PApiParams) *itineris.ApiParams {
	data := make(map[string]interface{})
	switch gparams.Encoding {
	case grpc.PDataEncoding_JSON_DEFAULT, grpc.PDataEncoding_JSON_STRING:
		err := json.Unmarshal(gparams.ParamsData, &data)
		if err != nil {
			return nil
		}
	case grpc.PDataEncoding_JSON_GZIP:
		buf, err := gzipDecode(gparams.ParamsData)
		fmt.Println(len(gparams.ParamsData), len(buf))
		fmt.Println("Encoding Gzip", err)
		if err != nil {
			return nil
		}
		err = json.Unmarshal(buf, &data)
		if err != nil {
			return nil
		}
	default:
		return nil
	}
	params := itineris.NewApiParams()
	for k, v := range data {
		params.SetParam(k, v)
	}
	return params
}

func gzipEncode(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	defer w.Close()
	_, err := w.Write(input)
	if err == nil {
		w.Close()
	}
	return buf.Bytes(), err
}

func gzipDecode(input []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewBuffer(input))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	var result []byte
	buf := make([]byte, 1024)
	for {
		count, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if count > 0 {
			result = append(result, buf[0:count]...)
		} else {
			break
		}
	}
	return result, nil
}
