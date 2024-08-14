# range func

https://go.dev/wiki/RangefuncExperiment

~~~
google\play\acquire.go
76:   for v := range a.message.Get(1) {

google\play\delivery.go
99:      for v := range d.Message.Get(15) {
110:      for v := range d.Message.Get(4) {

google\play\details.go
88:      for v := range v.Get(17) {

widevine\cdm.go
53:   for container := range license.Get(3) { // KeyContainer key
~~~
