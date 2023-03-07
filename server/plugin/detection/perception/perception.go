package perception

import (
	"flag"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/detection/model"
	"gocv.io/x/gocv/cuda"
	"image"
	"image/color"
	"math"
	"strings"
	"time"

	"gocv.io/x/gocv"
)

var (
	//modelFile      = flag.String("model", "yolov5s.onnx", "model path")
	modelImageSize = flag.Int("size", 640, "model image size")
	srcImage       = flag.String("image", "images/face.jpg", "input image")
)

func printDevices() {
	num := cuda.GetCudaEnabledDeviceCount() 
	for i := 0; i < num; i++ {
		cuda.PrintCudaDeviceInfo(i)
	}
}

func Yolov5(modelFile string, app string) {
	flag.Parse()

	printDevices()

	net := gocv.ReadNetFromONNX(modelFile)
	net.SetPreferableBackend(gocv.NetBackendCUDA)
	net.SetPreferableTarget(gocv.NetTargetCUDA)

	modelSize := image.Pt(*modelImageSize, *modelImageSize)

	unconnectedLayerIds := net.GetUnconnectedOutLayers()
	layerNames := []string{}
	for _, id := range unconnectedLayerIds {
		layer := net.GetLayer(id)
		layerNames = append(layerNames, layer.GetName())
	}

	var outs []gocv.Mat
	for {
		if global.GVA_DB == nil {
			time.Sleep(time.Second * 10)
			continue
		}
		db := global.GVA_DB.Model(&model.DetectionFileUploadAndDownload{})
		var fileLists []model.DetectionFileUploadAndDownload
		db = db.Where("app = '" + app + "'")
		db = db.Where("url_detection = '' or url_detection isnull")
		err := db.Order("created_at desc").Find(&fileLists).Error
		if err != nil {
			panic(err)
		}

		for ii := range fileLists {
			start0 := time.Now()
			*srcImage = fileLists[ii].Url
			fmt.Print(fileLists[ii].Url)

			// detect 100 times
			resized := gocv.NewMat()
			src := gocv.IMRead(*srcImage, gocv.IMReadColor)
			letterBox(src, &resized, modelSize)

			blob := gocv.BlobFromImage(resized, 1/255.0, modelSize, gocv.Scalar{}, true, false)

			start := time.Now()
			net.SetInput(blob, "")
			outs = net.ForwardLayers(layerNames)
			end := time.Now()
			fmt.Println("cost", end.Sub(start))
			//}

			sz := outs[0].Size()
			rows := sz[1]
			cols := sz[2]

			ptr, _ := outs[0].DataPtrFloat32()

			boxes := []image.Rectangle{}
			scores := []float32{}
			indices := []int{}
			classIndexLists := []int{}

			for j := 0; j < rows; j++ {
				i0 := j * cols
				i1 := j*cols + cols
				line := ptr[i0:i1]
				x := line[0]
				y := line[1]
				w := line[2]
				h := line[3]
				sc := line[4]
				confs := line[5:]
				bestId, bestScore := getBestFromConfs(confs)
				bestScore *= sc

				scores = append(scores, bestScore)
				boxes = append(boxes, calculateBoundingBox(src, []float32{x, y, w, h}))
				indices = append(indices, -1)
				classIndexLists = append(classIndexLists, bestId)
			}

			fmt.Println("Do NMS in", len(boxes), "boxes")
			gocv.NMSBoxes(boxes, scores, 0.25, 0.45, indices)

			nmsNumber := 0
			goodBoxes := []image.Rectangle{}
			goodScores := []float32{}
			goodClassIds := []int{}

			output := src.Clone()

			for _, v := range indices {
				if v < 0 {
					break
				} else {
					nmsNumber++
					goodBoxes = append(goodBoxes, boxes[v])
					goodScores = append(goodScores, scores[v])
					goodClassIds = append(goodClassIds, classIndexLists[v])

					gocv.Rectangle(&output, boxes[v], color.RGBA{0, 255, 0, 255}, 1)
					gocv.PutText(&output, fmt.Sprintf("yftech_%d:%.02f", classIndexLists[v], scores[v]),
						boxes[v].Min, gocv.FontHersheySimplex, 0.5, color.RGBA{255, 0, 0, 255}, 1)
				}
			}

			fmt.Println("After NMS", nmsNumber, "keeped")

			//w := gocv.NewWindow("detected")
			//w.ResizeWindow(modelSize.X, modelSize.Y)
			//w.IMShow(output)
			//w.WaitKey(-1)
			newurl := strings.Replace(fileLists[ii].Url, "file", "tmp", 1)
			gocv.IMWrite(newurl, output)
			db := global.GVA_DB.Model(&model.DetectionFileUploadAndDownload{})
			err = db.Where("url = ?", fileLists[ii].Url).Update("url_detection", newurl).Error
			if err != nil {
				return
			}
			end = time.Now()
			fmt.Println("all cost", end.Sub(start0))
			src.Close()
		}
		time.Sleep(time.Second)
	}
}

