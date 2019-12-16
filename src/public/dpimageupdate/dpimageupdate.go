package dpimageupdate

import (
	"alert"
	"datastructure"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/syyongx/php2go"
	"io"
	"io/ioutil"
	"k8sapi"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"user"
)

type MyError struct {
	MyError error
}

var (
	errors      = make(chan int)
	MyErrorChan = make(chan MyError)
)

func initCheck(rBody io.Reader) (err error, a datastructure.Request) {
	var (
		body []byte
	)
	// if the body exist
	if body, err = ioutil.ReadAll(rBody); err != nil {
		log.Printf("[InitCheck] Read Body ERR: %v\n", err)
		err = fmt.Errorf("[InitCheck] Read Body ERR: %v\n", err)
		return
	}

	// if the body can be turn to json
	if err = json.Unmarshal(body, &a); err != nil {
		log.Printf("[InitCheck] Unmarshal Body ERR: %v", err)
		err = fmt.Errorf("[InitCheck] Unmarshal Body ERR: %v", err)
		return
	}

	//judge the user if exist
	if err = user.User(a); err != nil {
		return
	}

	// log the parameter
	if parameter, err := json.Marshal(a); err == nil {
		log.Printf("[InitCheck] The Request Body: %v", string(parameter))
	}
	return
}

func ding(a datastructure.Request, content string, f [1]string) (err error) {
	//dingDing alert
	if err = alert.Ding(content, f, a.SendFormat); err != nil {
		log.Printf("[Ding] Dingding ERROR:[%s]", err)
		err = fmt.Errorf("[Ding] Dingding ERROR:[%v] %v", err,
			"\n DingAlert Filed, But Request Has Been Done, Do Not Worry !")
		return
	}
	return
}

func anonymousReplace(a datastructure.Request, f func(datastructure.Request) (err error)) (err error) {
	return f(a)
}

func DpUpdate(r *http.Request, token *viper.Viper) (err error) {
	var (
		a, newA         datastructure.Request
		f               [1]string
		bodyContentByte []byte
		content         string
	)
	//Check if body is right
	if err, a = initCheck(r.Body); err != nil {
		return
	}

	// get deployment info from apiServer
	if err, bodyContentByte = k8sapi.APIServerGet(a, token); err != nil {
		return
	}

	//replace the resource
	//the anonymous func is equivalent to func replace
	if err = anonymousReplace(a, func(a datastructure.Request) (err error) {
		// eliminate the Status from deployment
		if err, bodyContentByte = eliminateStatus(bodyContentByte); err != nil {
			return
		}
		//replace resource from deployment, include image, replicas
		if err, bodyContentByte, newA = replaceResource(a, bodyContentByte); err != nil {
			return
		}
		if newA.Replicas != "" {
			a = newA
		}
		return
	}); err != nil {
		return
	}

	// put the new deployment info to apiServer
	if err = k8sapi.APIServerPut(a, bodyContentByte, token); err != nil {
		return
	}

	//obtain the request content and phone number
	content, f = alert.Main(r.URL.String(), a)
	if err = ding(a, content, f); err != nil {
		return
	}
	return
}

