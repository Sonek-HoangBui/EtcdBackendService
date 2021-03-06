package String2Int64Service

import (
	"context"
	"errors"
	"log"

	"github.com/OpenStars/EtcdBackendService/String2Int64Service/s2i64kv/thrift/gen-go/OpenStars/Common/S2I64KV"
	"github.com/OpenStars/EtcdBackendService/String2Int64Service/s2i64kv/transports"
	"github.com/OpenStars/GoEndpointManager"
	"github.com/OpenStars/GoEndpointManager/GoEndpointBackendManager"
)

type String2Int64Service struct {
	host        string
	port        string
	sid         string
	epm         GoEndpointBackendManager.EndPointManagerIf
	etcdManager *GoEndpointManager.EtcdBackendEndpointManager
}

func (m *String2Int64Service) PutData(key string, value int64) error {

	if m.etcdManager != nil {
		h, p, err := m.etcdManager.GetEndpoint(m.sid)
		if err != nil {
			log.Println("EtcdManager get endpoints", "err", err)
		} else {
			m.host = h
			m.port = p
		}
	}

	client := transports.GetS2I64CompactClient(m.host, m.port)

	if client == nil || client.Client == nil {
		return errors.New("Can not connect to model")
	}

	tkey := S2I64KV.TKey(key)
	tvalue := &S2I64KV.TI64Value{
		Value: value,
	}

	_, err := client.Client.(*S2I64KV.TString2I64KVServiceClient).PutData(context.Background(), tkey, tvalue)

	if err != nil {
		return errors.New("String2Int64Service sid: " + m.sid + " address: " + m.host + ":" + m.port + " err: " + err.Error())
	}
	defer client.BackToPool()
	return nil
}

func (m *String2Int64Service) GetData(key string) (int64, error) {

	if m.etcdManager != nil {
		h, p, err := m.etcdManager.GetEndpoint(m.sid)
		if err != nil {
			log.Println("EtcdManager get endpoints", "err", err)
		} else {
			m.host = h
			m.port = p
		}
	}

	client := transports.GetS2I64CompactClient(m.host, m.port)

	if client == nil || client.Client == nil {
		return -1, errors.New("Can not connect to model")
	}

	tkey := S2I64KV.TKey(key)
	r, err := client.Client.(*S2I64KV.TString2I64KVServiceClient).GetData(context.Background(), tkey)

	if err != nil {

		return -1, errors.New("String2Int64Service sid: " + m.sid + " address: " + m.host + ":" + m.port + " err: " + err.Error())
	}
	defer client.BackToPool()
	if r == nil || r.Data == nil || r.ErrorCode != S2I64KV.TErrorCode_EGood || r.Data.Value <= 0 {

		return -1, errors.New("Can not found key")
	}
	return r.Data.Value, nil
}

func (m *String2Int64Service) CasData(key string, value int64) (sucess bool, oldvalue int64, err error) {
	if m.etcdManager != nil {
		h, p, err := m.etcdManager.GetEndpoint(m.sid)
		if err != nil {
			log.Println("EtcdManager get endpoints", "err", err)
		} else {
			m.host = h
			m.port = p
		}
	}

	client := transports.GetS2I64CompactClient(m.host, m.port)
	if client == nil || client.Client == nil {
		return false, -1, errors.New("Can not connect to model")
	}

	var aCas = &S2I64KV.TCasValue{OldValue: 0, NewValue_: value}
	r, err := client.Client.(*S2I64KV.TString2I64KVServiceClient).CasData(context.Background(), S2I64KV.TKey(key), aCas)
	if err != nil {
		return false, -1, errors.New("String2Int64Service sid: " + m.sid + " address: " + m.host + ":" + m.port + " err: " + err.Error())
	}
	defer client.BackToPool()

	if r != nil && r.GetOldValue() == 0 {
		return true, value, nil
	}
	return false, r.GetOldValue(), nil
}

func (m *String2Int64Service) handleEventChangeEndpoint(ep *GoEndpointBackendManager.EndPoint) {
	m.host = ep.Host
	m.port = ep.Port
	log.Println("Change config endpoint serviceID", ep.ServiceID, m.host, ":", m.port)
}

func NewString2Int64Service(serviceID string, etcdServers []string, defaultEndpoint GoEndpointBackendManager.EndPoint) String2Int64ServiceIf {
	s2isv := &String2Int64Service{
		host:        defaultEndpoint.Host,
		port:        defaultEndpoint.Port,
		sid:         defaultEndpoint.ServiceID,
		etcdManager: GoEndpointManager.GetEtcdBackendEndpointManagerSingleton(etcdServers),
	}

	if s2isv.etcdManager == nil {
		return s2isv
	}
	err := s2isv.etcdManager.SetDefaultEntpoint(serviceID, defaultEndpoint.Host, defaultEndpoint.Port)
	if err != nil {
		log.Println("SetDefaultEndpoint sid", serviceID, "err", err)
		return s2isv
	}
	// s2isv.etcdManager.GetAllEndpoint(serviceID)
	return s2isv
}
