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
	sendMessage = make(chan int)
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
		s, asu, su, sd  int64
		content         string
	)
	//Check if body is right
	if err, a = initCheck(r.Body); err != nil {
		return
	}

	// get deployment info from apiserver
	if err, bodyContentByte = k8sapi.APIServerGet(a, token); err != nil {
		return
	}

	// judge the criteria of gray deployment timeline,
	if err, s = secondTransform(a.Gray.DurationOfStay); err != nil {
		return
	}
	if err, asu = secondTransform(a.Gray.AVersionStepWiseUp); err != nil {
		return
	}
	if err, su = secondTransform(a.Gray.BVersionStepWiseUp); err != nil {
		return
	}
	if err, sd = secondTransform(a.Gray.BVersionStepWiseDown); err != nil {
		return
	}
	a.Gray.DurationOfStay = json.Number(strconv.FormatFloat(float64(s), 'f', -1, 64))
	a.Gray.AVersionStepWiseUp = json.Number(strconv.FormatFloat(float64(asu), 'f', -1, 64))
	a.Gray.BVersionStepWiseUp = json.Number(strconv.FormatFloat(float64(su), 'f', -1, 64))
	a.Gray.BVersionStepWiseDown = json.Number(strconv.FormatFloat(float64(sd), 'f', -1, 64))

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

	go goroutinesSend(r, a)

	//obtain the request content and phone number
	content, f = alert.Main(r.URL.String(), a)
	if err = ding(a, content, f); err != nil {
		return
	}
	return
}

// error handle routines
func errHandle() {
	for {
		select {
		case <-errors:
			err := <-MyErrorChan
			log.Println(err.MyError)
		}
	}
}

// error handle routines
func goroutinesSend(r *http.Request, a datastructure.Request) {
	var (
		f       [1]string
		content string
	)
	for {
		select {
		case <-sendMessage:
			content, f = alert.Main(r.URL.String(), a)
			if err := ding(a, content, f); err != nil {
				MyErrorChan <- MyError{err}
				errors <- 1
				return
			}
			return
		}
	}
}