func GrayDpUpdate(r *http.Request, token *viper.Viper) (err error) {
	var (
		a, newA         datastructure.Request
		f               [1]string
		bodyContentByte []byte
		s, su, sd       int64
		//s       int64
		content string
	)
	//Check if body is right
	if err, a = initCheck(r.Body); err != nil {
		return
	}

	// get deployment info from apiserver
	if err, bodyContentByte = k8sapi.APIServerGet(a, token); err != nil {
		return
	}

	//gray deployment timeline
	if err, s = secondTransform(a.Gray.DurationOfStay); err != nil {
		return
	}
	if err, su = secondTransform(a.Gray.TempStepWiseUp); err != nil {
		return
	}
	if err, sd = secondTransform(a.Gray.TempStepWiseDown); err != nil {
		return
	}
	fmt.Println("s, su, sd", s, su, sd)

	//a.DurationOfStay = s
	a.Gray.DurationOfStay = json.Number(strconv.FormatFloat(float64(s), 'f', -1, 64))
	a.Gray.TempStepWiseUp = json.Number(strconv.FormatFloat(float64(su), 'f', -1, 64))
	a.Gray.TempStepWiseDown = json.Number(strconv.FormatFloat(float64(sd), 'f', -1, 64))
	fmt.Println("durationOfStay: ?", a.Gray.DurationOfStay)

	// TieredRate if right
	tieredRate := a.Gray.TieredRate
	if strings.Contains(a.Gray.TieredRate.String(), "%") {
		if s, err = strconv.ParseInt(tieredRate.String()[0:len(tieredRate.String())-1], 10, 64); err != nil {
			log.Printf("[GrayDpUpdate] {%v} Is Not Number In %v", tieredRate[0:len(tieredRate)-1], tieredRate)
			err = fmt.Errorf("[GrayDpUpdate] {%v} Is Not Number In %v", tieredRate[0:len(tieredRate)-1], tieredRate)
			return
		}
		a.Gray.TieredRate = json.Number(strconv.FormatFloat(float64(s)*0.01, 'f', -1, 64))
	}

	//replace the resource
	//the anonymous func is equivalent to func replace
	if err = anonymousReplace(a, func(a datastructure.Request) (err error) {
		// eliminate the Status from deployment
		if err, bodyContentByte = eliminateStatus(bodyContentByte); err != nil {
			return
		}
		//replace resource from deployment, include image, replicas
		if err, bodyContentByte, newA = replaceResource(a, bodyContentByte); err != nil {
			return
		}
		if newA.Replicas != "" {
			a = newA
		}
		return
	}); err != nil {
		return
	}

	//gray deployment controller Goroutine
	go pauseGoroutine(a, bodyContentByte, token)

	//  handle the err of pauseGoroutine,if err exist
	go errHandle()

	//obtain the request content and phone number
	content, f = alert.Main(r.URL.String(), a)
	if err = ding(a, content, f); err != nil {
		return
	}
	return
}

func errHandle() {
	for {
		select {
		case <-errors:
			err := <-MyErrorChan
			log.Println(err.MyError)
		}
	}
}

