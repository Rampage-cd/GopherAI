package image

import (
	"GopherAI/common/image"
	"io"
	"log"
	"mime/multipart" //HTTP文件上传的标准类型
)

// 为HTTP层服务的，不是底层推理层
func RecognizeImage(file *multipart.FileHeader) (string, error) {

	modelPath := "/root/models/mobilenetv2/mobilenetv2-7.onnx"
	//指向ONNX模型文件，模型为MobileNetV2
	labelPath := "/root/imagenet_classes.txt" //共有1000个分类标签
	//在ImageNet上训练
	inputH, inputW := 224, 224
	//输入尺寸224*224

	recognizer, err := image.NewImageRecognizer(modelPath, labelPath, inputH, inputW)
	if err != nil {
		log.Println("NewImageRecognizer fail err is : ", err)
		return "", err
	}
	defer recognizer.Close() //防止内存泄漏（持有系统资源）

	src, err := file.Open() //multipart.FileHeader只是元数据，真正的数据需要Open()
	if err != nil {
		log.Println("file open fail err is : ", err)
		return "", err
	}
	defer src.Close() //关闭文件

	buf, err := io.ReadAll(src) //读取文件，得到原始图片数据
	if err != nil {
		log.Println("io.ReadAll fail err is : ", err)
		return "", err
	}

	return recognizer.PredictFromBuffer(buf) //从图片数据预测
}