func Yolov8seg(modelFile string, app string) {
	flag.Parse()
	th:=float32(0.1)
	//maskWeightCn := 1
	printDevices()

	net := gocv.ReadNetFromONNX(modelFile)
	net.SetPreferableBackend(gocv.NetBackendCUDA)
	net.SetPreferableTarget(gocv.NetTargetCUDA)

	modelSize := image.Pt(*modelImageSize, *modelImageSize)

	unconnectedLayerIds := net.GetUnconnectedOutLayers()
	layerNames := []string{}
	for _, id := range unconnectedLayerIds {
		layer := net.GetLayer(id)
		layerNames = append(layerNames, layer.GetName())
		//break
	}

	var outs []gocv.Mat
	for {
		if global.GVA_DB == nil {
			time.Sleep(time.Second * 10)
			continue
		}
		db := global.GVA_DB.Model(&model.DetectionFileUploadAndDownload{})
		var fileLists []model.DetectionFileUploadAndDownload
		db = db.Where("app = '" + app + "'")
		db = db.Where("url_detection = '' or url_detection isnull")
		err := db.Order("created_at desc").Find(&fileLists).Error
		if err != nil {
			global.GVA_DB.AutoMigrate(model.DetectionFileUploadAndDownload{})
			time.Sleep(time.Second * 10)
			continue
		}

		for ii := range fileLists {
			start0 := time.Now()
			*srcImage = fileLists[ii].Url
			fmt.Print(fileLists[ii].Url)

			// detect 100 times
			resized := gocv.NewMat()
			src := gocv.IMRead(*srcImage, gocv.IMReadColor)
			letterBox(src, &resized, modelSize)

			blob := gocv.BlobFromImage(resized, 1/255.0, modelSize, gocv.Scalar{}, true, false)

			start := time.Now()
			net.SetInput(blob, "")
			outs = net.ForwardLayers(layerNames)
			end := time.Now()
			fmt.Println("cost", end.Sub(start))

			//}

			sz := outs[0].Size()
			rows := sz[1]
			cols := sz[2]
			out0 := outs[0].Reshape(1, rows)
			outt := out0.T()

			sz1 := outs[1].Size()
			cns1 := sz1[1]
			rows1 := sz1[2]
			cols1 := sz1[3]
			//mask := outs[1].Reshape(1, cns1)
			//mask = mask.T()
			//mask = mask.Reshape(1, rows1)
			//fmt.Println("mask:", mask.Size())
			//fmt.Println("out1:", outs[1].Size())

			ptr, _ := outt.DataPtrFloat32()
			ptr1, _ := outs[1].DataPtrFloat32()

			boxes := []image.Rectangle{}
			resizedBoxes := []image.Rectangle{}
			scores := []float32{}
			indices := []int{}
			classIndexLists := []int{}
			maskweights := [][]float32{}

			for j := 0; j < cols; j++ {
				i0 := j * rows
				i1 := j*rows + rows
				line := ptr[i0:i1]
				x := line[0]
				y := line[1]
				w := line[2]
				h := line[3]
				confs := line[4 : rows-cns1]
				bestId, bestScore := getBestFromConfs(confs)
				if bestScore<th{
					continue
				}
				scores = append(scores, bestScore)
				boxes = append(boxes, calculateBoundingBox(src, []float32{x, y, w, h}))
				resizedBoxes = append(resizedBoxes, image.Rect(int(x-w/2)*cols1/modelSize.X, int(y-h/2)*rows1/modelSize.Y, int(x+w/2)*cols1/modelSize.X, int(y+h/2)*rows1/modelSize.Y))
				indices = append(indices, -1)
				classIndexLists = append(classIndexLists, bestId)
				maskweights = append(maskweights, line[rows-cns1:])
			}

			fmt.Println("Do NMS in", len(boxes), "boxes")
			if len(boxes)>1 {
				gocv.NMSBoxes(boxes, scores, 0.25, 0.45, indices)
			}

			nmsNumber := 0
			goodBoxes := []image.Rectangle{}
			goodScores := []float32{}
			goodClassIds := []int{}
			goodMaskWeights := [][]float32{}

			output := src.Clone()
			newurl := strings.Replace(fileLists[ii].Url, "file", "tmp", 1)
			for _, v := range indices {
				if v < 0 {
					break
				} else {
					nmsNumber++
					goodBoxes = append(goodBoxes, boxes[v])
					goodScores = append(goodScores, scores[v])
					goodClassIds = append(goodClassIds, classIndexLists[v])
					goodMaskWeights = append(goodMaskWeights, maskweights[v])

					gocv.Rectangle(&output, boxes[v], color.RGBA{0, 255, 0, 255}, 1)
					gocv.PutText(&output, fmt.Sprintf("yftech_%d:%.02f", classIndexLists[v], scores[v]),
						boxes[v].Min, gocv.FontHersheySimplex, 0.5, color.RGBA{255, 0, 0, 255}, 1)
					resizedBox := resizedBoxes[v]
					//fmt.Println(resizedBox.Min.X)
					maskt := gocv.NewMatWithSize(resizedBox.Dy(), resizedBox.Dx(), gocv.MatTypeCV8UC1)
					grid := rows1 * cols1
					for k := 0; k < grid; k++ {
						yy := k / cols1
						xx := k % cols1
						sum := float32(0.)
						if xx >= resizedBox.Min.X && xx < resizedBox.Max.X && yy >= resizedBox.Min.Y && yy < resizedBox.Max.Y {
							for l := 0; l < cns1; l++ {
								sum += maskweights[v][l] * ptr1[l*grid+k]
							}
							gt0 := 0
							//fmt.Println(sigmoid(sum))
							if sigmoid(sum) > 0.5 {
								gt0 = 255
							}
							maskt.SetSCharAt(yy-resizedBox.Min.Y, xx-resizedBox.Min.X, int8(gt0))
						}
					}
					masko := gocv.NewMat()
					gocv.Resize(maskt, &masko, image.Point{boxes[v].Dx(), boxes[v].Dy()}, 0, 0, gocv.InterpolationNearestNeighbor)
					//color_mask := gocv.NewMatWithSize(boxes[v].Dx(), boxes[v].Dy(), gocv.MatTypeCV8SC3)
					//color_mask.CopyToWithMask(&color_mask, masko)

					for u := 0; u < boxes[v].Dx(); u++ {
						for t := 0; t < boxes[v].Dy(); t++ {
							if masko.GetUCharAt(t, u) == 255 {
								output.SetUCharAt(t+boxes[v].Min.Y, (u+boxes[v].Min.X)*output.Channels()+0, 114)
								output.SetUCharAt(t+boxes[v].Min.Y, (u+boxes[v].Min.X)*output.Channels()+1, 114)
								output.SetUCharAt(t+boxes[v].Min.Y, (u+boxes[v].Min.X)*output.Channels()+2, 114)
								//output.SetUCharAt3(t+boxes[v].Min.Y, u+boxes[v].Min.X, 2,0)
							}
						}
					}
					//gocv.IMWrite(newurl+string(v)+".jpg", masko)
				}
			}

			fmt.Println("After NMS", nmsNumber, "keeped")

			//w := gocv.NewWindow("detected")
			//w.ResizeWindow(modelSize.X, modelSize.Y)
			//w.IMShow(output)
			//w.WaitKey(-1)

			gocv.IMWrite(newurl, output)
			//var filet []model.DetectionFileUploadAndDownload
			db := global.GVA_DB.Model(&model.DetectionFileUploadAndDownload{})
			err = db.Where("url = ?", fileLists[ii].Url).Update("url_detection", newurl).Error
			if err != nil {
				return
			}
			end = time.Now()
			fmt.Println("all cost", end.Sub(start0))
			src.Close()
			output.Close()
		}
		time.Sleep(time.Second * 10)
	}
}

