FROM golang:1.16 

RUN git clone https://github.com/4thlabs/perfmon
RUN cd perfmon/cmd/perf && go build
WORKDIR ./perfmon/cmd/perf

CMD [ "./perf" ]