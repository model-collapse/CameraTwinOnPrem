FROM modelcollapse/rasbian-opencv-go:v0.1.0

RUN git clone https://github.com/model-collapse/CameraTwinOnPrem.git
RUN cd CameraTwinOnPrem/CameraTwinGateway & go build -o camera-twin-onprem
RUN mkdir /deploy & cp CameraTwinOnPrem/CameraTwinGateway/camera-twin-onprem /deploy

ENV SERVICE_PORT 8964
ENV MQTT_PORT 8883
ENV MQTT_HOST mq

CMD ["/deploy/camera-twin-onprem"]