func sigmoid(x float32) float32 {
	xx := float32(1 / (1 + math.Exp(-float64(x))))
	return xx
}
func getBestFromConfs(confs []float32) (int, float32) {
	bestId := 0
	bestScore := float32(0)
	for i, v := range confs {
		if v > bestScore {
			bestId = i
			bestScore = v
		}
	}
	return bestId, bestScore
}

func letterBox(src gocv.Mat, dst *gocv.Mat, size image.Point) {
	k := math.Min(float64(size.X)/float64(src.Cols()), float64(size.Y)/float64(src.Rows()))
	newSize := image.Pt(int(k*float64(src.Cols())), int(k*float64(src.Rows())))

	tmp := gocv.NewMat()
	gocv.Resize(src, &tmp, newSize, 0, 0, gocv.InterpolationLinear)

	if dst.Cols() != size.X || dst.Rows() != size.Y {
		dstNew := gocv.NewMatWithSize(size.Y, size.X, src.Type())
		dstNew.CopyTo(dst)
	}

	rectOfTmp := image.Rect((dst.Cols()-newSize.X)/2, (dst.Rows()-newSize.Y)/2, (dst.Cols()-newSize.X)/2+newSize.X, (dst.Rows()-newSize.Y)/2+newSize.Y)

	regionOfDst := dst.Region(rectOfTmp)
	tmp.CopyTo(&regionOfDst)
}

// calculateBoundingBox calculate the bounding box of the detected object.
func calculateBoundingBox(frame gocv.Mat, row []float32) image.Rectangle {
	if len(row) < 4 {
		return image.Rect(0, 0, 0, 0)
	}
	gain := math.Min(float64(640)/float64(frame.Cols()), float64(640)/float64(frame.Rows()))
	padx := float32(float64(640)-float64(frame.Cols())*gain) / 2
	pady := float32(float64(640)-float64(frame.Rows())*gain) / 2

	x, y, w, h := row[0], row[1], row[2], row[3]
	left := int((x - padx - 0.5*w) / float32(gain))
	top := int((y - pady - 0.5*h) / float32(gain))
	width := int(w / float32(gain))
	height := int(h / float32(gain))

	return image.Rect(left, top, left+width, top+height)
}
