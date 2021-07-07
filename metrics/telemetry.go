// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

/*
 * AMF Statistics exposing to promethus
 *
 */

package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//AmfStats captures AMF level stats
type AmfStats struct {
	ngapMsg           *prometheus.CounterVec
	gnbSessionProfile *prometheus.GaugeVec
}

var amfStats *AmfStats

func initAmfStats() *AmfStats {
	return &AmfStats{
		ngapMsg: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "ngap_messages_total",
			Help: "ngap interface counters",
		}, []string{"amf_id", "msg_type", "direction", "result", "reason"}),

		gnbSessionProfile: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "gnb_session_profile",
			Help: "gNB session Profile",
		}, []string{"id", "ip", "state", "tac"}),
	}
}

func (ps *AmfStats) register() error {
	if err := prometheus.Register(ps.ngapMsg); err != nil {
		return err
	}
	if err := prometheus.Register(ps.gnbSessionProfile); err != nil {
		return err
	}
	return nil
}

func init() {
	amfStats = initAmfStats()

	if err := amfStats.register(); err != nil {
		log.Panicln("AMF Stats register failed")
	}
}

//InitMetrics initialises SMF stats
func InitMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9089", nil)
}

//IncrementNgapMsgStats increments message level stats
func IncrementNgapMsgStats(amfID, msgType, direction, result, reason string) {
	amfStats.ngapMsg.WithLabelValues(amfID, msgType, direction, result, reason).Inc()
}

//SetSessProfileStats maintains Session profile info
func SetGnbSessProfileStats(id, ip, state string, tac uint64) {
	amfStats.gnbSessionProfile.WithLabelValues(id, ip, state).Set(float64(tac))
}

