package webhook

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"k8s.io/api/admission/v1beta1"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/klog"
	// v1 "k8s.io/kubernetes/pkg/apis/core/v1"
)

// Server webhook server
type Server struct {
	Svr *http.Server
}

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()

	// (https://github.com/kubernetes/kubernetes/issues/57982)
	defaulter = runtime.ObjectDefaulter(runtimeScheme)
)

func init() {
	_ = corev1.AddToScheme(runtimeScheme)
	_ = admissionregistrationv1beta1.AddToScheme(runtimeScheme)
	// defaulting with webhooks:
	// https://github.com/kubernetes/kubernetes/issues/57982
	// _ = v1.AddToScheme(runtimeScheme)
}

// New 生成实例
func New(param WhSvrParameters) (*Server, error) {
	pair, err := tls.LoadX509KeyPair(param.CertFile, param.KeyFile)
	if err != nil {
		return nil, err
	}
	return &Server{
		Svr: &http.Server{
			Addr:      fmt.Sprintf(":%v", param.Port),
			TLSConfig: &tls.Config{Certificates: []tls.Certificate{pair}},
		},
	}, nil
}

func admissionRequired(ignoredList []string, admissionAnnotationKey string, metadata *metav1.ObjectMeta) bool {
	// skip special kubernetes system namespaces
	for _, namespace := range ignoredList {
		if metadata.Namespace == namespace {
			klog.Infof("Skip validation for %v for it's in special namespace:%v",
				metadata.Name, metadata.Namespace)
			return false
		}
	}

	annotations := metadata.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}

	var required bool
	switch strings.ToLower(annotations[admissionAnnotationKey]) {
	default:
		required = true
	case "n", "no", "false", "off":
		required = false
	}
	return required
}

func mutationRequired(ignoredList []string, metadata *metav1.ObjectMeta) bool {
	required := admissionRequired(ignoredList, admissionWebhookAnnotationMutateKey, metadata)
	annotations := metadata.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	status := annotations[admissionWebhookAnnotationStatusKey]

	if strings.ToLower(status) == "mutated" {
		required = false
	}

	klog.Infof("Mutation policy for %v/%v: required:%v", metadata.Namespace, metadata.Name, required)
	return required
}

func validationRequired(ignoredList []string, metadata *metav1.ObjectMeta) bool {
	required := admissionRequired(ignoredList, admissionWebhookAnnotationValidateKey, metadata)
	klog.Infof("Validation policy for %v/%v: required:%v", metadata.Namespace, metadata.Name, required)
	return required
}

func updateAnnotation(target map[string]string, added map[string]string) (patch []patchOperation) {
	for key, value := range added {
		if target == nil || target[key] == "" {
			target = map[string]string{}
			patch = append(patch, patchOperation{
				Op:   "add",
				Path: "/metadata/annotations",
				Value: map[string]string{
					key: value,
				},
			})
		} else {
			patch = append(patch, patchOperation{
				Op:    "replace",
				Path:  "/metadata/annotations/" + key,
				Value: value,
			})
		}
	}
	return patch
}

func updateLabels(target map[string]string, added map[string]string) (patch []patchOperation) {
	values := make(map[string]string)
	for key, value := range added {
		if target == nil || target[key] == "" {
			values[key] = value
		}
	}
	patch = append(patch, patchOperation{
		Op:    "add",
		Path:  "/metadata/labels",
		Value: values,
	})
	return patch
}

func createPatch(availableAnnotations map[string]string, annotations map[string]string,
	availableLabels map[string]string, labels map[string]string) ([]byte, error) {
	var patch []patchOperation

	patch = append(patch, updateAnnotation(availableAnnotations, annotations)...)
	patch = append(patch, updateLabels(availableLabels, labels)...)

	return json.Marshal(patch)
}

// validate deployments and services
func (whsvr *Server) validate(ar *v1beta1.AdmissionReview, log *bytes.Buffer) *v1beta1.AdmissionResponse {
	req := ar.Request
	var (
		availableLabels                 map[string]string
		objectMeta                      *metav1.ObjectMeta
		resourceNamespace, resourceName string
	)

	log.WriteString(fmt.Sprintf("\n======begin Admission for Namespace=[%v], Kind=[%v], Name=[%v]======",
		req.Namespace, req.Kind.Kind, req.Name))

	switch req.Kind.Kind {
	case "Deployment":
		var deployment appsv1.Deployment
		if err := json.Unmarshal(req.Object.Raw, &deployment); err != nil {
			log.WriteString(fmt.Sprintf("\nCould not unmarshal raw object: %v", err))
			klog.Errorf(log.String())
			return &v1beta1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
		resourceName, resourceNamespace, objectMeta = deployment.Name, deployment.Namespace, &deployment.ObjectMeta
		availableLabels = deployment.Labels
	case "Service":
		var service corev1.Service
		if err := json.Unmarshal(req.Object.Raw, &service); err != nil {
			log.WriteString(fmt.Sprintf("\nCould not unmarshal raw object: %v", err))
			klog.Errorf(log.String())
			return &v1beta1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
		resourceName, resourceNamespace, objectMeta = service.Name, service.Namespace, &service.ObjectMeta
		availableLabels = service.Labels
	//其他不支持的类型
	default:
		msg := fmt.Sprintf("\nNot support for this Kind of resource  %v", req.Kind.Kind)
		log.WriteString(msg)
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: msg,
			},
		}
	}

	if !validationRequired(ignoredNamespaces, objectMeta) {
		log.WriteString(fmt.Sprintf("Skipping validation for %s/%s due to policy check",
			resourceNamespace, resourceName))
		return &v1beta1.AdmissionResponse{
			Allowed: true,
		}
	}

	allowed := true
	var result *metav1.Status
	log.WriteString(fmt.Sprintf("available labels: %s ", availableLabels))
	log.WriteString(fmt.Sprintf("required labels: %s", requiredLabels))
	for _, rl := range requiredLabels {
		if _, ok := availableLabels[rl]; !ok {
			allowed = false
			result = &metav1.Status{
				Reason: "required labels are not set",
			}
			break
		}
	}

	return &v1beta1.AdmissionResponse{
		Allowed: allowed,
		Result:  result,
	}
}

