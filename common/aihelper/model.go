package aihelper

//对于不同AI大模型进行统一封装
//通过接口抽象，屏蔽OpenAI/Ollama等模型差异

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	//eino提供的大模型扩展组件
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"

	//通用的ChatModel接口定义
	"github.com/cloudwego/eino/components/model"

	//eino中用于表示消息结构的schema
	"github.com/cloudwego/eino/schema"
)

// 流式输出的回调函数类型
type StreamCallback func(msg string)

// AIModel 定义AI模型接口(上层业务只依赖AIModel，底层可以自由切换OpenAI/Ollama/其他模型)
type AIModel interface {
	//一次性生成完整回复（非流式调用）
	GenerateResponse(ctx context.Context, messages []*schema.Message) (*schema.Message, error)
	//逐步生成回复内容，并通过回调实时返回（流式调用）
	StreamResponse(ctx context.Context, messages []*schema.Message, cb StreamCallback) (string, error)
	//返回模型类型标识
	GetModelType() string
}

// =================== OpenAI 实现 ===================
type OpenAIModel struct {
	llm model.ToolCallingChatModel //是eino对ChatCompletion的统一抽象
}

func NewOpenAIModel(ctx context.Context) (*OpenAIModel, error) {
	key := os.Getenv("OPENAI_API_KEY")
	modelName := os.Getenv("OPENAI_MODEL_NAME")
	baseURL := os.Getenv("OPENAI_BASE_URL")

	//创建OpenAI ChatModel实例
	llm, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: baseURL,
		Model:   modelName,
		APIKey:  key,
	})
	if err != nil {
		return nil, fmt.Errorf("create openai model failed: %v", err)
	}
	return &OpenAIModel{llm: llm}, nil
}

func (o *OpenAIModel) GenerateResponse(ctx context.Context, messages []*schema.Message) (*schema.Message, error) {
	resp, err := o.llm.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("openai generate failed: %v", err)
	}
	return resp, nil
}

// 流式输出流程
// 1.调用Stream()建立流式连接
// 2.不断Recv()接收模型返回的delta
// 3.每收到一段就调用cb推送
// 4.同时在本地聚合完整结果
func (o *OpenAIModel) StreamResponse(ctx context.Context, messages []*schema.Message, cb StreamCallback) (string, error) {
	stream, err := o.llm.Stream(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("openai stream failed: %v", err)
	}
	defer stream.Close()

	var fullResp strings.Builder

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("openai stream recv failed: %v", err)
		}
		if len(msg.Content) > 0 {
			fullResp.WriteString(msg.Content) // 聚合

			cb(msg.Content) // 实时调用cb函数，方便主动发送给前端
		}
	}

	return fullResp.String(), nil //返回完整内容，方便后续存储
}

func (o *OpenAIModel) GetModelType() string { return "openai" }

// =================== Ollama 实现 ===================

// OllamaModel Ollama模型实现
type OllamaModel struct {
	llm model.ToolCallingChatModel
}

func NewOllamaModel(ctx context.Context, baseURL, modelName string) (*OllamaModel, error) {
	llm, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: baseURL,
		Model:   modelName,
	})
	if err != nil {
		return nil, fmt.Errorf("create ollama model failed: %v", err)
	}
	return &OllamaModel{llm: llm}, nil
}

func (o *OllamaModel) GenerateResponse(ctx context.Context, messages []*schema.Message) (*schema.Message, error) {
	resp, err := o.llm.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("ollama generate failed: %v", err)
	}
	return resp, nil
}

func (o *OllamaModel) StreamResponse(ctx context.Context, messages []*schema.Message, cb StreamCallback) (string, error) {
	stream, err := o.llm.Stream(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("ollama stream failed: %v", err)
	}
	defer stream.Close()
	var fullResp strings.Builder
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("openai stream recv failed: %v", err)
		}
		if len(msg.Content) > 0 {
			fullResp.WriteString(msg.Content) // 聚合
			cb(msg.Content)                   // 实时调用cb函数，方便主动发送给前端
		}
	}
	return fullResp.String(), nil //返回完整内容，方便后续存储
}

func (o *OllamaModel) GetModelType() string { return "ollama" }
