package main

import (
	"flag"
	"github.com/be-ys/pzem-004t-v3/pzem"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
)

const defaultHost = "0.0.0.0"
const defaultPort = 2112

var listenHost = flag.String("host", defaultHost, "Hostname to bind to")
var listenPort = flag.Int("port", defaultPort, "Port to listen request on")
var serialport = flag.String("serialPort", "", "Serial port used to communicate with PZEM-004T")
var shouldResetEnegy = flag.Bool("resetEnergy", false, "Should the energy value be reseted at start")

var powerMeter pzem.Probe

var powerMeterVoltage = promauto.NewGaugeFunc(prometheus.GaugeOpts{
	Name: "power_meter_voltage",
	Help: "Current voltage at power meter (V)",
},
	func() float64 {
		value, _ := powerMeter.Voltage()
		return float64(value)
	})

var powerMeterAmps = promauto.NewGaugeFunc(prometheus.GaugeOpts{
	Name: "power_meter_amps",
	Help: "Current amps consumed (A)",
}, func() float64 {
	value, _ := powerMeter.Intensity()
	return float64(value)
})

var powerMeterFrequency = promauto.NewGaugeFunc(prometheus.GaugeOpts{
	Name: "power_meter_frequency",
	Help: "Current measured frequency of AC electricity (Hz)",
}, func() float64 {
	value, _ := powerMeter.Frequency()
	return float64(value)
})

var powerMeterActivePower = promauto.NewGaugeFunc(prometheus.GaugeOpts{
	Name: "power_meter_active_power",
	Help: "Currently used power (W)",
}, func() float64 {
	value, _ := powerMeter.Power()
	return float64(value)
})

var powerMeterPowerFactor = promauto.NewGaugeFunc(prometheus.GaugeOpts{
	Name: "power_meter_power_factor",
	Help: "Measured power factor (no unit)",
}, func() float64 {
	value, _ := powerMeter.PowerFactor()
	return float64(value)
})

var powerMeterActiveEnergy = promauto.NewCounterFunc(prometheus.CounterOpts{
	Name: "power_meter_active_energy",
	Help: "Consumed Energy (kWh)",
}, func() float64 {
	value, _ := powerMeter.Energy()
	return float64(value)
})

func main() {

	flag.Parse()

	log.Printf("Starting PZEM-004T Prometheus exporter")
	log.Printf("Initialize PZEM-004T library with device on serial port '%s'", *serialport)

	var err error

	powerMeter, err = pzem.Setup(
		pzem.Config{
			Port:  *serialport,
			Speed: 9600,
		},
	)

	if err != nil {
		log.Fatalf("Could not initialize device on port '%s': %s", *serialport, err)
	}

	if *shouldResetEnegy {
		log.Printf("Trying to reset energy value...")
		err = powerMeter.ResetEnergy()

		if err != nil {
			log.Fatalf("Could not reset energy value in PZEM-004T: %s", err)
		}
		log.Printf("Energy value reseted !")
	}

	log.Printf("Register all metrics...")
	prometheus.MustRegister(powerMeterVoltage)
	prometheus.MustRegister(powerMeterAmps)
	prometheus.MustRegister(powerMeterFrequency)
	prometheus.MustRegister(powerMeterActivePower)
	prometheus.MustRegister(powerMeterPowerFactor)
	prometheus.MustRegister(powerMeterActiveEnergy)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe((*listenHost)+":"+strconv.Itoa(*listenPort), nil)

}