// mutate mutation process
func (whsvr *Server) mutate(ar *v1beta1.AdmissionReview, log *bytes.Buffer) *v1beta1.AdmissionResponse {
	req := ar.Request
	var (
		availableLabels, availableAnnotations map[string]string
		objectMeta                            *metav1.ObjectMeta
		resourceNamespace, resourceName       string
	)

	log.WriteString(fmt.Sprintf("\n======begin Admission for Namespace=[%v], Kind=[%v], Name=[%v]======",
		req.Namespace, req.Kind.Kind, req.Name))
	log.WriteString("\n>>>>>>" + req.Kind.Kind)

	switch req.Kind.Kind {
	case "Deployment":
		var deployment appsv1.Deployment
		if err := json.Unmarshal(req.Object.Raw, &deployment); err != nil {
			log.WriteString(fmt.Sprintf("\nCould not unmarshal raw object: %v", err))
			klog.Errorf(log.String())
			return &v1beta1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
		resourceName, resourceNamespace, objectMeta = deployment.Name, deployment.Namespace, &deployment.ObjectMeta
		availableLabels = deployment.Labels
	case "Service":
		var service corev1.Service
		if err := json.Unmarshal(req.Object.Raw, &service); err != nil {
			log.WriteString(fmt.Sprintf("\nCould not unmarshal raw object: %v", err))
			klog.Errorf(log.String())
			return &v1beta1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
		}
		resourceName, resourceNamespace, objectMeta = service.Name, service.Namespace, &service.ObjectMeta
		availableLabels = service.Labels
	//其他不支持的类型
	default:
		msg := fmt.Sprintf("\nNot support for this Kind of resource  %v", req.Kind.Kind)
		log.WriteString(msg)
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: msg,
			},
		}
	}

	if !mutationRequired(ignoredNamespaces, objectMeta) {
		log.WriteString(fmt.Sprintf("Skipping validation for %s/%s due to policy check", resourceNamespace, resourceName))
		return &v1beta1.AdmissionResponse{
			Allowed: true,
		}
	}

	annotations := map[string]string{admissionWebhookAnnotationStatusKey: "mutated"}
	patchBytes, err := createPatch(availableAnnotations, annotations, availableLabels, addLabels)
	if err != nil {
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	log.WriteString(fmt.Sprintf("AdmissionResponse: patch=%v\n", string(patchBytes)))
	return &v1beta1.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}

// Serve method for webhook server
func (whsvr *Server) Serve(w http.ResponseWriter, r *http.Request) {
	//记录日志
	var log bytes.Buffer

	//读取从ApiServer过来的数据放到body
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		log.WriteString("empty body")
		klog.Info(log.String())
		//返回状态码400
		//如果在Apiserver调用此Webhook返回是400，说明APIServer自己传过来的数据是空
		http.Error(w, log.String(), http.StatusBadRequest)
		return
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		log.WriteString(fmt.Sprintf("Content-Type=%s, expect `application/json`", contentType))
		klog.Errorf(log.String())
		//如果在Apiserver调用此Webhook返回是415，说明APIServer自己传过来的数据不是json格式，处理不了
		http.Error(w, log.String(), http.StatusUnsupportedMediaType)
		return
	}

	var admissionResponse *v1beta1.AdmissionResponse
	ar := v1beta1.AdmissionReview{}
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		//组装错误信息
		log.WriteString(fmt.Sprintf("\nCan't decode body,error info is :  %s", err.Error()))
		klog.Errorln(log.String())
		//返回错误信息，形式表现为资源创建会失败，
		admissionResponse = &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: log.String(),
			},
		}
	} else {
		fmt.Println(r.URL.Path)
		if r.URL.Path == "/mutate" {
			admissionResponse = whsvr.mutate(&ar, &log)
		} else if r.URL.Path == "/validate" {
			admissionResponse = whsvr.validate(&ar, &log)
		}
	}

	admissionReview := v1beta1.AdmissionReview{}
	if admissionResponse != nil {
		admissionReview.Response = admissionResponse
		if ar.Request != nil {
			admissionReview.Response.UID = ar.Request.UID
		}
	}

	resp, err := json.Marshal(admissionReview)
	if err != nil {
		log.WriteString(fmt.Sprintf("\nCan't encode response: %v", err))
		http.Error(w, log.String(), http.StatusInternalServerError)
	}
	klog.Infof("Ready to write response ...")
	if _, err := w.Write(resp); err != nil {
		log.WriteString(fmt.Sprintf("\nCan't write response: %v", err))
		http.Error(w, log.String(), http.StatusInternalServerError)
	}

	log.WriteString("\n======ended Admission already write to response======")
	//东八区时间
	datetime := time.Now().In(time.FixedZone("GMT", 8*3600)).Format("2006-01-02 15:04:05")
	//最后打印日志
	klog.Infof(datetime + " " + log.String())
}
