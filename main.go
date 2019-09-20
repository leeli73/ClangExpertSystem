package main
import(
	"os"
	"log"
	"fmt"
	"time"
	"strconv"
	"strings"
	"os/exec"
	"net/url"
	"net/http"
	"math/rand"
	"io/ioutil"
	"encoding/base64"
	"github.com/mndrix/golog"
)
type ErrorInfo struct {
	code string
	data string
}
type NewRule struct{
	question string
	operate string
	way string
}
var AllErrorInfo []ErrorInfo
var AllCPPCheckError []ErrorInfo
var AllNewRule []NewRule
func main(){
	Init()
	http.HandleFunc("/",Index)
	http.HandleFunc("/Code",Code)
	http.HandleFunc("/CheckCode",CheckCode)
	http.HandleFunc("/GetResult",GetResult)
	http.HandleFunc("/CheckCPP",CheckCPP)
	http.HandleFunc("/NewRule",NewRule1)
	http.HandleFunc("/admin",admin)
	http.HandleFunc("/work",work)
	log.Println("Listening Port 88...")
	if err := http.ListenAndServe(":88", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
func work(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	mytype := r.FormValue("type")
	if mytype == "del" {
		id, _ := strconv.Atoi(r.FormValue("id"))
		if len(AllNewRule) != 0 {
			if id < len(AllNewRule) {
				AllNewRule = append(AllNewRule[:id:id], AllNewRule[id+1:]...)
				w.Write([]byte("succecc"))
			} else {
				w.Write([]byte("error"))
			}
		}
	} else if mytype == "right" {
		id, _ := strconv.Atoi(r.FormValue("id"))
		if len(AllNewRule) != 0 {
			if id < len(AllNewRule) {
				temp := AllNewRule[id]
				AllNewRule = append(AllNewRule[:id:id], AllNewRule[id+1:]...)
				AddtoPl(temp.question,temp.operate,temp.way)
				Init()
				w.Write([]byte("succecc"))
			} else {
				w.Write([]byte("error"))
			}
		}
	}
}
func AddtoPl(q string,o string,w string){
	f, err := os.OpenFile("data/RunQuestion.pl", os.O_RDONLY,0600)
	defer f.Close()
	if err !=nil {
		log.Fatal(err)
		return
	} else {
		DataByte,err:=ioutil.ReadAll(f)
		if err != nil{
			log.Fatal(err)
			return
		}
		DataStr := string(DataByte)
		lines := strings.Split(DataStr,"\n")
		for count:=0 ; count < len(lines) ;count ++ {
			lines[count] = strings.TrimSpace(lines[count])
			lines[count] = strings.Replace(lines[count],"\\n","",-1)
		}
		var code1 []string
		code1Count := 0
		var code2 []string
		code2Count := 0
		var code3 []string
		code3Count := 0
		var plerror []string
		errorCount := 0
		for _,v := range lines {
			if v != "" {
				if v == "/*" || v == "*/" {
					continue
				} else if v[0:5] == "code1" {
					code1 = append(code1,v)
					code1Count++
				} else if v[0:5] == "code2" {
					code2 = append(code2,v)
					code2Count++
				} else if v[0:5] == "code3" {
					code3 = append(code3,v)
					code3Count++
				} else if v[0:5] == "error" {
					plerror = append(plerror,v)
					errorCount++
				}
			}
		}
		a := ""
		b := ""
		c := ""
		if code1Count < 10 {
			code1 = append(code1,fmt.Sprintf("code1000%d %s",code1Count+1,q))
			a = fmt.Sprintf("code1000%d",code1Count+1)
		} else {
			code1 = append(code1,fmt.Sprintf("code100%d %s",code1Count+1,q))
			a = fmt.Sprintf("code100%d",code1Count+1)
		}
		if code2Count < 10 {
			code2 = append(code2,fmt.Sprintf("code2000%d %s",code2Count+1,o))
			b = fmt.Sprintf("code2000%d",code2Count+1)
		} else {
			code2 = append(code2,fmt.Sprintf("code200%d %s",code2Count+1,o))
			b = fmt.Sprintf("code200%d",code2Count+1)
		}
		if code3Count < 10 {
			code3 = append(code3,fmt.Sprintf("code3000%d ByUser|$|%s",code3Count+1,w))
			c = fmt.Sprintf("code3000%d",code3Count+1)
		} else {
			code3 = append(code3,fmt.Sprintf("code300%d ByUser|$|%s",code3Count+1,w))
			c = fmt.Sprintf("code300%d",code3Count+1)
		}
		plerror = append(plerror,fmt.Sprintf("error(%s,%s,[%s]).",a,b,c))
		SaveData := "/*\n"
		for _,v := range code1 {
			SaveData = SaveData + v + "\n"
		}
		for _,v := range code2 {
			SaveData = SaveData + v + "\n"
		}
		for _,v := range code3 {
			SaveData = SaveData + v + "\n"
		}
		SaveData = SaveData + "*/\n"
		for _,v := range plerror {
			SaveData = SaveData + v + "\n"
		}
		os.Remove("data/RunQuestion.pl")
		f1,_ := os.Create("data/RunQuestion.pl")
		defer f1.Close()
		_,err=f1.Write([]byte(SaveData))
	}
}
func admin(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	id := r.FormValue("id")
	if id != "" {
		f, err := os.OpenFile("WWW/admin.html", os.O_RDONLY,0600)
		defer f.Close()
		if err !=nil {
			w.Write([]byte(err.Error()))
			return
		} else {
			HTMLByte,err:=ioutil.ReadAll(f)
			if err != nil{
				w.Write([]byte(err.Error()))
				return
			}
			code := ""
			for i,v := range AllNewRule {
				code = code + `<tr>
								<td>`+v.question+`</td>
								<td>`+v.operate+`</td>
								<td>`+v.way+`</td>
								<td>
									<div class="btn-group btn-group-sm">
										<button type="button" class="btn btn-primary" onclick="deleterule('`+strconv.Itoa(i)+`')">删除</button>
										<button type="button" class="btn btn-primary" onclick="rightrule('`+strconv.Itoa(i)+`')">确认</button>
									</div>
								</td>
							</tr>`
			}
			if code == "" {
				code = `<tr>
						<td>NULL</td>
						<td>NULL</td>
						<td>NULL</td>
						<td>
							<div class="btn-group btn-group-sm">
								<button type="button" class="btn btn-primary">删除</button>
								<button type="button" class="btn btn-primary">确认</button>
							</div>
						</td>
					</tr>`
			}
			w.Write([]byte(strings.Replace(string(HTMLByte),"{{table}}",code,-1)))
		}
	} else {
		w.Write([]byte("Error"))
	}
}
func NewRule1(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	newq := r.FormValue("newq")
	newo := r.FormValue("newo")
	neww := r.FormValue("neww")
	var temp NewRule
	temp.question = newq
	temp.operate = newo
	temp.way = neww
	AllNewRule = append(AllNewRule,temp)
	w.Write([]byte("我们已经记录你的反馈！感谢你的帮助！"))
}
func CheckCPP(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	codeBase64 := r.FormValue("code")
	
	decodeBytes, err := base64.StdEncoding.DecodeString(codeBase64)
    if err != nil {
        log.Fatalln(err)
    }
	code, err := url.QueryUnescape(string(decodeBytes))
	if err != nil{
		log.Println(err)
	}
	filename := ""
	myrand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i:=0;i<16;i++{
		filename = filename + strconv.Itoa(myrand.Intn(10))
	}
	f,err := os.Create("cache/"+filename+".cpp")
    defer f.Close()
    if err !=nil {
        log.Println(err)
    } else {
        _,err=f.Write([]byte(code))
        if err == nil {
			cppcheckout := ExecCommand("cppcheck "+"cache/"+filename+".cpp"+" --enable=all")
			outs := strings.Split(cppcheckout,"\n")
			resultStr := ""
			for _,v := range outs{
				for _,i := range AllCPPCheckError {
					if strings.Contains(v,i.data) {
						lineCount := GetBetweenStr(v,":","]")
						resultStr = resultStr + lineCount + " " + GetData_ProLog(i.code) + "\n"
					}
				}
			}
			if resultStr == "" {
				w.Write([]byte("not Find"))
			} else {
				w.Write([]byte(strings.Replace(resultStr,"\"","",-1)))
			}
		} else {
			w.Write([]byte("error"))
		}
		os.Remove("cache/"+filename+".cpp")
	}
}
func Init(){
	AllCPPCheckError = AllCPPCheckError[0:0]
	AllErrorInfo = AllErrorInfo[0:0]
	f, err := os.OpenFile("data/RunQuestion.pl", os.O_RDONLY,0600)
	defer f.Close()
	if err !=nil {
		log.Fatal(err)
		return
	} else {
		DataByte,err:=ioutil.ReadAll(f)
		if err != nil{
			log.Fatal(err)
			return
		}
		DataStr := string(DataByte)
		lines := strings.Split(DataStr,"\n")
		for count:=0 ; count < len(lines) ;count ++ {
			lines[count] = strings.TrimSpace(lines[count])
			lines[count] = strings.Replace(lines[count],"\\n","",-1)
		}
		for _,v := range lines {
			if v != "" {
				if v == "*/" {
					break
				} else if v != "/*"{
					temp := strings.Split(v," ")
					if len(temp) > 1 {
						var errorInfo ErrorInfo
						errorInfo.code = temp[0]
						errorInfo.data = temp[1]
						AllErrorInfo = append(AllErrorInfo,errorInfo)
					}
				}
			}
		}

	}
	f1, err := os.OpenFile("data/CPPCheck.pl", os.O_RDONLY,0600)
	defer f1.Close()
	if err !=nil {
		log.Fatal(err)
		return
	} else {
		DataByte,err:=ioutil.ReadAll(f1)
		if err != nil{
			log.Fatal(err)
			return
		}
		DataStr := string(DataByte)
		lines := strings.Split(DataStr,"\n")
		for count:=0 ; count < len(lines) ;count ++ {
			lines[count] = strings.TrimSpace(lines[count])
			lines[count] = strings.Replace(lines[count],"\\n","",-1)
		}
		for _,v := range lines {
			if v != "" {
				if v == "*/" {
					return
				} else if v != "/*"{
					temp := strings.Split(v,"@")
					if len(temp) > 1 {
						var errorInfo ErrorInfo
						errorInfo.code = temp[0]
						errorInfo.data = temp[1]
						AllCPPCheckError = append(AllCPPCheckError,errorInfo)
					}
				}
			}
		}

	}
}
func Index(w http.ResponseWriter, r *http.Request) {
	f, err := os.OpenFile("WWW/question.html", os.O_RDONLY,0600)
	defer f.Close()
	f1, _ := os.OpenFile("data/GCCChinese.pl", os.O_RDONLY,0600)
	TempByte,_:=ioutil.ReadAll(f1)
	TempArr := strings.Split(string(TempByte),"\n")
	defer f1.Close()
	if err !=nil {
		w.Write([]byte(err.Error()))
		return
	} else {
		HTMLByte,err:=ioutil.ReadAll(f)
		if err != nil{
			w.Write([]byte(err.Error()))
			return
		}
		questionStr := ""
		questionCount := 0
		operateStr := ""
		operateCount := 0
		for _,v := range AllErrorInfo {
			if v.code[0:5] == "code1" {
				questionCount++
				questionStr = questionStr + `<div class="radio">
												<label><input type="radio" name="`+v.code+`" id="`+v.code+`" value="`+v.code+`">`+v.data+`</label>
											</div>`
			} else if v.code[0:5] == "code2" {
				operateCount++
				operateStr = operateStr + `<div class="form-check">
												<label class="form-check-label">
												<input type="checkbox" class="form-check-input" value="`+v.code+`" id="`+v.code+`" name="`+v.code+`">`+v.data+`</label>
											</div>`
			}
		}
		HTML := strings.Replace(string(HTMLByte),"{{Question}}",questionStr,-1)
		HTML = strings.Replace(HTML,"{{Operate}}",operateStr,-1)
		HTML = strings.Replace(HTML,"{{QuestionCount}}",strconv.Itoa(questionCount),-1)
		HTML = strings.Replace(HTML,"{{OperateCount}}",strconv.Itoa(operateCount),-1)
		HTML = strings.Replace(HTML,"{{CountQurstion}}",strconv.Itoa(len(AllErrorInfo)),-1)
		HTML = strings.Replace(HTML,"{{CountError}}",strconv.Itoa(len(TempArr)),-1)
		w.Write([]byte(HTML))
	}
}
func Code(w http.ResponseWriter, r *http.Request) {
	f, err := os.OpenFile("WWW/code.html", os.O_RDONLY,0600)
	defer f.Close()
	if err !=nil {
		w.Write([]byte(err.Error()))
		return
	} else {
		HTMLByte,err:=ioutil.ReadAll(f)
		if err != nil{
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(HTMLByte)
	}
}
func CheckCode(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.FormValue("code")
	w.Write([]byte(code+`<br>`+GetCodeChinese(code)+`<br>`+GetCodeWay(code)))
}
func GetResult(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	question := r.FormValue("question")
	operate := r.FormValue("operate")
	operates := strings.Split(operate,",")
	if len(operates) > 1 {
		result := ""
		count := 0
		for _,v := range operates {
			if v == "" {
				continue
			}
			temp := GetResult_ProLog(question,v)
			if temp != "not find" {
				result = result + Result2Str(temp) + "|@|"
				count ++
			}
		}
		if result == "" {
			w.Write([]byte("0|@|not find"))
		} else {
			w.Write([]byte(strconv.Itoa(count)+"|@|"+result))
		}
	} else {
		temp := GetResult_ProLog(question,operates[0])
		if temp != "not find" {
			w.Write([]byte("1|@|"+Result2Str(temp)))
		} else {
			w.Write([]byte("0|@|"+temp))
		}
	}
}
func Result2Str(result string)string{
	result = strings.Replace(result,"[","",-1)
	result = strings.Replace(result,"]","",-1)
	temp := strings.Split(result,",")
	str := ""
	for _,v := range temp {
		for _,i := range AllErrorInfo {
			if v == i.code {
				str = str + i.data +"|#|"
				break
			}
		}
	}
	return str
}
func GetCodeChinese(code string) string {
	f, err := os.OpenFile("data/GCCChinese.pl", os.O_RDONLY,0600)
	defer f.Close()
	if err !=nil {
		return "Can't Find File!"
	} else {
		DataByte,err:=ioutil.ReadAll(f)
		if err != nil{
			return "Can't Open File"
		}
		m := golog.NewMachine().Consult(string(DataByte))
		solutions := m.ProveAll(`codeToChinese(`+code+`,X).`)
		for _, solution := range solutions {
			return fmt.Sprintf("%s",solution.ByName_("X"))
		}
	}
	return "not find code"
}
func GetCodeWay(code string) string {
	f, err := os.OpenFile("data/GCCWay.pl", os.O_RDONLY,0600)
	defer f.Close()
	if err !=nil {
		return "Can't Find File!"
	} else {
		DataByte,err:=ioutil.ReadAll(f)
		if err != nil{
			return "Can't Open File"
		}
		m := golog.NewMachine().Consult(string(DataByte))
		solutions := m.ProveAll(`codeToWay(`+code+`,X).`)
		for _, solution := range solutions {
			wayCode := fmt.Sprintf("%s",solution.ByName_("X"))
			f1, err := os.OpenFile("data/Table.txt", os.O_RDONLY,0600)
			defer f1.Close()
			if err !=nil {
				return "Can't Find File!"
			} else {
				DataByte1,err:=ioutil.ReadAll(f1)
				if err != nil{
					return "Can't Open File"
				}
				lines := strings.Split(string(DataByte1),"\n")
				for _,v := range lines {
					if v != "" {
						temp := strings.Split(v,"|@|")
						if temp[0] == wayCode {
							return temp[1]
						}
					}
				}
			}
		}
	}
	return "not find code"
}
func GetResult_ProLog(code1 string,code2 string) string {
	f, err := os.OpenFile("data/RunQuestion.pl", os.O_RDONLY,0600)
	defer f.Close()
	if err !=nil {
		return "Can't Find File!"
	} else {
		DataByte,err:=ioutil.ReadAll(f)
		if err != nil{
			return "Can't Open File"
		}
		m := golog.NewMachine().Consult(string(DataByte))
		solutions := m.ProveAll(`error(`+ code1 +`,`+ code2 +`,X).`)
		for _, solution := range solutions {
			return fmt.Sprintf("%s",solution.ByName_("X"))
		}
	}
	return "not find"
}
func ExecCommand(strCommand string)(string){
    cmd := exec.Command("/bin/bash", "-c", strCommand)
 
 
    stdout, _ := cmd.StderrPipe()
    if err := cmd.Start(); err != nil{
        log.Println("Execute failed when Start:" + err.Error())
        return ""
    }
 
    out_bytes, _ := ioutil.ReadAll(stdout)
    stdout.Close()
 
    if err := cmd.Wait(); err != nil {
        log.Println("Execute failed when Wait:" + err.Error())
        return ""
    }
    return string(out_bytes)
}
func GetData_ProLog(code1 string) string {
	f, err := os.OpenFile("data/CPPCheck.pl", os.O_RDONLY,0600)
	defer f.Close()
	if err !=nil {
		return "Can't Find File!"
	} else {
		DataByte,err:=ioutil.ReadAll(f)
		if err != nil{
			return "Can't Open File"
		}
		m := golog.NewMachine().Consult(string(DataByte))
		solutions := m.ProveAll(`error(`+ code1 +`,X,Y,Z).`)
		for _, solution := range solutions {
			return fmt.Sprintf("%s %s %s",solution.ByName_("X"),solution.ByName_("Y"),solution.ByName_("Z"))
		}
	}
	return "not find"
}
func GetBetweenStr(str, start, end string) string {
    n := strings.Index(str, start)
    if n == -1 {
        n = 0
    } else {
        n = n + len(start)  // 增加了else，不加的会把start带上
    }
    str = string([]byte(str)[n:])
    m := strings.Index(str, end)
    if m == -1 {
        m = len(str)
    }
    str = string([]byte(str)[:m])
    return str
}
func remove(slice []interface{}, elem interface{}) []interface{}{
    if len(slice) == 0 {
        return slice
    }
    for i, v := range slice {
        if v == elem {
            slice = append(slice[:i], slice[i+1:]...)
            return remove(slice,elem)
            break
        }
    }
    return slice
}