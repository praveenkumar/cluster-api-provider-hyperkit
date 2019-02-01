// NOTE: Boilerplate only.  Ignore this file.

// Package v1alpha1 contains API Schema definitions for the hyperkitproviderconfig v1alpha1 API group
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=github.com/praveenkumar/cluster-api-provider-hyperkit/pkg/apis/hyperkitproviderconfig
// +k8s:defaulter-gen=TypeMeta
// +groupName=hyperkitproviderconfig.k8s.io
package v1alpha1

import (
	"bytes"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/runtime/scheme"

	clusterv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: "hyperkitproviderconfig.k8s.io", Version: "v1alpha1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
)

// HyperkitProviderConfigCodec contains encoder/decoder to convert this types from/to serialize data
// +k8s:deepcopy-gen=false
type HyperkitProviderConfigCodec struct {
	encoder runtime.Encoder
	decoder runtime.Decoder
}

// NewScheme creates a new Scheme
func NewScheme() (*runtime.Scheme, error) {
	return SchemeBuilder.Build()
}

// NewCodec returns a encode/decoder for this API
func NewCodec() (*HyperkitProviderConfigCodec, error) {
	scheme, err := NewScheme()
	if err != nil {
		return nil, err
	}
	codecFactory := serializer.NewCodecFactory(scheme)
	encoder, err := newEncoder(&codecFactory)
	if err != nil {
		return nil, err
	}
	codec := HyperkitProviderConfigCodec{
		encoder: encoder,
		decoder: codecFactory.UniversalDecoder(SchemeGroupVersion),
	}
	return &codec, nil
}

// DecodeFromProviderSpec decodes a serialised ProviderConfig into an object
func (codec *HyperkitProviderConfigCodec) DecodeFromProviderSpec(providerConfig clusterv1alpha1.ProviderSpec, out runtime.Object) error {
	if providerConfig.Value != nil {
		_, _, err := codec.decoder.Decode(providerConfig.Value.Raw, nil, out)
		if err != nil {
			return fmt.Errorf("decoding failure: %v", err)
		}
	}
	return nil
}

// EncodeToProviderSpec encodes an object into a serialised ProviderConfig
func (codec *HyperkitProviderConfigCodec) EncodeToProviderSpec(in runtime.Object) (*clusterv1alpha1.ProviderSpec, error) {
	var buf bytes.Buffer
	if err := codec.encoder.Encode(in, &buf); err != nil {
		return nil, fmt.Errorf("encoding failed: %v", err)
	}
	return &clusterv1alpha1.ProviderSpec{
		Value: &runtime.RawExtension{Raw: buf.Bytes()},
	}, nil
}

// EncodeProviderStatus encodes an object into serialised data
func (codec *HyperkitProviderConfigCodec) EncodeProviderStatus(in runtime.Object) (*runtime.RawExtension, error) {
	var buf bytes.Buffer
	if err := codec.encoder.Encode(in, &buf); err != nil {
		return nil, fmt.Errorf("encoding failed: %v", err)
	}

	return &runtime.RawExtension{Raw: buf.Bytes()}, nil
}

// DecodeProviderStatus decodes a serialised providerStatus into an object
func (codec *HyperkitProviderConfigCodec) DecodeProviderStatus(providerStatus *runtime.RawExtension, out runtime.Object) error {
	if providerStatus != nil {
		_, _, err := codec.decoder.Decode(providerStatus.Raw, nil, out)
		if err != nil {
			return fmt.Errorf("decoding failure: %v", err)
		}
	}
	return nil
}

func newEncoder(codecFactory *serializer.CodecFactory) (runtime.Encoder, error) {
	serializerInfos := codecFactory.SupportedMediaTypes()
	if len(serializerInfos) == 0 {
		return nil, fmt.Errorf("unable to find any serlializers")
	}
	encoder := codecFactory.EncoderForVersion(serializerInfos[0].Serializer, SchemeGroupVersion)
	return encoder, nil
}