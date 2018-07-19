FROM dregistry:5000/bestv/otttest/common/alpine:3.6
LABEL maintainer "wang.min@bestv.com.cn"

# Copy our source code into the container.
WORKDIR /pixpress
ADD pixpress /pixpress/pixpress
ADD public /pixpress
ADD templates /pixpress/app/views/

# Expose a docker interface to our binary.
EXPOSE 7611
ENTRYPOINT ["/pixpress/pixpress", "web"]