// gray routines
func pauseGoroutine(a datastructure.Request, bodyContentByte []byte, token *viper.Viper) {
	var (
		specBodyContentByte               []byte
		minReadySeconds                   int64
		v, d, aTimeUp, bTimeUp, bTimeDown int
		m, t, st                          float64
		err                               error
		replace                           = make(chan int)
		replicas                          = make(chan int)
		deletes                           = make(chan int)
		grayInterval                      = make(chan int)
		replicasUpRate                    = make(chan int)
		replicasDownRate                  = make(chan int)
		SCReplicasUpRate                  = make(chan int)
		offLineInterval                   = make(chan int)
		spec                              datastructure.MySpec
	)

	// gray deploy tiered rate
	t, _ = a.Gray.TieredRate.Float64()

	// old and new version exist time of duration
	os, _ := a.Gray.DurationOfStay.Int64()
	d = int(os)

	//old version deploy increase stepwise interval
	avs, _ := a.Gray.AVersionStepWiseUp.Int64()
	aTimeUp = int(avs)

	//new version deploy increase stepwise interval
	osu, _ := a.Gray.BVersionStepWiseUp.Int64()
	bTimeUp = int(osu)

	//old version deploy reduce stepwise interval
	osd, _ := a.Gray.BVersionStepWiseDown.Int64()
	bTimeDown = int(osd)

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
		replicasUpRate <- 1
	}()

	// Interval goroutines
	go func() {
		for {
			select {
			case <-grayInterval:
				log.Println("[Paused] grayInterval")
				for {
					time.Sleep(time.Duration(d) * time.Second)
					break
				}
				replace <- 1
			case <-offLineInterval:
				log.Println("[Paused] offLineInterval")
				for {
					time.Sleep(time.Duration(d) * time.Second)
					break
				}
				replicasDownRate <- 1
			}
		}
	}()

	// replicas goroutines
	go func() {
		for {
			select {
			case <-SCReplicasUpRate:
				log.Println("[Paused] SCReplicasUpRate", t)
				for i := 1; i < int(math.Ceil(float64(10/int(t*10)))); i++ {
					st = t * float64(i)
					log.Println("[Paused] SCReplicasUpRate t :", st, i)
					if 1 < st {
						break
					}
					replicas <- 0
					replicas <- 3
					for {
						time.Sleep(time.Duration(aTimeUp) * time.Second)
						break
					}
				}
				offLineInterval <- 1
			case <-replicasUpRate:
				log.Println("[Paused] replicasUpRate", t)
				if err = k8sapi.APIServerPost(a, bodyContentByte, token); err != nil {
					MyErrorChan <- MyError{err}
					errors <- 1
				}
				for i := 1; i < int(math.Ceil(float64(10/int(t*10)))); i++ {
					st = t * float64(i)
					if st > 1 {
						break
					}
					replicas <- 0
					replicas <- 2
					log.Println("[Paused] replicasUpRate t :", st, i)
					for {
						time.Sleep(time.Duration(bTimeUp) * time.Second)
						break
					}
					log.Println("[Paused] replicasUpRate sleep down :")
				}
				grayInterval <- 1
			case <-replicasDownRate:
				log.Println("[Paused] replicasDownRate", t)
				for i := 1; i < int(math.Ceil(float64(10/int(t*10))))-1; i++ {
					st = t * float64(i)
					if st > 1 {
						break
					}
					replicas <- 0
					replicas <- 1
					for {
						time.Sleep(time.Duration(bTimeDown) * time.Second)
						break
					}
				}
				deletes <- 1
			}
		}
	}()

	// delete、replace、replicas goroutines
	for {
		select {
		case <-deletes:
			log.Println("[Paused] deletes")
			if err = k8sapi.APIServerDelete(a, token); err != nil {
				MyErrorChan <- MyError{err}
				errors <- 1
			}
		case <-replace:
			log.Println("[Paused] replace")
			log.Println("[Paused] replace", string(bodyContentByte))
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
			if err = k8sapi.APIServerPut(a, bodyContentByte, token); err != nil {
				MyErrorChan <- MyError{err}
				errors <- 1
			}
			log.Println("[Paused] replace ????", string(bodyContentByte))
			log.Println("SCReplicasUpRate???")
			SCReplicasUpRate <- 1
		case <-replicas:
			if v = <-replicas; v == 1 {
				log.Println("[Paused] replicasDown")
				f := func() (l float64) {
					m, _ = a.Replicas.Float64()
					return math.Ceil(m * (1 - st))
				}
				spec.Spec.Replicas = int(f())
				fmt.Println("replicasDown:", spec.Spec.Replicas)
			} else if v == 2 {
				log.Println("[Paused] replicasUp")
				f := func() (l float64) {
					m, _ = a.Replicas.Float64()
					return math.Ceil(m * st)
				}
				spec.Spec.Replicas = int(f())
			} else if v == 3 {
				log.Println("[Paused] SCReplicasUp")
				f := func() (l float64) {
					m, _ = a.Replicas.Float64()
					return math.Ceil(m * st)
				}
				spec.Spec.Replicas = int(f())
				if err, specBodyContentByte = ReplaceResourceReplicas(spec); err != nil {
					MyErrorChan <- MyError{err}
					errors <- 1
				}
				if err = k8sapi.APIServerPatch(a, specBodyContentByte, token, ""); err != nil {
					MyErrorChan <- MyError{err}
					errors <- 1
				}
				continue
			}
			if err, specBodyContentByte = ReplaceResourceReplicas(spec); err != nil {
				MyErrorChan <- MyError{err}
				errors <- 1
			}
			if err = k8sapi.APIServerPatch(a, specBodyContentByte, token, "temp"); err != nil {
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