func pauseGoroutine(a datastructure.Request, bodyContentByte []byte, token *viper.Viper) {
	var (
		minReadySeconds int64
		d, swd          int
		m, t            float64
		err             error
		interval        = make(chan int)
		replace         = make(chan int)
		replicas        = make(chan int)
		deletes         = make(chan int)
		grayInterval    = make(chan int)
		replicasRate    = make(chan int)
		spec            datastructure.MySpec
	)

	t, _ = a.Gray.TieredRate.Float64()
	os, _ := a.Gray.DurationOfStay.Int64()
	d = int(os)
	//osu, _ := a.Gray.TempStepWiseUp.Int64()
	//swu = int(osu)
	osd, _ := a.Gray.TempStepWiseDown.Int64()
	swd = int(osd)

	// create new deployment
	go func() {
		minReadySeconds, _ = a.MinReadySeconds.Int64()
		if minReadySeconds > 10 {
			minReadySeconds = 10
		}
		log.Printf("[Paused] CoolingTime Need TO %v Gray Update", d+int(minReadySeconds)*2)
		for {
			time.Sleep(time.Duration(minReadySeconds) * time.Second)
			break
		}
		if err = k8sapi.APIServerPost(a, bodyContentByte, token); err != nil {
			return
		}
		interval <- 1
	}()

	go func() {
		for {
			select {
			case <-interval:
				log.Println("[Paused] Interval")
				for {
					time.Sleep(time.Duration(d) * time.Second)
					break
				}
				replace <- 1
			case <-grayInterval:
				log.Println("[Paused] grayInterval")
				for {
					time.Sleep(time.Duration(d) * time.Second)
					break
				}
				replicasRate <- 1
			}
		}
	}()

	go func() {
		select {
		case <-replicasRate:
			log.Println("[Paused] replicasRate", t)
			for i := 1; i <= int(math.Ceil(float64(10/int(t*10)))); i++ {
				t = t * float64(i)
				replicas <- 1
				//if int(swd)/2 > 60 {
				//	swd = 60
				//}
				for {
					time.Sleep(time.Duration(swd) * time.Second)
					break
				}
			}
			deletes <- 1
		}
	}()

	for {
		select {
		case <-replace:
			log.Println("[Paused] replace")
			if err, bodyContentByte = eliminateStatus(bodyContentByte); err != nil {
				MyErrorChan <- MyError{err}
				errors <- 1
			}
			//replace resource from deployment, include image, replicas
			a.Name = "InstantDeployment"
			if err, bodyContentByte = ReplaceResourceName(a, bodyContentByte); err != nil {
				MyErrorChan <- MyError{err}
				errors <- 1
			}
			// put the new deployment info to apiServer
			if err = k8sapi.APIServerPut(a, bodyContentByte, token); err != nil {
				MyErrorChan <- MyError{err}
				errors <- 1
			}
			grayInterval <- 1
		case <-deletes:
			log.Println("[Paused] deletes")
			if err = k8sapi.APIServerDelete(a, token); err != nil {
				MyErrorChan <- MyError{err}
				errors <- 1
			}
		case <-replicas:
			log.Println("[Paused] replicas")
			f := func() (l float64) {
				m, _ = a.Replicas.Float64()
				return math.Ceil(m * (1 - t))
			}
			spec.Spec.Replicas = int(f())
			if err, bodyContentByte = ReplaceResourceReplicas(spec); err != nil {
				MyErrorChan <- MyError{err}
				errors <- 1
			}
			//fmt.Println("[Paused] replicas ", string(bodyContentByte))
			if err = k8sapi.APIServerPatch(a, bodyContentByte, token); err != nil {
				MyErrorChan <- MyError{err}
				errors <- 1
			}
		}
	}
}

func secondTransform(i json.Number) (err error, s int64) {
	fmt.Println(i)
	if i == "" {
		s = 60
	} else if strings.Contains(i.String(), "min") {
		//strconv.ParseInt(a.Paused.String(), 10, 64)
		if s, err = strconv.ParseInt(i.String()[0:len(i)-3], 10, 64); err != nil {
			fmt.Println(s, err)
			log.Printf("[GrayDpUpdate] {%v} Is Not Number In %v", i[0:len(i)-3], i)
			err = fmt.Errorf("[GrayDpUpdate] {%v} Is Not Number In %v", i[0:len(i)-3], i)
			return
		}
		s = s * 60
	} else if strings.Contains(i.String(), "s") {
		if s, err = strconv.ParseInt(i.String()[0:len(i)-1], 10, 64); err != nil {
			fmt.Println(s, err)
			log.Printf("[GrayDpUpdate] {%v} Is Not Number In %v", i[0:len(i)-1], i)
			err = fmt.Errorf("[GrayDpUpdate] {%v} Is Not Number In %v", i[0:len(i)-1], i)
			return
		}
	} else if php2go.IsNumeric(i) {
		s, _ = strconv.ParseInt(i.String(), 10, 64)
		log.Printf("[GrayDpUpdate] {%v} Has Not Unit, So Default Is Second", i)
	} else {
		log.Printf("[GrayDpUpdate] Paused: %v Is Null, "+
			"So GrayDeployment Paused Default Is 1 Minute. \n"+
			"Notice: GrayDeployment Are published Later More Than 1 Minute.", i)
		err = fmt.Errorf("[GrayDpUpdate] Paused: %v Is Null, "+
			"So GrayDeployment Paused Default Is 1 Minute. \n"+
			"Notice: GrayDeployment Are published Later More Than 1 Minute.", i)
	}
	return
}
