package image

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"sync"

	ort "github.com/yalue/onnxruntime_go"
	"golang.org/x/image/draw"
)

type ImageRecognizer struct {
	session      *ort.Session[float32] //ONNX推理会话
	inputName    string
	outputName   string
	inputH       int
	inputW       int
	labels       []string             //标签列表
	inputTensor  *ort.Tensor[float32] //复用的输入
	outputTensor *ort.Tensor[float32] //复用的输出
}

const (
	defaultInputName  = "data"
	defaultOutputName = "mobilenetv20_output_flatten0_reshape0"
) //默认的模型输入节点名和输出节点名

var (
	initOnce sync.Once //保证线程安全
	//sync.Once确保在并发环境下，某段代码在程序的整个生命周期中只被执行一次
	initErr error
)

// NewImageRecognizer 创建识别器（自动使用默认 input/output 名称）
func NewImageRecognizer(modelPath, labelPath string, inputH, inputW int) (*ImageRecognizer, error) {
	if inputH <= 0 || inputW <= 0 {
		inputH, inputW = 224, 224
	}

	// 初始化 ONNX 环境（全局一次）
	initOnce.Do(func() {
		initErr = ort.InitializeEnvironment()
	})
	if initErr != nil {
		return nil, fmt.Errorf("onnxruntime initialize error: %w", initErr)
	}

	// 预先创建输入输出 Tensor
	inputShape := ort.NewShape(1, 3, int64(inputH), int64(inputW)) //3表示RGB通道
	inData := make([]float32, inputShape.FlattenedSize())          //为输入Tensor分配连续内存
	inTensor, err := ort.NewTensor(inputShape, inData)
	if err != nil {
		return nil, fmt.Errorf("create input tensor failed: %w", err)
	}

	outShape := ort.NewShape(1, 1000)
	outTensor, err := ort.NewEmptyTensor[float32](outShape)
	if err != nil {
		inTensor.Destroy() //释放已创建的输入Tensor
		return nil, fmt.Errorf("create output tensor failed: %w", err)
	}

	// 创建 Session
	session, err := ort.NewSession[float32](
		modelPath,
		[]string{defaultInputName},
		[]string{defaultOutputName},
		[]*ort.Tensor[float32]{inTensor},
		[]*ort.Tensor[float32]{outTensor},
	)
	if err != nil {
		inTensor.Destroy()
		outTensor.Destroy()
		return nil, fmt.Errorf("create onnx session failed: %w", err)
	}

	// 读取 label 文件
	labels, err := loadLabels(labelPath) //从文件中加载分类标签
	if err != nil {
		session.Destroy()
		inTensor.Destroy()
		outTensor.Destroy()
		return nil, err
	}

	return &ImageRecognizer{
		session:      session,
		inputName:    defaultInputName,
		outputName:   defaultOutputName,
		inputH:       inputH,
		inputW:       inputW,
		labels:       labels,
		inputTensor:  inTensor,
		outputTensor: outTensor,
	}, nil
}

func (r *ImageRecognizer) Close() {
	if r.session != nil {
		_ = r.session.Destroy()
		r.session = nil
	}
	if r.inputTensor != nil {
		_ = r.inputTensor.Destroy()
		r.inputTensor = nil
	}
	if r.outputTensor != nil {
		_ = r.outputTensor.Destroy()
		r.outputTensor = nil
	}
	//存在则销毁
	//同时避免悬空指针
}

// 从图片文件进行预测
func (r *ImageRecognizer) PredictFromFile(imagePath string) (string, error) {
	file, err := os.Open(filepath.Clean(imagePath)) //打开图片文件
	if err != nil {
		return "", fmt.Errorf("image not found: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file) //自动识别格式并解码图片
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	return r.PredictFromImage(img)
	//调用统一的image预测逻辑
}

// 从字节流进行预测
func (r *ImageRecognizer) PredictFromBuffer(buf []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(buf)) //从内存字节解码图片
	if err != nil {
		return "", fmt.Errorf("failed to decode image from buffer: %w", err)
	}
	return r.PredictFromImage(img)
}

// 核心预测函数
func (r *ImageRecognizer) PredictFromImage(img image.Image) (string, error) {

	resizedImg := image.NewRGBA(image.Rect(0, 0, r.inputW, r.inputH))
	//创建目标尺寸RGBA图像

	draw.CatmullRom.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Over, nil)
	//目标图像，原始图像，覆盖绘制模式（高质量插值算法进行缩放）

	h, w := r.inputH, r.inputW
	ch := 3                         // R, G, B
	data := make([]float32, h*w*ch) //创建NCHW格式输入数据切片

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := resizedImg.At(x, y) //获取当前位置的像素颜色

			r, g, b, _ := c.RGBA() //获取RGBA

			rf := float32(r>>8) / 255.0 //转化为0-1的float32
			gf := float32(g>>8) / 255.0
			bf := float32(b>>8) / 255.0

			// NCHW format
			data[y*w+x] = rf
			data[h*w+y*w+x] = gf
			data[2*h*w+y*w+x] = bf //RGB三种通道，每种通道占用一层
		}
	}

	inData := r.inputTensor.GetData() //获取输入Tensor的底层内存
	copy(inData, data)

	if err := r.session.Run(); err != nil { //进行ONNX推理
		return "", fmt.Errorf("onnx run error: %w", err)
	}

	outData := r.outputTensor.GetData()
	if len(outData) == 0 {
		return "", errors.New("empty output from model")
	}

	maxIdx := 0
	maxVal := outData[0]
	for i := 1; i < len(outData); i++ {
		if outData[i] > maxVal { //找到最大概率以及其对应的下标
			maxVal = outData[i]
			maxIdx = i
		}
	}

	if maxIdx >= 0 && maxIdx < len(r.labels) {
		return r.labels[maxIdx], nil //返回预测类别
	}
	return "Unknown", nil //无法匹配标签
}

// 从文件中加载标签
func loadLabels(path string) ([]string, error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("open label file failed: %w", err)
	}
	defer f.Close() //打开标签文件，并在函数结束时关闭

	var labels []string
	sc := bufio.NewScanner(f) //创建逐行扫描器
	for sc.Scan() {           //逐行扫描
		line := sc.Text() //获取当前文本
		if line != "" {
			labels = append(labels, line) //忽略空行
		}
	}
	if err := sc.Err(); err != nil { //扫描过程中发生错误
		return nil, fmt.Errorf("read labels failed: %w", err)
	}
	if len(labels) == 0 {
		return nil, fmt.Errorf("no labels found in %s", path)
	}
	return labels, nil //返回标签列表
}
