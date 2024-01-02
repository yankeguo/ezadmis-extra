package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/yankeguo/ezadmis"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

const (
	KeyDeletionAllowed       = "ezadmis.yankeguo.github.io/deletion-allowed"
	KeyDeletionAllowedLegacy = "ezadmis.guoyk93.github.io/deletion-allowed"
)

func main() {
	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))

	s := ezadmis.NewWebhookServer(
		ezadmis.WebhookServerOptions{
			Debug: debug,
			Handler: func(
				ctx context.Context,
				request *admissionv1.AdmissionRequest,
				rw ezadmis.WebhookResponseWriter,
			) (err error) {
				var buf []byte
				if buf, err = request.OldObject.MarshalJSON(); err != nil {
					return
				}
				var ns corev1.Namespace
				if err = json.Unmarshal(buf, &ns); err != nil {
					return
				}
				if ns.Annotations != nil {
					if ok, _ := strconv.ParseBool(ns.Annotations[KeyDeletionAllowed]); ok {
						return
					}
					// legacy annotation
					if ok, _ := strconv.ParseBool(ns.Annotations[KeyDeletionAllowedLegacy]); ok {
						return
					}
				}
				rw.Deny("missing annotation '" + KeyDeletionAllowed + "', deletion of namespace is denied")
				return
			},
		},
	)

	err := s.ListenAndServeGracefully()

	if err != nil {
		log.Println("exited with error:", err.Error())
		os.Exit(1)
	}

}
